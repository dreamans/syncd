// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/govalidate"
    "github.com/tinystack/goutil/gois"
    "github.com/tinystack/syncd"
    userService "github.com/tinystack/syncd/service/user"
)

type UserParamValid struct {
    GroupId     int     `valid:"required" errmsg:"required=user group can not be empty"`
    Name        string  `valid:"required" errmsg:"required=user name group can not be empty"`
    Email       string  `valid:"required|email" errmsg:"required=user email can not be empty|email=email format incorrect"`
    Mobile      string  `valid:"mobile" errmsg:"mobile=mobile format incorrect"`
}

func UserUpdate(c *goweb.Context) error {
    params := UserParamValid{
        GroupId: c.PostFormInt("group_id"),
        Name: c.PostForm("name"),
        Email: c.PostForm("email"),
        Mobile: c.PostForm("mobile"),
    }
    if valid := govalidate.NewValidate(&params); !valid.Pass() {
        return syncd.RenderParamError(valid.LastFailed().Msg)
    }
    id := c.PostFormInt("id")
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
    offset, limit, keyword := c.QueryInt("offset"), c.QueryInt("limit"), c.Query("keyword")
    user := &userService.User{}
    list, total, err := user.List(keyword, offset, limit)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    var groupIdList []int
    for _, l := range list {
        groupIdList = append(groupIdList, l.GroupId)
    }
    if len(groupIdList) > 0 {
        group := &userService.Group{}
        groupNameList, err := group.GetNameByIds(groupIdList)
        if err != nil {
            syncd.RenderAppError(err.Error())
            return nil
        }
        for k, v := range list {
            if groupName, exists := groupNameList[v.GroupId]; exists {
                list[k].GroupName = groupName
            }
        }
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
    if err := user.Get(); err != nil {
        return syncd.RenderAppError(err.Error())
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
    checkType := c.Query("type")
    keyword := c.Query("keyword")
    id := c.QueryInt("id")

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

func UserSearch(c *goweb.Context) error {
    keyword := c.Query("keyword")
    if keyword == "" {
        return syncd.RenderJson(c, nil)
    }
    user := &userService.User{}
    if gois.IsEmail(keyword) {
        user.Email = keyword
    } else {
        user.Name = keyword
    }
    list, err := user.Search()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }

    var groupIds []int
    for _, l := range list {
        groupIds = append(groupIds, l.GroupId)
    }
    if len(groupIds) > 0 {
        group := &userService.Group{}
        groupNameList, err := group.GetNameByIds(groupIds)
        if err != nil {
            return syncd.RenderAppError(err.Error())
        }
        for k, l := range list {
            val, key := groupNameList[l.GroupId]
            if key {
                list[k].GroupName = val
            }
        }
    }

    return syncd.RenderJson(c, goweb.JSON{
        "list": list,
    })
}
