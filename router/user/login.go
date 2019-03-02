// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/module/user"
    "github.com/dreamans/syncd/util/gois"

)

type LoginBind struct {
    Username    string  `form:"username" binding:"required"`
    Password    string  `form:"password" binding:"required"`
}

func Logout(c *gin.Context) {
    login := &user.Login{
        UserId: c.GetInt("user_id"),
    }
    if err := login.Logout(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.Success(c)
}

func Login(c *gin.Context) {
    var form LoginBind
    if err := c.ShouldBind(&form); err != nil {
        render.ParamError(c, err.Error())
        return
    }

    login := &user.Login{
        Password: form.Password,
    }
    if gois.IsEmail(form.Username) {
        login.Email = form.Username
    } else {
        login.Username = form.Username
    }

    if err := login.Login(); err != nil {
        render.CustomerError(c, render.CODE_ERR_LOGIN_FAILED , err.Error())
        return
    }

    userInfo := map[string]interface{}{
        "token": login.Token,
    }

    render.JSON(c, userInfo)
}

func LoginStatus(c *gin.Context) {
    privilege, _ := c.Get("privilege")
    render.JSON(c, gin.H{
        "is_login": 1,
        "user_id": c.GetInt("user_id"),
        "username": c.GetString("username"),
        "email": c.GetString("email"),
        "truename": c.GetString("truename"),
        "mobile": c.GetString("mobile"),
        "role_name": c.GetString("role_name"),
        "privilege": privilege,
    })
}
