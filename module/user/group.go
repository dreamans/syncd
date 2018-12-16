// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    "github.com/tinystack/syncd/route"
    //privModel "github.com/tinystack/syncd/model/user/priv"
    userService "github.com/tinystack/syncd/service/user"
)

func init() {
    route.Register(route.API_USER_GROUP_UPDATE, updateUserGroup)
    route.Register(route.API_USER_GROUP_LIST, listUserGroup)
    route.Register(route.API_USER_GROUP_DETAIL, detailUserGroup)
    route.Register(route.API_USER_GROUP_PRIV, privUserGroup)
}

func updateUserGroup(c *goweb.Context) error {
    groupId := c.PostFormInt("id")
    groupName := c.PostForm("name")
    groupPriv := c.PostFormArray("priv")
    if groupName == "" {
        syncd.RenderParamError(c, "user group name can not empty")
        return nil
    }

    if err := userService.UpdateUserGroup(groupId, groupName, groupPriv); err != nil {
        syncd.RenderParamError(c, err.Error())
        return nil
    }

    syncd.RenderJson(c, nil)
    return nil
}

func listUserGroup(c *goweb.Context) error {
    keyword := c.Query("keyword")
    offset, limit := c.QueryInt("offset"), c.QueryInt("limit")

    list, total, err := userService.GetUserGroupList(keyword, offset, limit)
    if err != nil {
        syncd.RenderParamError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, goweb.JSON{
        "list": list,
        "total": total,
    })
    return nil
}

func detailUserGroup(c *goweb.Context) error {
    detail, err := userService.GetUserGroupDetail(c.QueryInt("id"))
    if err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, detail)
    return nil
}

func privUserGroup(c *goweb.Context) error {
    /*
    id := c.QueryInt("id")
    if id == 0 {
        syncd.RenderParamError(c, "id can not be empty")
        return nil
    }
    detail, ok := groupModel.Get(id)
    if !ok {
        syncd.RenderAppError(c, "get user group detail data failed")
        return nil
    }

    privModel.PrivList
    */
    return nil
}
