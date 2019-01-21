// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/tinystack/goutil/gostring"
    "github.com/tinystack/goweb"
    "github.com/dreamans/syncd"
    userService "github.com/dreamans/syncd/service/user"
)

func MyUpdate(c *goweb.Context) error {
    trueName, mobile, userId := c.PostForm("true_name"), c.PostForm("mobile"),  c.GetInt("user_id")

    user, err := userService.UserGetByPk(userId)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }

    user.TrueName = trueName
    user.Mobile = mobile
    user.DoNotUpdatePassword()
    if err := user.CreateOrUpdate(); err != nil {
        return syncd.RenderAppError(err.Error())
    }

    return syncd.RenderJson(c, nil)
}

func MyPasswordUpdate(c *goweb.Context) error {
    pass, newpass, userId := c.PostForm("password"), c.PostForm("newpassword"),  c.GetInt("user_id")

    user, err := userService.UserGetByPk(userId)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    password := gostring.StrMd5(gostring.JoinStrings(pass, user.Salt))
    if password != user.Password {
        return syncd.RenderCustomerError(syncd.CODE_ERR_USER_OR_PASS_WRONG, "current password incorrect")
    }

    user.Salt = gostring.StrRandom(10)
    user.Password = gostring.StrMd5(gostring.JoinStrings(newpass, user.Salt))
    if err := user.UpdatePassword(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}
