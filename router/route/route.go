// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
    "github.com/dreamans/syncd"
    "github.com/dreamans/syncd/router/user"
	"github.com/dreamans/syncd/router/server"
	"github.com/dreamans/syncd/router/project"
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
		api.GET("/user/detail", user.UserDetail)
		api.POST("/user/delete", user.UserDelete)

		api.POST("/project/space/add", project.SpaceAdd)
		api.POST("/project/space/update", project.SpaceUpdate)
		api.GET("/project/space/list", project.SpaceList)
		api.GET("/project/space/detail", project.SpaceDetail)
		api.POST("/project/space/delete", project.SpaceDelete)
		api.GET("/project/member/search", project.MemberSearch)
		api.POST("/project/member/add", project.MemberAdd)
		api.GET("/project/member/list", project.MemberList)
		api.POST("/project/member/remove", project.MemberRemove)
		api.POST("/project/add", project.ProjectAdd)
		api.POST("/project/update", project.ProjectUpdate)
		api.GET("/project/list", project.ProjectList)
		api.POST("/project/switchstatus", project.ProjectSwitchStatus)
		api.GET("/project/detail", project.ProjectDetail)
		api.POST("/project/delete", project.ProjectDelete)
    }
}
