// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package web

import (
    "net/http"
)

type Context struct {
    Serve          *Serve
    Request        *http.Request
    ResponseWriter http.ResponseWriter
    Params         Params
}

func (c *Context) Reset(w http.ResponseWriter, r *http.Request) {
    c.Request = r
    c.ResponseWriter = w
}

func (c *Context) Method() string {
    return c.Request.Method
}

func (c *Context) Param(key string) (val string) {
    return c.Params.ByName(key)
}
