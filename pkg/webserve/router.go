// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package webserve

type Route struct {
    Method string
    Path   string
    Name   string
}

type Router struct {
    serve  *Serve
    routes map[string]*Route
    trees  methodTrees
}

func NewRouter(ser *Serve) (r *Router) {
    r = &Router{
        serve: ser,
        routes: make(map[string]*Route),
        trees:  make(methodTrees, 0, 9),
    }
    return
}

func (r *Router) handle(method string, path string, handler HandlerFunc) {
    if path == "" {
        panic("andy: router path cannot be empty")
    }
    if path[0] != '/' {
        path = "/" + path
    }
    if path[len(path)-1] == '/' {
        path = path[:len(path)-1]
    }

    root := r.trees.get(method)
    if root == nil {
        root = &node{}
        r.trees = append(r.trees, methodTree{method: method, root: root})
    }
    root.addRoute(path, handler)

    r.routes[method+path] = &Route{
        Method: method,
        Path:   path,
    }
}
