// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/govalidate"
    "github.com/tinystack/syncd"
    "github.com/tinystack/syncd/route"
    userService "github.com/tinystack/syncd/service/user"
)

func init() {
    route.Register(route.API_USER_UPDATE, updateUser)
    route.Register(route.API_USER_LIST, listUser)
}

type UserParamValid struct {
    GroupId     int     `valid:"required" errmsg:"required=user group can not be empty"`
    Name        string  `valid:"required" errmsg:"required=user name group can not be empty"`
    Password    string  `valid:"required|str_range=6,20" errmsg:"required=password group can not be empty|str_range=password must between 6 and 20 characters"`
    Email       string  `valid:"required|email" errmsg:"required=user email can not be empty|email=email format incorrect"`
    Mobile      string  `valid:"mobile" errmsg:"mobile=mobile format incorrect"`
}

func updateUser(c *goweb.Context) error {
    params := UserParamValid{
        GroupId: c.PostFormInt("group_id"),
        Name: c.PostForm("name"),
        Email: c.PostForm("email"),
        Mobile: c.PostForm("mobile"),
    }
    if valid := govalidate.NewValidate(&params); !valid.Pass() {
        syncd.RenderParamError(c, valid.LastFailed().Msg)
        return nil
    }
    id := c.PostFormInt("id")
    user := &userService.User{
        ID: id,
        GroupId: params.GroupId,
        Name: params.Name,
        Password: params.Password,
        Email: params.Email,
        Mobile: params.Mobile,
        TrueName: c.PostForm("true_name"),
    }
    if err := user.CreateOrUpdate(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}

func listUser(c *goweb.Context) error {
    offset, limit, keyword := c.QueryInt("offset"), c.QueryInt("limit"), c.Query("keyword")
    user := &userService.User{}
    list, total, err := user.List(keyword, offset, limit)
    if err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, goweb.JSON{
        "list": list,
        "total": total,
    })
    return nil
}
