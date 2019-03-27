// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/module/project"
    "github.com/dreamans/syncd/util/gostring"
)

type SpaceFormBind struct {
    Name        string  `form:"name" binding:"required"`
    Description string  `form:"description"`
}

func SpaceDelete(c *gin.Context) {
    id := gostring.Str2Int(c.PostForm("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    member := &project.Member{
        UserId: c.GetInt("user_id"),
        SpaceId: id,
    }
    if in := member.MemberInSpace(); !in {
        render.CustomerError(c, render.CODE_ERR_NO_PRIV, "user is not in the project space")
        return
    }
    space := &project.Space{
        ID: id,
    }
    if err := space.Delete(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.JSON(c, nil)
}

func SpaceDetail(c *gin.Context) {
    id := gostring.Str2Int(c.Query("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    member := &project.Member{
        UserId: c.GetInt("user_id"),
        SpaceId: id,
    }
    if in := member.MemberInSpace(); !in {
        render.CustomerError(c, render.CODE_ERR_NO_PRIV, "user is not in the project space")
        return
    }
    space := &project.Space{
        ID: id,
    }
    if err := space.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.JSON(c, space)
}

func SpaceList(c *gin.Context) {
    var query QueryBind
    if err := c.ShouldBind(&query); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    member := &project.Member{
        UserId: c.GetInt("user_id"),
    }
    spaceIds, err := member.SpaceIdsByUserId()
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    space := &project.Space{}
    list, err := space.List(spaceIds, query.Keyword, query.Offset, query.Limit)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    total, err := space.Total(spaceIds, query.Keyword)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.JSON(c, gin.H{
        "list": list,
        "total": total,
    })
}

func SpaceAdd(c *gin.Context) {
    spaceCreateOrUpdate(c, 0)
}

func SpaceUpdate(c *gin.Context) {
    id := gostring.Str2Int(c.PostForm("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    member := &project.Member{
        UserId: c.GetInt("user_id"),
        SpaceId: id,
    }
    if in := member.MemberInSpace(); !in {
        render.CustomerError(c, render.CODE_ERR_NO_PRIV, "user is not in the project space")
        return
    }
    spaceCreateOrUpdate(c, id)
}

func spaceCreateOrUpdate(c *gin.Context, id int) {
    var spaceForm SpaceFormBind
    if err := c.ShouldBind(&spaceForm); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    space := &project.Space{
        ID: id,
        Name: spaceForm.Name,
        Description: spaceForm.Description,
    }
    if err := space.CreateOrUpdate(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    if id == 0 {
        member := &project.Member{
            UserId: c.GetInt("user_id"),
            SpaceId: space.ID,
        }
        if err := member.Create(); err != nil {
            render.AppError(c, err.Error())
            return
        }
    }
    render.Success(c)
}
