// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route
import (
    "net/http"

    "github.com/tinystack/goweb"
    "github.com/dreamans/syncd"
    serverModule "github.com/dreamans/syncd/module/server"
    userModule "github.com/dreamans/syncd/module/user"
    projectModule "github.com/dreamans/syncd/module/project"
    deployModule "github.com/dreamans/syncd/module/deploy"
)

func init() {
    handler()
}

func handler() {
    h := map[string]goweb.HandlerFunc{
        syncd.API_ROOT: func(c *goweb.Context) error {
            c.Json(http.StatusOK, goweb.JSON{
                "code": 0,
                "message": "welcome to visit syncd API service",
            })
            return nil
        },

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
        syncd.API_USER_MY_UPDATE: userModule.MyUpdate,
        syncd.API_USER_MY_PASSWORD: userModule.MyPasswordUpdate,

        syncd.API_USER_GROUP_EXISTS: userModule.GroupExists,
        syncd.API_USER_GROUP_NEW: userModule.GroupNew,
        syncd.API_USER_GROUP_UPDATE: userModule.GroupEdit,
        syncd.API_USER_GROUP_PRIV: userModule.GroupPlainPriv,
        syncd.API_USER_GROUP_LIST: userModule.GroupList,
        syncd.API_USER_GROUP_DETAIL: userModule.GroupDetail,
        syncd.API_USER_GROUP_DELETE: userModule.GroupDelete,

        // project api
        syncd.API_PROJECT_SPACE_NEW: projectModule.SpaceNew,
        syncd.API_PROJECT_SPACE_UPDATE: projectModule.SpaceEdit,
        syncd.API_PROJECT_SPACE_LIST: projectModule.SpaceList,
        syncd.API_PROJECT_SPACE_DETAIL: projectModule.SpaceDetail,
        syncd.API_PROJECT_SPACE_DELETE: projectModule.SpaceDelete,
        syncd.API_PROJECT_SPACE_EXISTS: projectModule.SpaceExists,

        syncd.API_PROJECT_USER_ADD: projectModule.UserAdd,
        syncd.API_PROJECT_USER_LIST: projectModule.UserList,
        syncd.API_PROJECT_USER_REMOVE: projectModule.UserRemove,
        syncd.API_PROJECT_USER_SEARCH: projectModule.UserSearch,

        syncd.API_PROJECT_NEW: projectModule.ProjectNew,
        syncd.API_PROJECT_UPDATE: projectModule.ProjectEdit,
        syncd.API_PROJECT_DETAIL: projectModule.ProjectDetail,
        syncd.API_PROJECT_DELETE: projectModule.ProjectDelete,
        syncd.API_PROJECT_LIST: projectModule.ProjectList,
        syncd.API_PROJECT_EXISTS: projectModule.ProjectExists,
        syncd.API_PROJECT_STATUS_CHANGE: projectModule.ProjectChangeStatus,
        syncd.API_PROJECT_REPO_RESET: projectModule.RepoReset,
        syncd.API_PROJECT_SERVER_CHECK: projectModule.ServerCheck,

        // deploy api
        syncd.API_DEPLOY_APPLY_SPACE_LIST: deployModule.ApplySpaceList,
        syncd.API_DEPLOY_APPLY_PROJECT_LIST: deployModule.ApplyProjectList,
        syncd.API_DEPLOY_APPLY_PROJECT_DETAIL: deployModule.ApplyProjectDetail,
        syncd.API_DEPLOY_APPLY_TAGLIST: deployModule.ApplyRepoTagList,
        syncd.API_DEPLOY_APPLY_SUBMIT: deployModule.ApplySubmit,
        syncd.API_DEPLOY_APPLY_COMMITLIST: deployModule.ApplyRepoCommitList,
        syncd.API_DEPLOY_APPLY_LIST: deployModule.ApplyList,
        syncd.API_DEPLOY_APPLY_DETAIL: deployModule.ApplyDetail,
        syncd.API_DEPLOY_APPLY_AUDIT: deployModule.ApplyAudit,
        syncd.API_DEPLOY_APPLY_UNAUDIT: deployModule.ApplyUnAudit,
        syncd.API_DEPLOY_APPLY_DISCARD: deployModule.ApplyDiscard,
        syncd.API_DEPLOY_APPLY_PROJECT_ALL: deployModule.ApplyProjectAll,
        syncd.API_DEPLOY_APPLY_UPDATE: deployModule.ApplyUpdate,
        syncd.API_DEPLOY_APPLY_LOG: deployModule.ApplyLog,
        syncd.API_DEPLOY_DEPLOY_START: deployModule.DeployStart,
        syncd.API_DEPLOY_DEPLOY_STATUS: deployModule.DeployStatus,
        syncd.API_DEPLOY_DEPLOY_STOP: deployModule.DeployStop,
    }

    for k, v := range h {
        register(k, v)
    }
}
