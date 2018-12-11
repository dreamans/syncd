// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package syncd

import (
    "net/http"

    "github.com/tinystack/goweb"
)

func RenderParamError(c *goweb.Context, msg string) {
    c.Json(http.StatusOK, goweb.JSON{
        "code": CODE_ERR_PARAM,
        "message": msg,
    })
}

func RenderJson(c *goweb.Context, data interface{}) {
    c.Json(http.StatusOK, goweb.JSON{
        "code": CODE_OK,
        "message": "success",
        "data": data,
    })
}
