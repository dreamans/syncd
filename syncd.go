// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package syncd

import (
    "os"
    "time"

    "github.com/tinystack/goweb"
    "github.com/tinystack/golog"
    "github.com/tinystack/syncd/route"
    _ "github.com/tinystack/syncd/module"
)

type Syncd struct {
    config  *Config
    serve   *goweb.Serve
}

func NewSyncd(cfg *Config) *Syncd {
    var logger *golog.Logger
    if cfg.Serve.Log != "" {
        logger = golog.New(golog.NewFileHandler(cfg.Serve.Log))
    } else {
        logger = golog.New(os.Stderr)
    }
    syncd := &Syncd{
        config: cfg,
        serve: goweb.New(cfg.Serve.Addr),
    }
    syncd.serve.ReadTimeout = time.Second * time.Duration(cfg.Serve.ReadTimeout)
    syncd.serve.WriteTimeout = time.Second * time.Duration(cfg.Serve.WriteTimeout)
    syncd.serve.IdleTimeout = time.Second * time.Duration(cfg.Serve.IdleTimeout)
    syncd.serve.Logger = logger

    return syncd
}

func (s *Syncd) Start() error {
    s.regRoute()
    s.regDb()

    return s.serve.Start()
}

func (s *Syncd) regRoute() {
    s.serve.BeforeHandler = route.BeforeHandler
    s.serve.AfterHandler = route.AfterHandler
    s.serve.ServerErrorHandler = route.ServerErrorHandler
    s.serve.NotFoundHandler = route.NotFoundHandler
    s.serve.MethodNotAllowHandler = route.NotFoundHandler

    rg := route.RouteGroup()
    for _, r := range rg {
        switch r.Method {
        case "GET":
            s.serve.GET(r.Path, r.Handler)
        case "POST":
            s.serve.POST(r.Path, r.Handler)
        case "OPTIONS":
            s.serve.OPTIONS(r.Path, r.Handler)
        }
    }
}

func (s *Syncd) regDb() {

}
