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

func RunApiServe(cfg *ApiServeConfig) error {
    apiServe := webserve.New()
    if cfg.Addr != "" {
        apiServe.Addr = cfg.Addr
    }

    return nil
}
