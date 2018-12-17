// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/govalidate"
    "github.com/tinystack/syncd"
    "github.com/tinystack/syncd/route"
    serverService "github.com/tinystack/syncd/service/server"
)

func init() {
    route.Register(route.API_SERVER_UPDATE, updateServer)
    route.Register(route.API_SERVER_LIST, listServer)
    route.Register(route.API_SERVER_DETAIL, detailServer)
    route.Register(route.API_SERVER_DELETE, deleteServer)
}

type ServerParamValid struct {
    GroupId     int         `valid:"required" errmsg:"required=sverver group cannot be empty"`
    Name        string      `valid:"required" errmsg:"required=server name cannot be empty"`
    Ip          string      `valid:"required" errmsg:"required=server Ip cannot be empty"`
    SshPort     int         `valid:"required|int_min=1|int_max=65535" errmsg:"required=ssh port cannot be empty|int_min=ssh port must be between 1 and 65535|int_max=ssh port must be between 1 and 65535"`
}

func updateServer(c *goweb.Context) error {
    params := ServerParamValid{
        GroupId: c.PostFormInt("group_id"),
        Name: c.PostForm("name"),
        Ip: c.PostForm("ip"),
        SshPort: c.PostFormInt("ssh_port"),
    }
    if valid := govalidate.NewValidate(&params); !valid.Pass() {
        syncd.RenderParamError(c, valid.LastFailed().Msg)
        return nil
    }
    server := &serverService.Server{
        ID: c.PostFormInt("id"),
        GroupId: params.GroupId,
        Name: params.Name,
        Ip: params.Ip,
        SshPort: params.SshPort,
    }
    if err := server.CreateOrUpdate(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}

func listServer(c *goweb.Context) error {
    groupId, offset, limit := c.QueryInt("group_id"), c.QueryInt("offset"), c.QueryInt("limit")
    keyword := c.Query("keyword")

    server := &serverService.Server{}
    list, total, err := server.List(keyword, groupId, offset, limit)
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

func detailServer(c *goweb.Context) error {
    server := &serverService.Server{
        ID: c.QueryInt("id"),
    }
    if err := server.Get(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, server)
    return nil
}

func deleteServer(c *goweb.Context) error {
    server := &serverService.Server{
        ID: c.PostFormInt("id"),
    }
    if err := server.Delete(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}
