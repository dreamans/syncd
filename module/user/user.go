// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/govalidate"
    "github.com/tinystack/syncd"
    userService "github.com/tinystack/syncd/service/user"
)

type UserParamValid struct {
    GroupId     int     `valid:"required" errmsg:"required=user group can not be empty"`
    Name        string  `valid:"required" errmsg:"required=user name group can not be empty"`
    Email       string  `valid:"required|email" errmsg:"required=user email can not be empty|email=email format incorrect"`
    Mobile      string  `valid:"mobile" errmsg:"mobile=mobile format incorrect"`
}

func UserNew(c *goweb.Context) error {
    return userUpdate(c, 0)
}

func UserEdit(c *goweb.Context) error {
    id := c.PostFormInt("id")
    if id == 0 {
     return syncd.RenderParamError("id can not empty")
    }
    return userUpdate(c, id)
}

func userUpdate(c *goweb.Context, id int) error {
    params := UserParamValid{
        GroupId: c.PostFormInt("group_id"),
        Name: c.PostForm("name"),
        Email: c.PostForm("email"),
        Mobile: c.PostForm("mobile"),
    }
    if valid := govalidate.NewValidate(&params); !valid.Pass() {
        return syncd.RenderParamError(valid.LastFailed().Msg)
    }
    password := c.PostForm("password")
    if id == 0 {
        if password == "" || len(password) != 32 {
            return syncd.RenderParamError("password incorrect")
        }
    }
    var (
        existsUser *userService.User
        exists bool
        err error
    )
    existsUser = &userService.User{
        ID: id,
        Name: params.Name,
    }
    exists, err = existsUser.CheckUserExists()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if exists {
        return syncd.RenderCustomerError(syncd.CODE_ERR_DATA_REPEAT, "user name exists")
    }

    existsUser = &userService.User{
        ID: id,
        Email: params.Email,
    }
    exists, err = existsUser.CheckUserExists()
    if err != nil {
        syncd.RenderAppError(err.Error())
        return nil
    }
    if exists {
        return syncd.RenderCustomerError(syncd.CODE_ERR_DATA_REPEAT, "email exists")
    }

    user := &userService.User{
        ID: id,
        GroupId: params.GroupId,
        Name: params.Name,
        Password: password,
        Email: params.Email,
        Mobile: params.Mobile,
        TrueName: c.PostForm("true_name"),
        LockStatus: c.PostFormInt("lock_status"),
    }
    if err := user.CreateOrUpdate(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

func UserList(c *goweb.Context) error {
    offset, limit, keyword := c.QueryInt("offset"), c.GetInt("limit"), c.Query("keyword")
    user := &userService.User{}
    list, total, err := user.List(keyword, offset, limit)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }

    list, err = userService.GroupUserListFillGroupName(list)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }

    return syncd.RenderJson(c, goweb.JSON{
        "list": list,
        "total": total,
    })
}

func UserDetail(c *goweb.Context) error {
    user := &userService.User{
        ID: c.QueryInt("id"),
    }
    if err := user.Detail(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if user.ID == 0 {
        return syncd.RenderAppError("user not exists")
    }
    return syncd.RenderJson(c, goweb.JSON{
        "id": user.ID,
        "group_id": user.GroupId,
        "name": user.Name,
        "email": user.Email,
        "true_name": user.TrueName,
        "mobile": user.Mobile,
        "lock_status": user.LockStatus,
    })
}

func UserExists(c *goweb.Context) error {
    checkType, keyword, id := c.Query("type"), c.Query("keyword"), c.QueryInt("id")
    user := &userService.User{
        ID: id,
    }
    switch checkType {
    case "name":
        user.Name = keyword
    case "email":
        user.Email = keyword
    }
    exists, err := user.CheckUserExists()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "exists": exists,
    })
}

func UserDelete(c *goweb.Context) error {
    user := &userService.User{
        ID: c.PostFormInt("id"),
    }
    if err := user.Delete(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

