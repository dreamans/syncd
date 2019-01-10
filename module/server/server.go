// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/govalidate"
    "github.com/dreamans/syncd"
    serverService "github.com/dreamans/syncd/service/server"
)

type ServerParamValid struct {
    GroupId     int         `valid:"required" errmsg:"required=sverver group cannot be empty"`
    Name        string      `valid:"required" errmsg:"required=server name cannot be empty"`
    Ip          string      `valid:"required" errmsg:"required=server Ip cannot be empty"`
    SshPort     int         `valid:"required|int_min=1|int_max=65535" errmsg:"required=ssh port cannot be empty|int_min=ssh port must be between 1 and 65535|int_max=ssh port must be between 1 and 65535"`
}

func ServerNew(c *goweb.Context) error {
    return serverUpdate(c, 0)
}

func ServerEdit(c *goweb.Context) error {
    id := c.PostFormInt("id")
    if id == 0 {
        return syncd.RenderParamError("id can not be empty")
    }
    return serverUpdate(c, id)
}

func serverUpdate(c *goweb.Context, id int) error {
    params := ServerParamValid{
        GroupId: c.PostFormInt("group_id"),
        Name: c.PostForm("name"),
        Ip: c.PostForm("ip"),
        SshPort: c.PostFormInt("ssh_port"),
    }
    if valid := govalidate.NewValidate(&params); !valid.Pass() {
        return syncd.RenderParamError(valid.LastFailed().Msg)
    }
    server := &serverService.Server{
        ID: id,
        GroupId: params.GroupId,
        Name: params.Name,
        Ip: params.Ip,
        SshPort: params.SshPort,
    }
    if err := server.CreateOrUpdate(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

func ServerList(c *goweb.Context) error {
    groupId, offset, limit := c.QueryInt("group_id"), c.QueryInt("offset"), c.GetInt("limit")
    keyword := c.Query("keyword")
    server := &serverService.Server{}
    list, total, err := server.List(keyword, groupId, offset, limit)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "list": list,
        "total": total,
    })
}

func ServerDetail(c *goweb.Context) error {
    server := &serverService.Server{
        ID: c.QueryInt("id"),
    }
    if err := server.Get(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, server)
}

func ServerDelete(c *goweb.Context) error {
    server := &serverService.Server{
        ID: c.PostFormInt("id"),
    }
    if err := server.Delete(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}
