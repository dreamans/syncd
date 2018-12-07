// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package webserve

import (
    "log"
    "net/http"
    "sync"
)

type (
    Serve struct {
        Addr                  string
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

func New() (ser *Serve) {
    ser = &Serve{
        Addr: "localhost:8068",
    }
    ser.pool.New = func() interface{} {
        return ser.NewContext()
    }
    ser.router = NewRouter(ser)
    ser.ServerErrorHandler = ser.DefaultServerErrorHandler
    ser.NotFoundHandler = ser.DefaultNotFoundHandler
    ser.MethodNotAllowHandler = ser.DefaulMethodNotAllowHandler
    return
}

func (ser *Serve) NewContext() *Context {
    return &Context{}
}

func (ser *Serve) Routes() map[string]*Route {
    return ser.router.routes
}

func (ser *Serve) Start() (err error) {
    err = http.ListenAndServe(ser.Addr, ser)
    return
}

func (ser *Serve) DefaultServerErrorHandler(err error, c *Context, code int) {
    log.Println(err)
}

func (ser *Serve) DefaultNotFoundHandler(c *Context) (err error) {
    log.Println("Not Found")
    return
}

func (ser *Serve) DefaulMethodNotAllowHandler(c *Context) (err error) {
    log.Println("Not Allowed")
    return
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

func (ser *Serve) Handle(method string, path string, handler HandlerFunc) {
    ser.router.handle(method, path, handler)
}

func (ser *Serve) GET(path string, handler HandlerFunc) {
    ser.Handle("GET", path, handler)
}

func (ser *Serve) POST(path string, handler HandlerFunc) {
    ser.Handle("POST", path, handler)
}
