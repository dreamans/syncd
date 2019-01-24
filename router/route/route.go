// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
    "github.com/dreamans/syncd"
    "github.com/dreamans/syncd/router/user"
    "github.com/dreamans/syncd/router/server"
)

func RegisterRoute() {
    api := syncd.App.Gin.Group("/api")
    {
        api.POST("/login", user.Login)
        api.GET("/login/status", user.LoginStatus)

        api.POST("/server/group/add", server.GroupAdd)
        api.GET("/server/group/list", server.GroupList)
        api.POST("/server/group/delete", server.GroupDelete)
        api.GET("/server/group/detail", server.GroupDetail)
        api.POST("/server/group/update", server.GroupUpdate)
        api.POST("/server/add", server.ServerAdd)
        api.POST("/server/update", server.ServerUpdate)
        api.GET("/server/list", server.ServerList)
        api.POST("/server/delete", server.ServerDelete)
        api.GET("/server/detail", server.ServerDetail)

        api.POST("/user/role/add", user.RoleAdd)
        api.POST("/user/role/update", user.RoleUpdate)
        api.GET("/user/role/list", user.RoleList)
        api.GET("/user/role/detail", user.RoleDetail)
        api.POST("/user/role/delete", user.RoleDelete)
        api.POST("/user/add", user.UserAdd)
        api.POST("/user/update", user.UserUpdate)
        api.GET("/user/list", user.UserList)
        api.GET("/user/exists", user.UserExists)
    }
}
