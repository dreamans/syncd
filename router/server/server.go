// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/module/server"
    "github.com/dreamans/syncd/util/gostring"
)

type QueryBind struct {
    Keyword	string  `form:"keyword"`
    Offset	int     `form:"offset"`
    Limit	int     `form:"limit" binding:"required,gte=1,lte=999"`
}

type ServerForm struct {
    GroupId	int     `form:"group_id" binding:"required"`
    Name	string  `form:"name" binding:"required"`
    Ip		string	`form:"ip" binding:"required"`
    SSHPort	int     `form:"ssh_port" binding:"required,gte=1,lte=65535"`
}

func ServerAdd(c *gin.Context) {
    serverCreateOrUpdate(c, 0)
}

func ServerUpdate(c *gin.Context) {
    id := gostring.Str2Int(c.PostForm("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    serverCreateOrUpdate(c, id)
}

func serverCreateOrUpdate(c *gin.Context, id int) {
    var serverForm ServerForm
    if err := c.ShouldBind(&serverForm); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    server := &server.Server{
        ID: id,
        GroupId: serverForm.GroupId,
        Name: serverForm.Name,
        Ip: serverForm.Ip,
        SSHPort: serverForm.SSHPort,
    }
    if err := server.CreateOrUpdate(); err != nil {
        render.AppError(c, err.Error())
        return
    }
	render.Success(c)
}

func ServerList(c *gin.Context) {
    var query QueryBind
    if err := c.ShouldBind(&query); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    ser := &server.Server{}
    list, err := ser.List(query.Keyword, query.Offset, query.Limit)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    total, err := ser.Total(query.Keyword)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.JSON(c, gin.H{
        "list": list,
        "total": total,
    })
}

func ServerDelete(c *gin.Context) {
    id := gostring.Str2Int(c.PostForm("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    server := &server.Server{
        ID: id,
    }
    if err := server.Delete(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.JSON(c, nil)
}

func ServerDetail(c *gin.Context) {
    id := gostring.Str2Int(c.Query("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    ser := &server.Server{
        ID: id,
    }
    if err := ser.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    group := &server.Group{
        ID: ser.GroupId,
    }
    if err := group.Detail(); err == nil {
        ser.GroupName = group.Name
    }

    render.JSON(c, ser)
}
