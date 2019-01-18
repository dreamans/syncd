// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/goutil/gois"
    "github.com/dreamans/syncd"
    projectService "github.com/dreamans/syncd/service/project"
    userService "github.com/dreamans/syncd/service/user"
)

func UserSearch(c *goweb.Context) error {
    keyword := c.Query("keyword")
    if keyword == "" {
        return syncd.RenderJson(c, nil)
    }
    user := &userService.User{}
    if gois.IsEmail(keyword) {
        user.Email = keyword
    } else {
        user.Name = keyword
    }
    list, err := user.Search()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    list, err = userService.GroupUserListFillGroupName(list)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "list": list,
    })
}

func UserAdd(c *goweb.Context) error {
    spaceId, userId := c.PostFormInt("space_id"), c.PostFormInt("user_id")
    if spaceId == 0 || userId == 0 {
        return syncd.RenderParamError("param empty")
    }
    user := &projectService.User{
        SpaceId: spaceId,
        UserId: userId,
    }
    exists, err := user.CheckUserInSpace()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if exists {
        return syncd.RenderCustomerError(syncd.CODE_ERR_DATA_REPEAT, "user have exists")
    }
    if err := user.Add(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

func UserList(c *goweb.Context) error {
    var (
        err error
        userList []userService.UserItem
    )
    offset, limit, spaceId := c.QueryInt("offset"), c.GetInt("limit"), c.QueryInt("spaceId")
    projectUser := &projectService.User{}
    list, total, err := projectUser.List(spaceId, offset, limit)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    var userIdList []int
    for _, l := range list {
        userIdList = append(userIdList, l.UserId)
    }
    userList, err = userService.UserGetListByIds(userIdList)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    userList, err = userService.GroupUserListFillGroupName(userList)
    userListMap := map[int]userService.UserItem{}
    for _, l := range userList {
        userListMap[l.ID] = l
    }
    projectUserList := []map[string]interface{}{}
    for _, l := range list {
        user, exists := userListMap[l.UserId]
        var item map[string]interface{}
        if !exists {
            item = map[string]interface{}{
                "id": l.ID,
                "name": "--",
                "email": "--",
                "group_name": "--",
                "lock_status": 0,
            }
        } else {
            item = map[string]interface{}{
                "id": l.ID,
                "name": user.Name,
                "email": user.Email,
                "group_name": user.GroupName,
                "lock_status": user.LockStatus,
            }
        }
        projectUserList = append(projectUserList, item)
    }
    return syncd.RenderJson(c, goweb.JSON{
        "list": projectUserList,
        "total": total,
    })
}

func UserRemove(c *goweb.Context) error {
    user := &projectService.User{
        ID: c.PostFormInt("id"),
    }
    if err := user.Delete(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}
