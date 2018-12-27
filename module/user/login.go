// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    "github.com/tinystack/syncd/route"
    userService "github.com/tinystack/syncd/service/user"
)

func init() {
    route.Register(route.API_USER_LOGIN, loginUser)
}

func loginUser(c *goweb.Context) error {
    name, pass := c.PostForm("name"), c.PostForm("pass")

    if name == "" || pass == "" {
        return syncd.RenderParamError(c, "username or password name can not empty")
    }
    login := &userService.Login{
        Name: name,
        Pass: pass,
    }
    if err := login.Login(); err != nil {
        return syncd.RenderCustomerError(c, syncd.CODE_ERR_LOGIN_FAILED, err.Error())
    }

    userDetail := login.GetUserDetail()

    return syncd.RenderJson(c, goweb.JSON{
        "user_id": userDetail.ID,
        "name": userDetail.Name,
        "email": userDetail.Email,
        "token": login.GetToken(),
    })
}
