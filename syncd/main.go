// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
    "log"
    "flag"
    "fmt"
    "os"

    "github.com/Unknwon/goconfig"
    "github.com/dreamans/syncd"
    "github.com/dreamans/syncd/route"
    handlerModule "github.com/dreamans/syncd/module/handler"
)

var (
    help bool
    syncdIni string
    version bool
    configFile *goconfig.ConfigFile
)

func init() {
    flag.BoolVar(&help, "h", false, "This help")
    flag.StringVar(&syncdIni, "c", "./etc/syncd.ini", "Set configuration file `file`")
    flag.BoolVar(&version, "v", false, "Version number")

    flag.Usage = usage
    flag.Parse()
}

func usage() {
    fmt.Printf("Usage: syncd [-c filename]\n\nOptions:\n")
    flag.PrintDefaults()
}

func initCfg() {
    var err error
    configFile, err = goconfig.LoadConfigFile(syncdIni)
    if err != nil {
        fmt.Printf("can not load config file '%s', %s\n", syncdIni, err.Error())
        os.Exit(1)
    }
    outputInfo("Config Loaded", syncdIni)
}

func configOrDefault(section, key, useDefault string) string {
    val, err := configFile.GetValue(section, key)
    if err != nil {
        return useDefault
    }
    return val
}

func configIntOrDefault(section, key string, useDefault int) int {
    val, err := configFile.Int(section, key)
    if err != nil {
        return useDefault
    }
    return val
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
    outputInfo("Version", syncd.VERSION)
}

func outputInfo(tag string, value interface{}) {
    fmt.Printf("%-18s    %v\n", tag + ":", value)
}

func run() {
    if help {
        flag.Usage()
        os.Exit(0)
    }
    if version {
        fmt.Printf("syncd/%s\n", syncd.VERSION)
        os.Exit(0)
    }

    welcome()

    initCfg()

    cfg := &syncd.Config{
        Serve: &syncd.ServeConfig{
            Addr: configOrDefault("serve", "addr", ":8868"),
            ReadTimeout: configIntOrDefault("serve", "read_timeout", 300),
            WriteTimeout: configIntOrDefault("serve", "write_timeout", 300),
            IdleTimeout: configIntOrDefault("serve", "idle_timeout", 300),
        },
        Db: &syncd.DbConfig{
            Host: configOrDefault("database", "host", "127.0.0.1"),
            Port: configOrDefault("database", "port", "3306"),
            Charset: "utf8mb4",
            User: configOrDefault("database", "user", ""),
            Pass: configOrDefault("database", "password", ""),
            DbName: configOrDefault("database", "dbname", ""),
            TablePrefix: "syd_",
        },
        Log: &syncd.LogConfig{
            Path: configOrDefault("log", "path", "stdout"),
        },
        Syncd: &syncd.SyncdConfig{
            Dir: configOrDefault("syncd", "workspace", "/tmp/.syncd"),
            Cipher: configOrDefault("syncd", "cipher_key", "pG1L62EM0cPIIOwusQsbcV8Cs6j/M1RxLoXIZylWUC4="),
        },
        Mail: &syncd.MailConfig{
            Enable: configIntOrDefault("mail", "enable", 0),
            Smtp: configOrDefault("mail", "smtp_host", ""),
            Port: configIntOrDefault("mail", "smtp_port", 465),
            User: configOrDefault("mail", "smtp_user", ""),
            Pass: configOrDefault("mail", "smtp_pass", ""),
        },
    }

    syd := syncd.NewSyncd(cfg)

    syd.InitEnv()
    outputInfo("Workspace", syncd.DataDir)

    syd.RegisterLog()
    outputInfo("Log", cfg.Log.Path)

    syd.RegisterOrm()
    outputInfo("Database", cfg.Db.Host)

    syd.RegisterMail()
    outputInfo("Mail Enable", cfg.Mail.Enable)

    routes := route.RouteGroup()
    for _, r := range routes {
        syd.RegisterRoute(r.Method, r.Path, r.Handler)
    }
    syd.RegisterServeHandler(syncd.ServeHandler{
        BeforeHandler: handlerModule.BeforeHandler,
        AfterHandler: handlerModule.AfterHandler,
        ServerErrorHandler: handlerModule.ServerErrorHandler,
        NotFoundHandler: handlerModule.NotFoundHandler,
        MethodNotAllowHandler: handlerModule.NotFoundHandler,
    })

    outputInfo("HTTP Service", cfg.Serve.Addr)

    fmt.Println("Start Running...")
    if err := syd.Start(); err != nil {
        log.Fatalln(err.Error())
    }
}

func main() {
    run()
}
