// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/module/user"
    "github.com/dreamans/syncd/util/gostring"
    "github.com/dreamans/syncd/util/goslice"
)

type RoleFormBind struct {
    ID          int     `form:"id"`
    Name        string  `form:"name" binding:"required"`
    Privilege   []int   `form:"privilege"`
}

func RolePrivList(c *gin.Context) {
    render.JSON(c, user.PrivList)
}

func RoleDelete(c *gin.Context) {
    id := gostring.Str2Int(c.PostForm("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    role := &user.Role{
        ID: id,
    }
    if err := role.Delete(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.JSON(c, nil)
}

func RoleDetail(c *gin.Context) {
    id := gostring.Str2Int(c.Query("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    role := &user.Role{
        ID: id,
    }
    if err := role.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.JSON(c, role)
}

func RoleList(c *gin.Context) {
    var query QueryBind
    if err := c.ShouldBind(&query); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    role := &user.Role{}
    list, err := role.List(query.Keyword, query.Offset, query.Limit)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    total, err := role.Total(query.Keyword)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    var roleList []map[string]interface{}
    for _, l := range list {
        roleList = append(roleList, map[string]interface{}{
            "id": l.ID,
            "name": l.Name,
            "ctime": l.Ctime,
        })
    }

    render.JSON(c, gin.H{
        "list": roleList,
        "total": total,
    })
}

func RoleAdd(c *gin.Context) {
    roleCreateOrUpdate(c, 0)
}

func RoleUpdate(c *gin.Context) {
    id := gostring.Str2Int(c.PostForm("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    roleCreateOrUpdate(c, id)
}

func roleCreateOrUpdate(c *gin.Context, id int) {
    var roleForm RoleFormBind
    if err := c.ShouldBind(&roleForm); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    role := user.Role{
        ID: roleForm.ID,
        Name: roleForm.Name,
        Privilege: goslice.FilterSliceInt(roleForm.Privilege),
    }
    if err := role.CreateOrUpdate(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.Success(c)
}
