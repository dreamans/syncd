// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/govalidate"
    "github.com/tinystack/goutil"
    "github.com/tinystack/syncd"
    "github.com/tinystack/syncd/route"
    userService "github.com/tinystack/syncd/service/user"
)

func init() {
    route.Register(route.API_USER_UPDATE, updateUser)
    route.Register(route.API_USER_LIST, listUser)
    route.Register(route.API_USER_DETAIL, detailUser)
    route.Register(route.API_USER_EXISTS, existsUser)
    route.Register(route.API_USER_DELETE, deleteUser)
    route.Register(route.API_USER_SEARCH, searchUser)
}

type UserParamValid struct {
    GroupId     int     `valid:"required" errmsg:"required=user group can not be empty"`
    Name        string  `valid:"required" errmsg:"required=user name group can not be empty"`
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
    password := c.PostForm("password")
    if id == 0 {
        if password == "" || len(password) != 32 {
            syncd.RenderParamError(c, "password incorrect")
            return nil
        }
    }

    // check name is exists
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
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    if exists {
        syncd.RenderAppError(c, "user name exists")
        return nil
    }
    // check email is exists
    existsUser = &userService.User{
        ID: id,
        Email: params.Email,
    }
    exists, err = existsUser.CheckUserExists()
    if err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    if exists {
        syncd.RenderAppError(c, "email exists")
        return nil
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
    var groupIdList []int
    for _, l := range list {
        groupIdList = append(groupIdList, l.GroupId)
    }
    if len(groupIdList) > 0 {
        group := &userService.Group{}
        groupNameList, err := group.GetNameByIds(groupIdList)
        if err != nil {
            syncd.RenderAppError(c, err.Error())
            return nil
        }
        for k, v := range list {
            if groupName, exists := groupNameList[v.GroupId]; exists {
                list[k].GroupName = groupName
            }
        }
    }

    syncd.RenderJson(c, goweb.JSON{
        "list": list,
        "total": total,
    })
    return nil
}

func detailUser(c *goweb.Context) error {
    user := &userService.User{
        ID: c.QueryInt("id"),
    }
    if err := user.Get(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, user)
    return nil
}

func existsUser(c *goweb.Context) error {
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
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, goweb.JSON{
        "exists": exists,
    })
    return nil
}

func deleteUser(c *goweb.Context) error {
    user := &userService.User{
        ID: c.PostFormInt("id"),
    }
    if err := user.Delete(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}

func searchUser(c *goweb.Context) error {
    keyword := c.Query("keyword")
    if keyword == "" {
        syncd.RenderJson(c, nil)
        return nil
    }

    user := &userService.User{}
    if goutil.IsEmail(keyword) {
        user.Email = keyword
    } else {
        user.Name = keyword
    }
    list, err := user.Search()
    if err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }

    //append group_name
    var groupIds []int
    for _, l := range list {
        groupIds = append(groupIds, l.GroupId)
    }
    if len(groupIds) > 0 {
        group := &userService.Group{}
        groupNameList, err := group.GetNameByIds(groupIds)
        if err != nil {
            syncd.RenderAppError(c, err.Error())
            return nil
        }
        for k, l := range list {
            val, key := groupNameList[l.GroupId]
            if key {
                list[k].GroupName = val
            }
        }
    }

    syncd.RenderJson(c, goweb.JSON{
        "list": list,
    })
    return nil
}
