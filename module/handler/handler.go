// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handler

import (
    "net/url"
    "net/http"
    "strings"

    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    "github.com/tinystack/goutil/gostring"
    "github.com/tinystack/goutil/goslice"
    "github.com/tinystack/goutil/goaes"
    userService "github.com/tinystack/syncd/service/user"
)

func setCrossDomainHeader(c *goweb.Context) {
    var origin string
    if referer := c.Request.Referer(); referer != "" {
        if u, err := url.Parse(referer); err == nil {
            origin = gostring.JoinStrings(u.Scheme, "://", u.Host)
        }
    }
    c.SetHeader("Access-Control-Allow-Origin", origin)
    c.SetHeader("Access-Control-Allow-Credentials", "true")
    c.SetHeader("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    c.SetHeader("Access-Control-Allow-Headers", "Origin, Content-Type, Cookie, Accept, X-Requested-With, Cache-Control, Authorization, X-Token")
}

func apiPrivCheck(c *goweb.Context) error {
    var err error
    var isLogin bool
    authToken, err := c.GetCookie("SYD_AUTH_TOKEN")

    c.Set("user_id", 0)
    c.Set("user_name", "")
    c.Set("email", "")

    // priv check
    loginReqPath := gostring.JoinSepStrings(" ", c.GetRequestMethod(), c.GetRequestPath())
    if err == nil && authToken != "" {
        authTokenBytes, err := gostring.Base64UrlDecode(authToken)
        tokenValBytes, err := goaes.Decrypt(syncd.CipherKey, authTokenBytes)
        if err == nil {
            tokenArr := strings.Split(string(tokenValBytes), "\t")
            if len(tokenArr) == 2 {
                userId := gostring.Str2Int(tokenArr[0])
                tokenKey := tokenArr[1]
                if userId > 0 && tokenKey != "" {
                    token := &userService.Token{
                        UserId: userId,
                        Token: tokenKey,
                    }
                    if status := token.ValidateToken(); status {
                        user := &userService.User{
                            ID: userId,
                        }
                        if err = user.Get(); err == nil {
                            if user.ID > 0 {
                                isLogin = true
                                c.Set("user_id", userId)
                                c.Set("user_name", user.Name)
                                c.Set("email", user.Email)
                                group := &userService.Group{
                                    ID: user.GroupId,
                                }
                                if err := group.Get(); err != nil {
                                    return syncd.RenderAppError(err.Error())
                                }
                                c.Set("priv", group.Priv)

                                // check priv
                                noCheckPrivList := []string{
                                    syncd.API_USER_LOGIN, syncd.API_USER_LOGOUT, syncd.API_USER_LOGIN_STATUS,
                                }
                                noCheckPriv := goslice.InSliceString(loginReqPath, noCheckPrivList)
                                if !noCheckPriv {
                                    havePriv := userService.CheckHavePriv(loginReqPath, group.Priv)
                                    if !havePriv {
                                        return syncd.RenderCustomerError(syncd.CODE_ERR_NO_PRIV, "no priv")
                                    }
                                }

                            }
                        }
                    }
                }
            }
        }
    }

    if !isLogin && loginReqPath != syncd.API_USER_LOGIN {
        return syncd.RenderCustomerError(syncd.CODE_ERR_NO_LOGIN, "no login")
    }

    return nil
}

func BeforeHandler(c *goweb.Context) error {
    setCrossDomainHeader(c)
    if err := apiPrivCheck(c); err != nil {
        return err
    }
    return nil
}

func AfterHandler(c *goweb.Context) error {
    return nil
}

func NotFoundHandler(c *goweb.Context) error {
    c.Json(http.StatusNotFound, goweb.JSON{
        "code": http.StatusNotFound,
        "message": "Page Not Found",
    })
    return nil
}

func ServerErrorHandler(err error, c *goweb.Context, code int) {
    errMsg := strings.Split(err.Error(), "=>")
    if len(errMsg) == 2 {
        code, msg := gostring.Str2Int(errMsg[0]), errMsg[1]
        c.Json(http.StatusOK, goweb.JSON{
            "code": code,
            "message": msg,
        })
    } else {
        c.Json(code, goweb.JSON{
            "code": code,
            "message": err.Error(),
        })
    }
}
