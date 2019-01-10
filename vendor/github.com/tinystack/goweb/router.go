// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goweb

type Route struct {
    Method string
    Path   string
}

type Router struct {
    serve   *Serve
    routes  map[string]*Route
    trees   methodTrees
    methods map[string]HandlerFunc
}

func NewRouter(ser *Serve) (r *Router) {
    r = &Router{
        serve: ser,
        routes: make(map[string]*Route),
        trees:  make(methodTrees, 0, 9),
    }
    return
}

func (r *Router) methodHandler(method string, handler HandlerFunc) {
    r.methods[method] = handler
    r.routes[method] = &Route{
        Method: method,
        Path: "*",
    }
}

func (r *Router) handler(method string, path string, handler HandlerFunc) {
    if path == "" {
        panic("router path cannot be empty")
    }
    if path[0] != '/' {
        path = "/" + path
    }
    if len(path) > 1 && path[len(path)-1] == '/' {
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
