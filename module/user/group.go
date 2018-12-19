// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/tinystack/goutil"
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    "github.com/tinystack/syncd/route"
    userService "github.com/tinystack/syncd/service/user"
)

func init() {
    route.Register(route.API_USER_GROUP_UPDATE, updateUserGroup)
    route.Register(route.API_USER_GROUP_LIST, listUserGroup)
    route.Register(route.API_USER_GROUP_DETAIL, detailUserGroup)
    route.Register(route.API_USER_GROUP_PRIV, privUserGroup)
    route.Register(route.API_USER_GROUP_DELETE, deleteGroup)
}

func updateUserGroup(c *goweb.Context) error {
    id, name, priv := c.PostFormInt("id"), c.PostForm("name"), c.PostFormArray("priv")
    if name == "" {
        syncd.RenderParamError(c, "user group name can not empty")
        return nil
    }
    userGroup := &userService.Group{
        ID: id,
        Name: name,
        Priv: goutil.StrSlice2IntSlice(goutil.StrFilterSliceEmpty(priv)),
    }
    if err := userGroup.CreateOrUpdate(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}

func listUserGroup(c *goweb.Context) error {
    offset, limit, keyword := c.QueryInt("offset"), c.QueryInt("limit"), c.Query("keyword")
    userGroup := &userService.Group{}
    list, total, err := userGroup.List(keyword, offset, limit)
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

func detailUserGroup(c *goweb.Context) error {
    userGroup := &userService.Group{
        ID: c.QueryInt("id"),
    }
    if err := userGroup.Get(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, userGroup)
    return nil
}

func deleteGroup(c *goweb.Context) error {
    userGroup := &userService.Group{
        ID: c.PostFormInt("id"),
    }
    if err := userGroup.Delete(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}

func privUserGroup(c *goweb.Context) error {
    /*
    id := c.QueryInt("id")
    if id == 0 {
        syncd.RenderParamError(c, "id can not be empty")
        return nil
    }
    detail, ok := groupModel.Get(id)
    if !ok {
        syncd.RenderAppError(c, "get user group detail data failed")
        return nil
    }

    privModel.PrivList
    */
    return nil
}
