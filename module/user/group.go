// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/tinystack/goutil/gostring"
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    userService "github.com/tinystack/syncd/service/user"
)

func GroupUpdate(c *goweb.Context) error {
    id, name, priv := c.PostFormInt("id"), c.PostForm("name"), c.PostFormArray("priv")
    if name == "" {
        return syncd.RenderParamError("user group name can not empty")
    }
    userGroup := &userService.Group{
        ID: id,
        Name: name,
        Priv: gostring.StrSlice2IntSlice(gostring.StrFilterSliceEmpty(priv)),
    }
    if err := userGroup.CreateOrUpdate(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

func GroupList(c *goweb.Context) error {
    offset, limit, keyword := c.QueryInt("offset"), c.QueryInt("limit"), c.Query("keyword")
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
    if err := userGroup.Get(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, userGroup)
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
