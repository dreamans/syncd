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

func GroupNew(c *goweb.Context) error {
    return groupUpdate(c, 0)
}

func GroupEdit(c *goweb.Context) error {
     id := c.PostFormInt("id")
     if id == 0 {
         return syncd.RenderParamError("user group id can not empty")
     }
     return groupUpdate(c, id)
}

func groupUpdate(c *goweb.Context, id int) error {
    name, priv := c.PostForm("name"), c.PostFormArray("priv")
    if name == "" {
        return syncd.RenderParamError("user group name can not empty")
    }
    userGroup := &userService.Group{
        ID: id,
        Name: name,
        Priv: gostring.StrSlice2IntSlice(gostring.StrFilterSliceEmpty(priv)),
    }
    exists, err := userGroup.CheckGroupExists()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if exists {
        return syncd.RenderCustomerError(syncd.CODE_ERR_DATA_REPEAT, "group exists")
    }
    if err := userGroup.CreateOrUpdate(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

func GroupList(c *goweb.Context) error {
    offset, limit, keyword := c.QueryInt("offset"), c.GetInt("limit"), c.Query("keyword")
    userGroup := &userService.Group{}
    list, total, err := userGroup.List(keyword, offset, limit)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "list": list,
        "total": total,
    })
}

func GroupDetail(c *goweb.Context) error {
    userGroup := &userService.Group{
        ID: c.QueryInt("id"),
    }
    if err := userGroup.Detail(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, userGroup)
}

func GroupPlainPriv(c *goweb.Context) error {
    return syncd.RenderJson(c, goweb.JSON{
        "list": userService.PrivList,
    })
}

func GroupDelete(c *goweb.Context) error {
    userGroup := &userService.Group{
        ID: c.PostFormInt("id"),
    }
    if err := userGroup.Delete(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

func GroupExists(c *goweb.Context) error {
    keyword, id := c.Query("keyword"), c.QueryInt("id")
    group := userService.Group{
        ID: id,
        Name: keyword,
    }
    exists, err := group.CheckGroupExists()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "exists": exists,
    })
}