// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package syncd

import (
    "time"
    "io"
    "os"
    "errors"
    "encoding/base64"

    "github.com/tinystack/goutil/gopath"
    "github.com/tinystack/goweb"
    "github.com/tinystack/golog"
    "github.com/jinzhu/gorm"
)

type Syncd struct {
    config  *Config
    serve   *goweb.Serve
}

var (
    Logger          *golog.Logger
    Orm             *gorm.DB
    DbInstance      *DB
    DataDir         string
    TmpDir          string
    CipherKey       []byte
)

func NewSyncd(cfg *Config) *Syncd {
    syncd := &Syncd{
        config: cfg,
        serve: goweb.New(cfg.Serve.Addr),
    }
    syncd.serve.ReadTimeout = time.Second * time.Duration(cfg.Serve.ReadTimeout)
    syncd.serve.WriteTimeout = time.Second * time.Duration(cfg.Serve.WriteTimeout)
    syncd.serve.IdleTimeout = time.Second * time.Duration(cfg.Serve.IdleTimeout)

    syncd.serve.BeforeHandler = beforeHandler
    syncd.serve.AfterHandler = afterHandler
    syncd.serve.ServerErrorHandler = serverErrorHandler
    syncd.serve.NotFoundHandler = notFoundHandler
    syncd.serve.MethodNotAllowHandler = notFoundHandler

    return syncd
}

func (s *Syncd) Start() error {
    return s.serve.Start()
}

func (s *Syncd) RegisterRoute(method, path string, handler goweb.HandlerFunc) {
    s.serve.Handler(method, path, handler)
}

func (s *Syncd) UnRegisterRoute() {}

func (s *Syncd) RegisterOrm() {
    DbInstance = NewDatabase(s.config.Db)
    if err := DbInstance.Open(); err != nil {
        panic(err)
    }
    Orm = DbInstance.DbHandler
}

func (s *Syncd) RegisterLog() {
    var loggerHandler io.Writer
    switch s.config.Log.Path {
    case "stdout":
        loggerHandler = os.Stdout
    case "stderr":
        loggerHandler = os.Stderr
    case "":
        loggerHandler = os.Stdout
    default:
        loggerHandler = golog.NewFileHandler(s.config.Log.Path)
    }
    Logger = golog.New(loggerHandler)
}

func (s *Syncd) InitEnv() {
    DataDir = s.config.Syncd.Dir
    if s.config.Syncd.Dir == "" {
        path, err := gopath.CurrentPath()
        if err != nil {
            panic(err)
        }
        DataDir = path + "/data"
    }
    if err := gopath.CreatePath(DataDir); err != nil {
        panic(err)
    }

    TmpDir = DataDir + "/tmp"
    if err := gopath.CreatePath(TmpDir); err != nil {
        panic(err)
    }

    if s.config.Syncd.Cipher == "" {
        panic(errors.New("syncd config 'Cipher' not setting"))
    }

    dec, err := base64.StdEncoding.DecodeString(s.config.Syncd.Cipher)
    if err != nil {
        panic(err)
    }
    CipherKey = dec
}
