// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
    "net/url"
    "net/http"

    "github.com/tinystack/goweb"
    "github.com/tinystack/goutil"
)

func BeforeHandler(c *goweb.Context) error {
    var origin string
    if referer := c.Request.Referer(); referer != "" {
        if u, err := url.Parse(referer); err == nil {
            origin = goutil.JoinStrings(u.Scheme, "://", u.Host)
        }
    }
    c.SetHeader("Access-Control-Allow-Origin", origin)
    c.SetHeader("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    c.SetHeader("Access-Control-Allow-Headers", "Origin, Content-Type, Cookie, Accept, X-Requested-With, Cache-Control, Authorization, X-Token")

    return nil
}

func AfterHandler(c *goweb.Context) error {
    return nil
}

func ServerErrorHandler(error error, c *goweb.Context, code int) {
    c.Logger().Error("server error occurs, code[%v], error[%v]", code, error)
    c.Json(http.StatusInternalServerError, goweb.JSON{
        "code": http.StatusInternalServerError,
        "message": "Internal Server Error",
    })
}

func NotFoundHandler(c *goweb.Context) error {
    c.Json(http.StatusNotFound, goweb.JSON{
        "code": http.StatusNotFound,
        "message": "Page Not Found",
    })
    return nil
}
