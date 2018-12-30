// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route
import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    serverModule "github.com/tinystack/syncd/module/server"
    userModule "github.com/tinystack/syncd/module/user"
    projectModule "github.com/tinystack/syncd/module/project"
)

func init() {
    handler()
}

func handler() {
    h := map[string]goweb.HandlerFunc{
        // server group api
        syncd.API_SERVER_GROUP_NEW: serverModule.GroupNew,
        syncd.API_SERVER_GROUP_UPDATE: serverModule.GroupEdit,
        syncd.API_SERVER_GROUP_LIST: serverModule.GroupList,
        syncd.API_SERVER_GROUP_DETAIL: serverModule.GroupDetail,
        syncd.API_SERVER_GROUP_DELETE: serverModule.GroupDelete,

        // server api
        syncd.API_SERVER_NEW: serverModule.ServerNew,
        syncd.API_SERVER_UPDATE: serverModule.ServerEdit,
        syncd.API_SERVER_LIST: serverModule.ServerList,
        syncd.API_SERVER_DETAIL: serverModule.ServerDetail,
        syncd.API_SERVER_DELETE: serverModule.ServerDelete,

        // user api
        syncd.API_USER_EXISTS: userModule.UserExists,
        syncd.API_USER_NEW: userModule.UserNew,
        syncd.API_USER_UPDATE: userModule.UserEdit,
        syncd.API_USER_LIST: userModule.UserList,
        syncd.API_USER_DETAIL: userModule.UserDetail,
        syncd.API_USER_DELETE: userModule.UserDelete,

        syncd.API_USER_LOGIN: userModule.Login,
        syncd.API_USER_LOGOUT: userModule.Logout,
        syncd.API_USER_LOGIN_STATUS: userModule.LoginStatus,

        syncd.API_USER_GROUP_NEW: userModule.GroupNew,
        syncd.API_USER_GROUP_UPDATE: userModule.GroupEdit,
        syncd.API_USER_GROUP_PRIV: userModule.GroupPlainPriv,
        syncd.API_USER_GROUP_LIST: userModule.GroupList,
        syncd.API_USER_GROUP_DETAIL: userModule.GroupDetail,
        syncd.API_USER_GROUP_DELETE: userModule.GroupDelete,

        // project api
        syncd.
    }

    for k, v := range h {
        register(k, v)
    }
}

