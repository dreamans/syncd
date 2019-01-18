// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handler

import (
    "net/http"
    "net/url"
    "strings"

    "github.com/dreamans/syncd"
    userService "github.com/dreamans/syncd/service/user"
    "github.com/tinystack/goutil/goaes"
    "github.com/tinystack/goutil/goslice"
    "github.com/tinystack/goutil/gostring"
    "github.com/tinystack/goweb"
)

func BeforeHandler(c *goweb.Context) error {
    resetReqParams(c)
    setCrossDomainHeader(c)
    if err := apiPrivCheck(c); err != nil {
        return err
    }
    return nil
}

func apiPrivCheck(c *goweb.Context) error {
    var err error
    authToken, _ := c.GetCookie("SYD_AUTH_TOKEN")

    c.Set("user_id", 0)
    c.Set("group_id", 0)
    c.Set("group_name", "")
    c.Set("user_name", "")
    c.Set("email", "")
    c.Set("mobile", "")
    c.Set("true_name", "")

    loginReqPath := gostring.JoinSepStrings(" ", c.GetRequestMethod(), c.GetRequestPath())

    if loginReqPath == syncd.API_USER_LOGIN {
        return nil
    }
    if authToken == "" {
        if loginReqPath == syncd.API_USER_LOGIN_STATUS {
            return nil
        }
        return syncd.RenderCustomerError(syncd.CODE_ERR_NO_LOGIN, "no login")
    }

    authTokenBytes, err := gostring.Base64UrlDecode(authToken)
    if err != nil {
        return syncd.RenderAppError("token check failed, " + err.Error())
    }
    tokenValBytes, err := goaes.Decrypt(syncd.CipherKey, authTokenBytes)
    if err != nil {
        return syncd.RenderAppError("token check failed, " + err.Error())
    }

    tokenArr := strings.Split(string(tokenValBytes), "\t")
    if len(tokenArr) != 2 {
        return syncd.RenderAppError("token check failed, len wrong")
    }

    token := &userService.Token{
        UserId: gostring.Str2Int(tokenArr[0]),
        Token: tokenArr[1],
    }
    if status := token.ValidateToken(); !status {
        return syncd.RenderCustomerError(syncd.CODE_ERR_NO_LOGIN, "token check failed, maybe your account is logged in on another device or token expired")
    }
    user := &userService.User{
        ID: token.UserId,
    }
    if err := user.Detail(); err != nil {
        return syncd.RenderAppError("token check failed, " + err.Error())
    }
    c.Set("user_id", user.ID)
    c.Set("user_name", user.Name)
    c.Set("email", user.Email)
    c.Set("group_id", user.GroupId)
    c.Set("mobile", user.Mobile)
    c.Set("true_name", user.TrueName)

    group := &userService.Group{
        ID: user.GroupId,
    }
    if err := group.Detail(); err != nil {
        return syncd.RenderAppError("token check failed, " + err.Error())
    }
    c.Set("priv", group.Priv)
    c.Set("group_name", group.Name)

    havePriv := userService.CheckHavePriv(loginReqPath, group.Priv)
    if !havePriv {
        noNeedPrivCheck := []string{
            syncd.API_USER_LOGIN_STATUS,
            syncd.API_USER_LOGOUT,
            syncd.API_USER_MY_UPDATE,
            syncd.API_USER_MY_PASSWORD,
            syncd.API_ROOT,
        }
        if in := goslice.InSliceString(loginReqPath, noNeedPrivCheck); !in {
            return syncd.RenderCustomerError(syncd.CODE_ERR_NO_PRIV, "no priv")
        }
    }
    go token.UpdateExpirationTime()
    return nil
}

func resetReqParams(c *goweb.Context) {
    limit := c.QueryInt("limit")
    if limit > 999 {
        limit = 999
    }
    if limit == 0 {
        limit = 10
    }
    c.Set("limit", limit)
}

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
