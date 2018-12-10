// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
    "strings"
    "fmt"

    "github.com/tinystack/goweb"
)

type (
    Route struct {
        Method      string
        Path        string
        Handler     goweb.HandlerFunc
    }
)

var routeGroup []*Route

func Register(key string, handler goweb.HandlerFunc) {
    arrMap := strings.Split(key, " ")
    if len(arrMap) != 2 {
        panic(fmt.Sprintf("register router map failed, want \"Method /path\", have \"%v\"", key))
    }
    method := arrMap[0]
    path := arrMap[1]

    routeGroup = append(routeGroup, &Route{
        Method: method,
        Path: path,
        Handler: handler,
    })
}

func RouteGroup() []*Route {
    return routeGroup
}
