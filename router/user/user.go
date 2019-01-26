// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/module/user"
    "github.com/dreamans/syncd/util/gostring"
)

type QueryBind struct {
    Keyword     string  `form:"keyword"`
    Offset      int     `form:"offset"`
    Limit       int     `form:"limit" binding:"required,gte=1,lte=999"`
}

type UserForm struct {
    ID          int     `form:"id"`
    RoleId      int     `form:"role_id" binding:"required"`
    Username    string  `form:"username" binding:"required"`
    Password    string  `form:"password"`
    Email       string  `form:"email" binding:"required"`
    Truename    string  `form:"truename"`
    Mobile      string  `form:"mobile"`
    Status      int     `form:"status"`
}

type UserExistsQuery struct {
    ID          int     `form:"id"`
    Username    string  `form:"username"`
    Email       string  `form:"email"`
}

func UserDelete(c *gin.Context) {
	id := gostring.Str2Int(c.PostForm("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
	}
	u := &user.User{
		ID: id,
	}
	if err := u.Delete(); err != nil {
		render.AppError(c, err.Error())
        return
	}
	render.Success(c)
}

func UserDetail(c *gin.Context) {
    id := gostring.Str2Int(c.Query("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    u := &user.User{
        ID: id,
    }
    if err := u.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    detail := map[string]interface{}{
        "id": u.ID,
        "role_id": u.RoleId,
        "username": u.Username,
        "email": u.Email,
        "truename": u.Truename,
        "mobile": u.Mobile,
        "status": u.Status,
    }
    render.JSON(c, detail)
}

func UserExists(c *gin.Context) {
    var query UserExistsQuery
    if err := c.ShouldBind(&query); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    u := &user.User{
        ID: query.ID,
        Username: query.Username,
        Email: query.Email,
    }
    exists, err := u.UserCheckExists()
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.JSON(c, gin.H{
        "exists": exists,
    })
}

func UserList(c *gin.Context) {
    var query QueryBind
    if err := c.ShouldBind(&query); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    u := &user.User{}
    list, err := u.List(query.Keyword, query.Offset, query.Limit)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    total, err := u.Total(query.Keyword)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    var userList []map[string]interface{}
    for _, l := range list {
        userList = append(userList, map[string]interface{}{
            "id": l.ID,
            "role_name": l.RoleName,
            "username": l.Username,
            "truename": l.Truename,
            "email": l.Email,
            "status": l.Status,
            "last_login_time": l.LastLoginTime,
            "last_login_ip": l.LastLoginIp,
        })
    }

    render.JSON(c, gin.H{
        "list": list,
        "total": total,
    })
}

func UserAdd(c *gin.Context) {
    var userForm UserForm
    if err := c.ShouldBind(&userForm); err != nil {
        render.ParamError(c, err.Error())
        return
    }

    if len(userForm.Password) != 32 {
        render.ParamError(c, "password param incorrect")
        return
    }

    userCreateOrUpdate(c, userForm)
}

func UserUpdate(c *gin.Context) {
    var userForm UserForm
    if err := c.ShouldBind(&userForm); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    if userForm.ID == 0 {
        render.ParamError(c, "id cannot empty")
        return
    }
    if userForm.Password != "" && len(userForm.Password) != 32 {
        render.ParamError(c, "password param incorrect")
        return
    }

    userCreateOrUpdate(c, userForm)
}

func userCreateOrUpdate(c *gin.Context, userForm UserForm) {
    var (
        checkUsername, checkEmail *user.User
        exists bool
        err error
    )
    checkUsername = &user.User{
        ID: userForm.ID,
        Username: userForm.Username,
    }
    exists, err = checkUsername.UserCheckExists()
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    if exists {
        render.RepeatError(c, "username have exists")
        return
    }

    checkEmail = &user.User{
        ID: userForm.ID,
        Email: userForm.Email,
    }
    exists, err = checkEmail.UserCheckExists()
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    if exists {
        render.RepeatError(c, "email have exists")
        return
    }

    u := &user.User{
        ID: userForm.ID,
        RoleId: userForm.RoleId,
        Username: userForm.Username,
        Password: userForm.Password,
        Email: userForm.Email,
        Truename: userForm.Truename,
        Mobile: userForm.Mobile,
        Status: userForm.Status,
    }
    if err := u.CreateOrUpdate(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.Success(c)
}