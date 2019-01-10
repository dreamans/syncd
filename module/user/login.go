// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/tinystack/goutil/gois"
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    userService "github.com/tinystack/syncd/service/user"
)

func Login(c *goweb.Context) error {
    name, pass := c.PostForm("name"), c.PostForm("pass")
    if name == "" || pass == "" {
        return syncd.RenderParamError("username or password name can not empty")
    }
    login := &userService.Login{
        Pass: pass,
        Ip: c.ClientIP(),
    }
    if gois.IsEmail(name) {
        login.Email = name
    } else {
        login.Name = name
    }

    if err := login.Login(); err != nil {
        return syncd.RenderCustomerError(syncd.CODE_ERR_LOGIN_FAILED, err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "user_id": login.UserDetail.ID,
        "name": login.UserDetail.Name,
        "email": login.UserDetail.Email,
        "token": login.Token,
    })
}

func LoginStatus(c *goweb.Context) error {
    userId := c.GetInt("user_id")
    return syncd.RenderJson(c, goweb.JSON{
        "is_login": userId > 0,
        "group_id": c.GetInt("group_id"),
        "group_name": c.GetString("group_name"),
        "user_id": userId,
        "name": c.GetString("user_name"),
        "email": c.GetString("email"),
        "mobile": c.GetString("mobile"),
        "true_name": c.GetString("true_name"),
        "priv": c.GetIntSlice("priv"),
    })
}

func Logout(c *goweb.Context) error {
    userId := c.GetInt("user_id")
    token := &userService.Token{
        UserId: userId,
    }
    if err := token.DeleteByUserId(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

