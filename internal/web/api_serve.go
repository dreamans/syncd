// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package web

import (
    "github.com/tinystack/syncd/pkg/webserve"
)

type ApiServeConfig struct {
    Addr    string
}

type ApiServe struct {
    serve   *webserve.Serve
}

var apiGetHandler map[string]func(c *webserve.Context) error = make(map[string]func(c *webserve.Context) error)

var apiPostHandler map[string]func(c *webserve.Context) error = make(map[string]func(c *webserve.Context) error)

func RunApiServe(cfg *ApiServeConfig) error {
    apiServe := NewApiServe(cfg)
    apiServe.Router()
    if err := apiServe.Run(); err != nil {
        return err
    }
    return nil
}

func NewApiServe(cfg *ApiServeConfig) *ApiServe {
    apiServe := &ApiServe{
        serve: webserve.New(),
    }
    if cfg.Addr != "" {
        apiServe.serve.Addr = cfg.Addr
    }

    apiServe.serve.BeforeHandler = beforeHandler
    apiServe.serve.AfterHandler = afterHandler
    apiServe.serve.ServerErrorHandler = serverErrorHandler
    apiServe.serve.NotFoundHandler = notFoundHandler
    apiServe.serve.MethodNotAllowHandler = notFoundHandler

    return apiServe
}

func (srv *ApiServe) Run() error {
    return srv.serve.Start()
}

func (srv *ApiServe) Router() {
    for path, handler := range apiGetHandler {
        srv.serve.GET(path, handler)
    }
    for path, handler := range apiPostHandler {
        srv.serve.POST(path, handler)
    }
}
