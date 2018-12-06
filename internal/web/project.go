// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package web

import (
    "github.com/tinystack/syncd/pkg/webserve"
)

func init() {
    apiGetHandler["/project/new"] = func (c *webserve.Context) error {
        c.Json(200, webserve.JSON{
            "code": "0",
            "message": "Hijacker",
        })

        return nil
    }
}
