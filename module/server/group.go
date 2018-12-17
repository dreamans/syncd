// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/goutil"
    "github.com/tinystack/syncd"
    "github.com/tinystack/syncd/route"
    serverService "github.com/tinystack/syncd/service/server"
)

func init() {
    route.Register(route.API_SERVER_GROUP_UPDATE, updateServerGroup)
    route.Register(route.API_SERVER_GROUP_LIST, listServerGroup)
    route.Register(route.API_SERVER_GROUP_DETAIL, detailServerGroup)
    route.Register(route.API_SERVER_GROUP_DELETE, deleteServerGroup)
    route.Register(route.API_SERVER_GROUP_MULTI, multiServerGroup)
}

func updateServerGroup(c *goweb.Context) error {
    id, name := c.PostFormInt("id"), c.PostForm("name")
    if name == "" {
        syncd.RenderParamError(c, "name can not be empty")
        return nil
    }
    serverGroup := &serverService.Group{
        ID: id,
        Name: name,
    }
    if err := serverGroup.CreateOrUpdate(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}

func listServerGroup(c *goweb.Context) error {
    offset, limit := c.QueryInt("offset"), c.QueryInt("limit")
    keyword := c.Query("keyword")

    serverGroup := &serverService.Group{}
    list, total, err := serverGroup.List(keyword, offset, limit)
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

func multiServerGroup(c *goweb.Context) error {
    ids := goutil.StrSplit2IntSlice(c.Query("ids"), ",")
    serverGroup := &serverService.Group{}
    list, err := serverGroup.GetMultiById(ids)
    if err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, goweb.JSON{
        "list": list,
    })
    return nil
}

func detailServerGroup(c *goweb.Context) error {
    serverGroup := &serverService.Group{
        ID: c.QueryInt("id"),
    }
    if err := serverGroup.Detail(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, serverGroup)
    return nil
}

func deleteServerGroup(c *goweb.Context) error {
    serverGroup := &serverService.Group{
        ID: c.PostFormInt("id"),
    }
    if err := serverGroup.Delete(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}
