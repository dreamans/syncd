// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/goutil/gostring"
    "github.com/tinystack/syncd"
    serverService "github.com/tinystack/syncd/service/server"
)

func GroupUpdate(c *goweb.Context) error {
    id, name := c.PostFormInt("id"), c.PostForm("name")
    if name == "" {
        return syncd.RenderParamError("name can not be empty")
    }
    serverGroup := &serverService.Group{
        ID: id,
        Name: name,
    }
    if err := serverGroup.CreateOrUpdate(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

func GroupList(c *goweb.Context) error {
    offset, limit := c.QueryInt("offset"), c.QueryInt("limit")
    keyword := c.Query("keyword")
    serverGroup := &serverService.Group{}
    list, total, err := serverGroup.List(keyword, offset, limit)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "list": list,
        "total": total,
    })
}

func GroupMulti(c *goweb.Context) error {
    ids := gostring.StrSplit2IntSlice(c.Query("ids"), ",")
    serverGroup := &serverService.Group{}
    list, err := serverGroup.GetMultiById(ids)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "list": list,
    })
}

func GroupDetail(c *goweb.Context) error {
    serverGroup := &serverService.Group{
        ID: c.QueryInt("id"),
    }
    if err := serverGroup.Detail(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, serverGroup)
}

func GroupDelete(c *goweb.Context) error {
    serverGroup := &serverService.Group{
        ID: c.PostFormInt("id"),
    }
    if err := serverGroup.Delete(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}
