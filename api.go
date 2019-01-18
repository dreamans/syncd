// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package syncd

const (

    API_ROOT = "GET /"

    API_SERVER_GROUP_NEW = "POST /api/server/group/new"
    API_SERVER_GROUP_UPDATE = "POST /api/server/group/update"
    API_SERVER_GROUP_LIST = "GET /api/server/group/list"
    API_SERVER_GROUP_DETAIL = "GET /api/server/group/detail"
    API_SERVER_GROUP_DELETE = "POST /api/server/group/delete"

    API_SERVER_NEW = "POST /api/server/new"
    API_SERVER_UPDATE = "POST /api/server/update"
    API_SERVER_LIST = "GET /api/server/list"
    API_SERVER_DETAIL = "GET /api/server/detail"
    API_SERVER_DELETE = "POST /api/server/delete"

    API_USER_EXISTS = "GET /api/user/checkexists"
    API_USER_NEW = "POST /api/user/new"
    API_USER_UPDATE = "POST /api/user/update"
    API_USER_LIST = "GET /api/user/list"
    API_USER_DETAIL = "GET /api/user/detail"
    API_USER_DELETE = "POST /api/user/delete"

    API_USER_LOGIN = "POST /api/user/login"
    API_USER_LOGOUT = "POST /api/user/logout"
    API_USER_LOGIN_STATUS = "GET /api/user/login/status"
    API_USER_MY_UPDATE = "POST /api/user/my/update"
    API_USER_MY_PASSWORD = "POST /api/user/my/password"

    API_USER_GROUP_EXISTS = "GET /api/user/group/checkexists"
    API_USER_GROUP_NEW = "POST /api/user/group/new"
    API_USER_GROUP_UPDATE = "POST /api/user/group/update"
    API_USER_GROUP_PRIV = "GET /api/user/group/priv"
    API_USER_GROUP_LIST = "GET /api/user/group/list"
    API_USER_GROUP_DETAIL = "GET /api/user/group/detail"
    API_USER_GROUP_DELETE = "POST /api/user/group/delete"

    API_PROJECT_SPACE_NEW = "POST /api/project/space/new"
    API_PROJECT_SPACE_UPDATE = "POST /api/project/space/update"
    API_PROJECT_SPACE_LIST = "GET /api/project/space/list"
    API_PROJECT_SPACE_DETAIL = "GET /api/project/space/detail"
    API_PROJECT_SPACE_DELETE = "POST /api/project/space/delete"
    API_PROJECT_SPACE_EXISTS = "GET /api/project/space/exists"

    API_PROJECT_USER_ADD = "POST /api/project/user/add"
    API_PROJECT_USER_LIST = "GET /api/project/user/list"
    API_PROJECT_USER_REMOVE = "POST /api/project/user/remove"
    API_PROJECT_USER_SEARCH = "GET /api/project/user/search"

    API_PROJECT_NEW = "POST /api/project/new"
    API_PROJECT_UPDATE = "POST /api/project/update"
    API_PROJECT_DETAIL = "GET /api/project/detail"
    API_PROJECT_DELETE = "POST /api/project/delete"
    API_PROJECT_LIST = "GET /api/project/list"
    API_PROJECT_EXISTS = "GET /api/project/exists"
    API_PROJECT_STATUS_CHANGE = "POST /api/project/status/change"
    API_PROJECT_REPO_RESET = "POST /api/project/repo/reset"
    API_PROJECT_SERVER_CHECK = "GET /api/project/server/check"

    API_DEPLOY_APPLY_SPACE_LIST = "GET /api/deploy/apply/space/list"
    API_DEPLOY_APPLY_PROJECT_LIST = "GET /api/deploy/apply/project/list"
    API_DEPLOY_APPLY_PROJECT_ALL = "GET /api/deploy/apply/project/all"
    API_DEPLOY_APPLY_PROJECT_DETAIL = "GET /api/deploy/apply/project/detail"
    API_DEPLOY_APPLY_TAGLIST = "GET /api/deploy/apply/repo/taglist"
    API_DEPLOY_APPLY_COMMITLIST = "GET /api/deploy/apply/repo/commitlist"
    API_DEPLOY_APPLY_SUBMIT = "POST /api/deploy/apply/submit"
    API_DEPLOY_APPLY_LIST = "GET /api/deploy/apply/list"
    API_DEPLOY_APPLY_DETAIL = "GET /api/deploy/apply/detail"
    API_DEPLOY_APPLY_AUDIT = "POST /api/deploy/apply/audit"
    API_DEPLOY_APPLY_UNAUDIT = "POST /api/deploy/apply/unaudit"
    API_DEPLOY_APPLY_DISCARD = "POST /api/deploy/apply/discard"
    API_DEPLOY_APPLY_UPDATE = "POST /api/deploy/apply/update"
    API_DEPLOY_APPLY_LOG = "GET /api/deploy/apply/log"
    API_DEPLOY_DEPLOY_START = "POST /api/deploy/deploy/start"
    API_DEPLOY_DEPLOY_STATUS = "GET /api/deploy/deploy/status"
    API_DEPLOY_DEPLOY_STOP = "POST /api/deploy/deploy/stop"
)

