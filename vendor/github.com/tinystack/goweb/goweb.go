// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goweb

import (
    "net/http"
    "sync"
    "time"
)

type (
    Serve struct {
        Addr                  string
        ReadTimeout           time.Duration
        WriteTimeout          time.Duration
        IdleTimeout           time.Duration
        ServerErrorHandler    func(error, *Context, int)
        NotFoundHandler       HandlerFunc
        MethodNotAllowHandler HandlerFunc
        BeforeHandler         HandlerFunc
        AfterHandler          HandlerFunc
        pool                  sync.Pool
        router                *Router
    }
)

type HandlerFunc func(*Context) error

func New(addr string) (ser *Serve) {
    ser = &Serve{
        Addr: addr,
    }
    ser.pool.New = func() interface{} {
        return ser.NewContext()
    }
    ser.router = NewRouter(ser)
    ser.ServerErrorHandler = ser.EmptyServeErrorHandler
    ser.NotFoundHandler = ser.EmptyNotFoundHandler
    ser.MethodNotAllowHandler = ser.EmptyNotFoundHandler
    return
}

func (ser *Serve) NewContext() *Context {
    return &Context{
        Serve: ser,
    }
}

func (ser *Serve) Routes() map[string]*Route {
    return ser.router.routes
}

func (ser *Serve) Start() (err error) {
    server := &http.Server{
        Addr: ser.Addr,
        Handler: ser,
        ReadTimeout: ser.ReadTimeout,
        WriteTimeout: ser.WriteTimeout,
        IdleTimeout: ser.IdleTimeout,
    }
    err = server.ListenAndServe()
    return
}

func (ser *Serve) EmptyServeErrorHandler(err error, c *Context, code int) {

}

func (ser *Serve) EmptyNotFoundHandler(c *Context) (err error) {
    return nil
}

func (ser *Serve) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    c := ser.pool.Get().(*Context)
    c.Reset(w, r)
    ser.requestHTTPHandle(c)
    ser.pool.Put(c)
}

func (ser *Serve) requestHTTPHandle(c *Context) {
    httpMethod := c.Request.Method
    path := c.Request.URL.Path

    if mh, ok := ser.router.methods[httpMethod]; ok {
        if err := mh(c); err != nil {
            ser.ServerErrorHandler(err, c, http.StatusInternalServerError)
            return
        }
    }

    if l := len(path); l > 1 && path[l-1] == '/' {
        path = path[:l-1]
    }
    t := ser.router.trees
    for i, tl := 0, len(t); i < tl; i++ {
        if t[i].method != httpMethod {
            continue
        }
        root := t[i].root
        handler, params, _ := root.getValue(path)
        if handler != nil {
            c.Params = params
            handlers := []HandlerFunc{}
            if ser.BeforeHandler != nil {
                handlers = append(handlers, ser.BeforeHandler)
            }
            handlers = append(handlers, handler)
            if ser.AfterHandler != nil {
                handlers = append(handlers, ser.AfterHandler)
            }
            for _, h := range handlers {
                if err := h(c); err != nil {
                    ser.ServerErrorHandler(err, c, http.StatusInternalServerError)
                    return
                }
            }
            return
        }
        break
    }
    if ser.MethodNotAllowHandler != nil {
        for _, tree := range ser.router.trees {
            if tree.method == httpMethod {
                continue
            }
            if handler, _, _ := tree.root.getValue(path); handler != nil {
                ser.MethodNotAllowHandler(c)
                return
            }
        }
    }

    if ser.NotFoundHandler != nil {
        ser.NotFoundHandler(c)
    }

    return
}

func (ser *Serve) Handler(method string, path string, handler HandlerFunc) {
    if method == "*" {
        ser.router.methodHandler(path, handler)
    } else {
        ser.router.handler(method, path, handler)
    }
}

func (ser *Serve) GET(path string, handler HandlerFunc) {
    ser.Handler(http.MethodGet, path, handler)
}

func (ser *Serve) POST(path string, handler HandlerFunc) {
    ser.Handler(http.MethodPost, path, handler)
}

func (ser *Serve) HEAD(path string, handler HandlerFunc) {
    ser.Handler(http.MethodHead, path, handler)
}

func (ser *Serve) OPTIONS(path string, handler HandlerFunc) {
    ser.Handler(http.MethodOptions, path, handler)
}

func (ser *Serve) PUT(path string, handler HandlerFunc) {
    ser.Handler(http.MethodPut, path, handler)
}

func (ser *Serve) DELETE(path string, handler HandlerFunc) {
    ser.Handler(http.MethodDelete, path, handler)
}

func (ser *Serve) CONNECT(path string, handler HandlerFunc) {
    ser.Handler(http.MethodConnect, path, handler)
}

func (ser *Serve) TRACE(path string, handler HandlerFunc) {
    ser.Handler(http.MethodTrace, path, handler)
}

func (ser *Serve) PATCH(path string, handler HandlerFunc) {
    ser.Handler(http.MethodPatch, path, handler)
}
