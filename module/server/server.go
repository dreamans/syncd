// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/govalidate"
    "github.com/tinystack/syncd"
    "github.com/tinystack/syncd/route"
    serverModel "github.com/tinystack/syncd/model/server"
)

func init() {
    route.Register(route.API_SERVER_UPDATE, updateServer)
    route.Register(route.API_SERVER_LIST, listServer)
    route.Register(route.API_SERVER_DETAIL, detailServer)
    route.Register(route.API_SERVER_DELETE, deleteServer)
    route.Register(route.API_SERVER_MULTI, multiServer)
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
    serverId := c.PostFormInt("id")

    s := serverModel.Server{
        GroupId: params.GroupId,
        Name: params.Name,
        Ip: params.Ip,
        SshPort: params.SshPort,
    }
    var ok bool
    if serverId > 0 {
        ok = serverModel.Update(serverId, s)
    } else {
        ok = serverModel.Create(&s)
    }
    if !ok {
        syncd.RenderAppError(c, "server data update failed")
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}

func listServer(c *goweb.Context) error {
    offset, limit := c.QueryInt("offset"), c.QueryInt("limit")
    list, ok := serverModel.List("id, group_id, name, ip, ssh_port", offset, limit)
    if !ok {
        syncd.RenderAppError(c, "get server list data failed")
        return nil
    }

    total, ok := serverModel.Total()
    if !ok {
        syncd.RenderAppError(c, "get server total count failed")
        return nil
    }

    syncd.RenderJson(c, goweb.JSON{
        "list": list,
        "total": total,
    })
    return nil
}

func multiServer(c *goweb.Context) error {
    groupId := c.QueryInt("group_id")
    syncd.Logger.Info("%v", groupId)
    list, ok := serverModel.Multi("id, group_id, name, ip, ssh_port", groupId)
    if !ok {
        syncd.RenderAppError(c, "get server list data failed")
        return nil
    }
    syncd.RenderJson(c, goweb.JSON{
        "list": list,
    })
    return nil
}

func detailServer(c *goweb.Context) error {
    id := c.QueryInt("id")
    if id == 0 {
        syncd.RenderParamError(c, "id can not be empty")
        return nil
    }
    detail, ok := serverModel.Get(id)
    if !ok {
        syncd.RenderAppError(c, "get server detail data failed")
        return nil
    }
    syncd.RenderJson(c, goweb.JSON{
        "detail": detail,
    })
    return nil
}

func deleteServer(c *goweb.Context) error {
    id := c.PostFormInt("id")
    if id == 0 {
        syncd.RenderParamError(c, "id can not be empty")
        return nil
    }
    ok := serverModel.Delete(id)
    if !ok {
        syncd.RenderAppError(c, "delete server data failed")
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}
