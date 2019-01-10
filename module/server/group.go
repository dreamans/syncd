// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "github.com/tinystack/goweb"
    "github.com/dreamans/syncd"
    serverService "github.com/dreamans/syncd/service/server"
)

func GroupNew(c *goweb.Context) error {
    name := c.PostForm("name")
    if name == "" {
        return syncd.RenderParamError("name can not be empty")
    }
    return groupUpdate(c, 0, name)
}

func GroupEdit(c *goweb.Context) error {
    id, name := c.PostFormInt("id"), c.PostForm("name")
    if name == "" || id == 0 {
        return syncd.RenderParamError("id or name can not be empty")
    }
    return groupUpdate(c, id, name)
}

func groupUpdate(c *goweb.Context, id int, name string) error {
    serverGroup := &serverService.Group{
        ID: id,
        Name: name,
    }
    if err := serverGroup.CreateOrUpdate(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
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

func GroupList(c *goweb.Context) error {
    offset, limit, keyword := c.QueryInt("offset"), c.GetInt("limit"), c.Query("keyword")
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

func GroupDelete(c *goweb.Context) error {
    serverGroup := &serverService.Group{
        ID: c.PostFormInt("id"),
    }
    if err := serverGroup.Delete(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}
