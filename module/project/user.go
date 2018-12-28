// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    projectService "github.com/tinystack/syncd/service/project"
    userService "github.com/tinystack/syncd/service/user"
)

func UserAdd2Space(c *goweb.Context) error {
    spaceId, userId := c.PostFormInt("space_id"), c.PostFormInt("user_id")
    if spaceId == 0 || userId == 0 {
        return syncd.RenderParamError("param empty")
    }
    user := &projectService.User{
        SpaceId: spaceId,
        UserId: userId,
    }
    exists, err := user.CheckUserExists()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if exists {
        return syncd.RenderCustomerError(syncd.CODE_ERR_DATA_REPEAT, "用户已经存在")
    }
    if err := user.Add(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

func UserList(c *goweb.Context) error {
    offset, limit := c.QueryInt("offset"), c.QueryInt("limit")
    projectUser := &projectService.User{}
    list, total, err := projectUser.List(offset, limit)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    var userIdList []int
    for _, l := range list {
        userIdList = append(userIdList, l.UserId)
    }
    user := &userService.User{}
    userList, err := user.GetListByIds(userIdList)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }

    var groupIdList []int
    for _, l := range userList {
        groupIdList = append(groupIdList, l.GroupId)
    }
    group := &userService.Group{}
    groupNameList, err := group.GetNameByIds(groupIdList)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if len(groupNameList) > 0 {
        for k, v := range userList {
            if groupName, exists := groupNameList[v.GroupId]; exists {
                userList[k].GroupName = groupName
            }
        }
    }
    userListMap := map[int]userService.UserItem{}
    for _, l := range userList {
        userListMap[l.ID] = l
    }

    projectUserList := []map[string]interface{}{}
    for _, l := range list {
        user, exists := userListMap[l.UserId];
        if !exists {
            continue
        }
        item := map[string]interface{}{
            "id": l.ID,
            "name": user.Name,
            "email": user.Email,
            "group_name": user.GroupName,
            "lock_status": user.LockStatus,
        }
        projectUserList = append(projectUserList, item)
    }

    return syncd.RenderJson(c, goweb.JSON{
        "list": projectUserList,
        "total": total,
    })
}

func UserRemoveFromSpace(c *goweb.Context) error {
    user := &projectService.User{
        ID: c.PostFormInt("id"),
    }
    if err := user.Delete(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}
