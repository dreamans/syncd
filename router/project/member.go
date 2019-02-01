// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/module/user"
    "github.com/dreamans/syncd/module/project"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/util/gostring"
)

type MemberAddQueryBind struct {
    MemberId    int     `form:"member_id" binding:"required"`
    SpaceId     int     `form:"space_id" binding:"required"`
}

type MemberListQueryBind struct {
    SpaceId     int     `form:"space_id" binding:"required"`
    Offset      int     `form:"offset"`
    Limit       int     `form:"limit" binding:"required,gte=1,lte=999"`
}

func MemberRemove(c *gin.Context) {
    id := gostring.Str2Int(c.PostForm("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }

    member := &project.Member{
        ID: id,
    }
    if err := member.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    m := &project.Member{
        UserId: c.GetInt("user_id"),
        SpaceId: member.SpaceId,
    }
    if in := m.MemberInSpace(); !in {
        render.CustomerError(c, render.CODE_ERR_NO_PRIV, "user is not in the project space")
        return
    }

    if err := member.Delete(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.JSON(c, nil)
}

func MemberList(c *gin.Context) {
    var query MemberListQueryBind
    if err := c.ShouldBind(&query); err != nil {
        render.ParamError(c, err.Error())
        return
    }

    mb := &project.Member{
        UserId: c.GetInt("user_id"),
        SpaceId: query.SpaceId,
    }
    if in := mb.MemberInSpace(); !in {
        render.CustomerError(c, render.CODE_ERR_NO_PRIV, "user is not in the project space")
        return
    }

    m := &project.Member{}
    memberList, err := m.List(query.SpaceId, query.Offset, query.Limit)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    total, err := m.Total(query.SpaceId)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    render.JSON(c, gin.H{
        "list": memberList,
        "total": total,
    })
}

func MemberSearch(c *gin.Context) {
    keyword := c.Query("keyword")
    u := &user.User{}
    list, err := u.List(keyword, 0, 20)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    userList := []map[string]interface{}{}
    for _, l := range list {
        userList = append(userList, map[string]interface{}{
            "id": l.ID,
            "username": l.Username,
            "email": l.Email,
            "status": l.Status,
            "role_name": l.RoleName,
        })
    }
    render.JSON(c, userList)
}

func MemberAdd(c *gin.Context) {
    var query MemberAddQueryBind
    if err := c.ShouldBind(&query); err != nil {
        render.ParamError(c, err.Error())
        return
    }

    mb := &project.Member{
        UserId: c.GetInt("user_id"),
        SpaceId: query.SpaceId,
    }
    if in := mb.MemberInSpace(); !in {
        render.CustomerError(c, render.CODE_ERR_NO_PRIV, "user is not in the project space")
        return
    }

    var (
        err error
        exists bool
    )

    u := &user.User{
        ID: query.MemberId,
    }
    exists , err = u.Exists()
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    if !exists {
        render.NoDataError(c, "user not exists")
        return
    }

    member := project.Member{
        UserId: query.MemberId,
        SpaceId: query.SpaceId,
    }
    exists, err = member.Exists()
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    if exists {
        render.RepeatError(c, "user have exists in project space")
        return
    }
    if err := member.Create(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    render.Success(c)
}
