// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    "github.com/tinystack/syncd/route"
    "github.com/tinystack/syncd/model/server/group"
)

func init() {
    route.Register(route.API_SERVER_GROUP_UPDATE, updateServerGroup)
}

func updateServerGroup(c *goweb.Context) error {
    groupId := c.PostFormInt("id")
    groupName := c.PostForm("name")
    if groupName == "" {
        syncd.RenderParamError(c, "group name can not empty")
        return nil
    }
    var ok bool
    g := group.Group{
        Name: groupName,
    }
    if groupId > 0 {
        ok = group.Update(groupId, g)
    } else {
        ok = group.Create(&g)
    }
    if !ok {
        syncd.RenderAppError(c, "server group data update failed")
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}
