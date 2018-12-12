// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package syncd

import (
    "net/http"

    "github.com/tinystack/goweb"
)

const (
    CODE_OK         = 0
    CODE_ERR_SYSTEM = 1000
    CODE_ERR_APP    = 1001
    CODE_ERR_PARAM  = 1002
)

func RenderParamError(c *goweb.Context, msg string) {
    c.Json(http.StatusOK, goweb.JSON{
        "code": CODE_ERR_PARAM,
        "message": msg,
    })
}

func RenderAppError(c *goweb.Context, msg string) {
    c.Json(http.StatusOK, goweb.JSON{
        "code": CODE_ERR_APP,
        "message": "success",
    })
}

func RenderJson(c *goweb.Context, data interface{}) {
    c.Json(http.StatusOK, goweb.JSON{
        "code": CODE_OK,
        "message": "success",
        "data": data,
    })
}
