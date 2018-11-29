// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
    "log"

    "github.com/tinystack/syncd/internal/web"
)

func main() {
    cfg := &web.ApiServeConfig{
        Addr: "localhost:8068",
    }
    if err := web.RunApiServe(cfg); err != nil {
        log.Fatalln(err.Error())
    }
}

