// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goweb

import (
    "time"
    "strings"
    "net"
    "net/url"
    "net/http"
    "io/ioutil"
    "strconv"
)

const multipartMemory = 32 << 20 // 32 MB

type Context struct {
    Serve           *Serve
    Request         *http.Request
    ResponseWriter  http.ResponseWriter
    Params          Params
    Keys            map[string]interface{}
}

func (c *Context) Reset(w http.ResponseWriter, r *http.Request) {
    c.Request = r
    c.ResponseWriter = w
    c.Params = c.Params[0:0]
    c.Keys = nil
}

func (c *Context) CloseCallback(cbFn func(), timeout int) {
    notify := c.ResponseWriter.(http.CloseNotifier).CloseNotify()
    go func() {
        select {
        case <-notify:
        case <-time.After(time.Second * time.Duration(timeout)):
        }
        cbFn()
    }()
}

//***  Request Data  ***//

func (c *Context) GetRequestPath() string {
    return c.Request.URL.Path
}

func (c *Context) GetRequestMethod() string {
    return c.Request.Method
}

func (c *Context) GetHeader(key string) string {
    return c.Request.Header.Get(key)
}

func (c *Context) SetHeader(key, value string) {
    if value == "" {
        c.ResponseWriter.Header().Del(key)
        return
    }
    c.ResponseWriter.Header().Set(key, value)
}

func (c *Context) ContentType() string {
    return filterFlags(c.GetHeader("Content-Type"))
}

func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
    if path == "" {
        path = "/"
    }
    http.SetCookie(c.ResponseWriter, &http.Cookie{
        Name:     name,
        Value:    url.QueryEscape(value),
        MaxAge:   maxAge,
        Path:     path,
        Domain:   domain,
        Secure:   secure,
        HttpOnly: httpOnly,
    })
}

func (c *Context) GetCookie(name string) (string, error) {
    cookie, err := c.Request.Cookie(name)
    if err != nil {
        return "", err
    }
    val, _ := url.QueryUnescape(cookie.Value)
    return val, nil
}

func (c *Context) SetStatusCode(code int) {
    c.ResponseWriter.WriteHeader(code)
}

func (c *Context) ClientIP() string {
    clientIP := c.GetHeader("X-Forwarded-For")
    clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
    if clientIP == "" {
        clientIP = strings.TrimSpace(c.GetHeader("X-Real-Ip"))
    }
    if clientIP != "" {
        return clientIP
    }
    if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
        return ip
    }
    return ""
}

//***  Keys  ***//

func (c *Context) Set(key string, value interface{}) {
    if c.Keys == nil {
        c.Keys = make(map[string]interface{})
    }
    c.Keys[key] = value
}

func (c *Context) Get(key string) (value interface{}, exists bool) {
    value, exists = c.Keys[key]
    return
}

func (c *Context) GetString(key string) (s string) {
    if val, ok := c.Get(key); ok && val != nil {
        s, _ = val.(string)
    }
    return
}

func (c *Context) GetBool(key string) (b bool) {
    if val, ok := c.Get(key); ok && val != nil {
        b, _ = val.(bool)
    }
    return
}

func (c *Context) GetInt(key string) (i int) {
    if val, ok := c.Get(key); ok && val != nil {
        i, _ = val.(int)
    }
    return
}

func (c *Context) GetInt64(key string) (i64 int64) {
    if val, ok := c.Get(key); ok && val != nil {
        i64, _ = val.(int64)
    }
    return
}

func (c *Context) GetFloat64(key string) (f64 float64) {
    if val, ok := c.Get(key); ok && val != nil {
        f64, _ = val.(float64)
    }
    return
}

func (c *Context) GetTime(key string) (t time.Time) {
    if val, ok := c.Get(key); ok && val != nil {
        t, _ = val.(time.Time)
    }
    return
}

func (c *Context) GetDuration(key string) (d time.Duration) {
    if val, ok := c.Get(key); ok && val != nil {
        d, _ = val.(time.Duration)
    }
    return
}

func (c *Context) GetStringSlice(key string) (ss []string) {
    if val, ok := c.Get(key); ok && val != nil {
        ss, _ = val.([]string)
    }
    return
}

func (c *Context) GetIntSlice(key string) (ss []int) {
    if val, ok := c.Get(key); ok && val != nil {
        ss, _ = val.([]int)
    }
    return
}

func (c *Context) GetStringMap(key string) (sm map[string]interface{}) {
    if val, ok := c.Get(key); ok && val != nil {
        sm, _ = val.(map[string]interface{})
    }
    return
}

func (c *Context) GetStringMapString(key string) (sms map[string]string) {
    if val, ok := c.Get(key); ok && val != nil {
        sms, _ = val.(map[string]string)
    }
    return
}

//***  Input Data  ***//

func (c *Context) Param(key string) (val string) {
    return c.Params.ByName(key)
}

func (c *Context) Query(key string) string {
    value, _ := c.GetQuery(key)
    return value
}

func (c *Context) QueryInt(key string) int {
    value, _ := strconv.Atoi(c.Query(key))
    return value
}

func (c *Context) QueryFloat64(key string) float64 {
    value, _ := strconv.ParseFloat(c.Query(key), 64)
    return value
}

func (c *Context) GetQuery(key string) (string, bool) {
    if values, ok := c.GetQueryArray(key); ok {
        return values[0], ok
    }
    return "", false
}

func (c *Context) GetQueryArray(key string) ([]string, bool) {
    if values, ok := c.Request.URL.Query()[key]; ok && len(values) > 0 {
        return values, true
    }
    return []string{}, false
}

func (c *Context) PostForm(key string) string {
    value, _ := c.GetPostForm(key)
    return value
}

func (c *Context) PostFormInt(key string) int {
    value, _ := strconv.Atoi(c.PostForm(key))
    return value
}

func (c *Context) PostFormFloat64(key string) float64 {
    value, _ := strconv.ParseFloat(c.PostForm(key), 64)
    return value
}

func (c *Context) GetPostForm(key string) (string, bool) {
    if values, ok := c.GetPostFormArray(key); ok {
        return values[0], ok
    }
    return "", false
}

func (c *Context) PostFormArray(key string) []string {
    values, _ := c.GetPostFormArray(key)
    return values
}

func (c *Context) GetPostFormArray(key string) ([]string, bool) {
    req := c.Request
    req.ParseMultipartForm(multipartMemory)
    if values := req.PostForm[key]; len(values) > 0 {
        return values, true
    }
    if req.MultipartForm != nil && req.MultipartForm.File != nil {
        if values := req.MultipartForm.Value[key]; len(values) > 0 {
            return values, true
        }
    }
    return []string{}, false
}

func (c *Context) GetRawData() ([]byte, error) {
    return ioutil.ReadAll(c.Request.Body)
}

//***  Render Data  ***//

func (c *Context) Json(code int, obj interface{}) error {
    c.SetHeader("Content-Type", "application/json; charset=utf-8")
    c.SetStatusCode(code)
    if err := renderJson(c.ResponseWriter, obj); err != nil {
        return err
    }
    return nil
}
