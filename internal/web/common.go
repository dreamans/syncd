// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package web

import (
    "log"

    "github.com/tinystack/syncd/pkg/webserve"
)

func beforeHandler(c *webserve.Context) error {

    return nil
}

func afterHandler(c *webserve.Context) error {
    log.Println("access " + c.Method())

    return nil
}

func serverErrorHandler(error error, c *webserve.Context, code int) {

}

func notFoundHandler(c *webserve.Context) error {

    return nil
}
