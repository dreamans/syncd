// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/module/user"
    "github.com/dreamans/syncd/util/gostring"
)

type UserSettingForm struct {
    Truename    string  `form:"truename"`
    Mobile      string  `form:"mobile"`
}

type UserPasswordForm struct {
    Password         string  `form:"password"`
    NewPassword      string  `form:"new_password"`
}

func MyUserSetting(c *gin.Context) {
    var form UserSettingForm
    if err := c.ShouldBind(&form); err != nil {
        render.ParamError(c, err.Error())
        return
    }

    u := &user.User{
        ID: c.GetInt("user_id"),
        Truename: form.Truename,
        Mobile: form.Mobile,
    }
    if err := u.UserSettingUpdate(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.Success(c)
}

func MyUserPassword(c *gin.Context) {
    var form UserPasswordForm
    if err := c.ShouldBind(&form); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    user := &user.User{
        ID: c.GetInt("user_id"),
    }
    if err := user.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    if user.Password != gostring.StrMd5(gostring.JoinStrings(form.Password, user.Salt)) {
        render.CustomerError(c, render.CODE_ERR_USER_OR_PASS_WRONG, "current password wrong")
        return
    }

    user.Password = form.NewPassword
    if err := user.UpdatePassword(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    render.Success(c)
}
