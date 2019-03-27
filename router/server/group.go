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

type GroupForm struct {
    ID		int     `form:"id"`
    Name	string  `form:"name" binding:"required"`
}

func GroupAdd(c *gin.Context) {
    var groupForm GroupForm
    if err := c.ShouldBind(&groupForm); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    group := &server.Group{
        Name: groupForm.Name,
    }
    if err := group.Create(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.Success(c)
}

func GroupList(c *gin.Context) {
    var query QueryBind
    if err := c.ShouldBind(&query); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    group := &server.Group{}
    list, err := group.List(query.Keyword, query.Offset, query.Limit)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    total, err := group.Total(query.Keyword)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.JSON(c, gin.H{
        "list": list,
        "total": total,
    })
}

func GroupDelete(c *gin.Context) {
    id := gostring.Str2Int(c.PostForm("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    group := &server.Group{
        ID: id,
    }
    if err := group.Delete(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.JSON(c, nil)
}

func GroupDetail(c *gin.Context) {
    id := gostring.Str2Int(c.Query("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    group := &server.Group{
        ID: id,
    }
    if err := group.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.JSON(c, group)
}

func GroupUpdate(c *gin.Context) {
    var groupForm GroupForm
    if err := c.ShouldBind(&groupForm); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    if groupForm.ID == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    group := &server.Group{
        ID: groupForm.ID,
        Name: groupForm.Name,
    }
    if err := group.Update(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.Success(c)
}
