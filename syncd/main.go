// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
    "log"
    "flag"
    "fmt"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/util/gopath"
    "github.com/dreamans/syncd"
    "github.com/dreamans/syncd/router/route"
    "github.com/Unknwon/goconfig"
)

var (
    helpFlag        bool
    syncdIniFlag    string
    versionFlag     bool
    configFile      *goconfig.ConfigFile
)

func init() {
    gin.SetMode(gin.ReleaseMode)

    flag.BoolVar(&helpFlag, "h", false, "This help")
    flag.StringVar(&syncdIniFlag, "c", "", "Set configuration file `file`")
    flag.BoolVar(&versionFlag, "v", false, "Version number")

    flag.Usage = usage
    flag.Parse()
}

func usage() {
    fmt.Printf("Usage: syncd [-c filename]\n\nOptions:\n")
    flag.PrintDefaults()
}

func initCfg() {
    var err error
    syncdIni := findSyncdIniFile()
    configFile, err = goconfig.LoadConfigFile(syncdIni)
    if err != nil {
        log.Fatalf("load config file failed, %s\n", err.Error())
    }
    outputInfo("Config Loaded", syncdIni)
}

func configIntOrDefault(section, key string, useDefault int) int {
    val, err := configFile.Int(section, key)
    if err != nil {
        return useDefault
    }
    return val
}

func configOrDefault(section, key, useDefault string) string {
    val, err := configFile.GetValue(section, key)
    if err != nil {
        return useDefault
    }
    return val
}

func findSyncdIniFile() string {
    if syncdIniFlag != "" {
        return syncdIniFlag
    }
    currPath, _ := gopath.CurrentPath()
    parentPath, _ := gopath.CurrentParentPath()
    scanPath := []string{
        "/etc",
        currPath,
        fmt.Sprintf("%s/etc", currPath),
        fmt.Sprintf("%s/etc", parentPath),
    }

    for _, path := range scanPath {
        iniFile := path + "/syncd.ini"
        if gopath.Exists(iniFile) && gopath.IsFile(iniFile) {
            return iniFile
        }
    }

    return "./syncd.ini"
}

func outputInfo(tag string, value interface{}) {
    fmt.Printf("%-18s    %v\n", tag + ":", value)
}

func welcome() {
    fmt.Println("                                          __")
    fmt.Println("   _____   __  __   ____     _____   ____/ /")
    fmt.Println("  / ___/  / / / /  / __ \\   / ___/  / __  / ")
    fmt.Println(" (__  )  / /_/ /  / / / /  / /__   / /_/ /  ")
    fmt.Println("/____/   \\__, /  /_/ /_/   \\___/   \\__,_/   ")
    fmt.Println("        /____/                              ")
    fmt.Println("")
    outputInfo("Service", "syncd")
    outputInfo("Version", syncd.Version)
}

func main() {
    if helpFlag {
        flag.Usage()
        os.Exit(0)
    }
    if versionFlag {
        fmt.Printf("syncd/%s\n", syncd.Version)
        os.Exit(0)
    }

    welcome()

    initCfg()

    cfg := &syncd.Config{
        Serve: &syncd.ServeConfig{
            Addr: configOrDefault("serve", "addr", "8868"),
            FeServeEnable: configIntOrDefault("serve", "fe_serve_enable", 1),
            ReadTimeout: configIntOrDefault("serve", "read_timeout", 300),
            WriteTimeout: configIntOrDefault("serve", "write_timeout", 300),
            IdleTimeout: configIntOrDefault("serve", "idle_timeout", 300),
        },
        Db: &syncd.DbConfig{
            Unix: configOrDefault("database", "unix", ""),
            Host: configOrDefault("database", "host", ""),
            Port: configIntOrDefault("database", "port", 3306),
            Charset: "utf8mb4",
            User: configOrDefault("database", "user", ""),
            Pass: configOrDefault("database", "password", ""),
            DbName: configOrDefault("database", "dbname", ""),
            MaxIdleConns: configIntOrDefault("database", "max_idle_conns", 100),
            MaxOpenConns: configIntOrDefault("database", "max_open_conns", 200),
            ConnMaxLifeTime: configIntOrDefault("database", "conn_max_life_time", 500),
        },
        Log: &syncd.LogConfig{
            Path: configOrDefault("log", "path", "stdout"),
        },
        Syncd: &syncd.SyncdConfig{
            LocalSpace: configOrDefault("syncd", "local_space", "~/.syncd"),
            RemoteSpace: configOrDefault("syncd", "remote_space", "~/.syncd"),
            Cipher: configOrDefault("syncd", "cipher_key", ""),
            AppHost: configOrDefault("syncd", "app_host", ""),
        },
        Mail: &syncd.MailConfig{
            Enable: configIntOrDefault("mail", "enable", 0),
            Smtp: configOrDefault("mail", "smtp_host", ""),
            Port: configIntOrDefault("mail", "smtp_port", 465),
            User: configOrDefault("mail", "smtp_user", ""),
            Pass: configOrDefault("mail", "smtp_pass", ""),
        },
    }

    outputInfo("Log", cfg.Log.Path)
    outputInfo("Mail Enable", cfg.Mail.Enable)
    outputInfo("HTTP Service", cfg.Serve.Addr)

    if err := syncd.App.Init(cfg); err != nil {
        log.Fatal(err)
    }

    route.RegisterRoute()

    fmt.Println("Start Running...")
    if err := syncd.App.Start(); err != nil {
        log.Fatal(err)
    }
}
