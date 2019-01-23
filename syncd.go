// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package syncd

import (
    "errors"
    "encoding/base64"
    "fmt"
    "os"
    "io"

    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/util/golog"
)

var (
    App	*syncd
)

const (
    Version = "v2.0.0"
)

func init() {
    App = newSyncd()
}

type syncd struct {
    Gin				*gin.Engine
    DB				*DB
    Logger			*golog.Logger
    LocalSpace		string
    RemoteSpace		string
    CipherKey		[]byte
    config			*Config
}

func newSyncd() *syncd {
    return &syncd{
        Gin: gin.New(),
    }
}

func (s *syncd) Init(cfg *Config) error {
    s.config = cfg

    if err := s.registerOrm(); err != nil {
        return err
    }

    s.registerLog()

    if err := s.initEnv(); err != nil {
        return err
    }

    return nil
}

func (s *syncd) Start() error {
    return s.Gin.Run(s.config.Serve.Addr)
}

func  (s *syncd) registerOrm() error {
    s.DB = NewDatabase(s.config.Db)
    return s.DB.Open()
}

func (s *syncd) registerLog() {
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
    s.Logger = golog.New(loggerHandler)
}

func (s *syncd) initEnv() error {
    s.LocalSpace = s.config.Syncd.LocalSpace
    s.RemoteSpace = s.config.Syncd.RemoteSpace
    if s.config.Syncd.Cipher == "" {
        return errors.New("syncd config 'Cipher' not setting")
    }
    dec, err := base64.StdEncoding.DecodeString(s.config.Syncd.Cipher)
    if err != nil {
        return errors.New(fmt.Sprintf("decode Cipher failed, %s", err.Error()))
    }
    s.CipherKey = dec

    return nil
}
