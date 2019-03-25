// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/dreamans/syncd/util/goslice"
    reqApi "github.com/dreamans/syncd/router/route/api"
)

var apiToPrivMap = map[string][]int{}

func init() {
    for priv, apiList := range privToApiMap{
        for _, api := range apiList {
            privMap, _ := apiToPrivMap[api]
            apiToPrivMap[api] = append(privMap, priv)
        }
    }
}

func CheckHavePriv(api string, priv []int) bool {
    privMap, exists := apiToPrivMap[api]
    if !exists {
        return false
    }
    return len(goslice.SliceIntersectInt(privMap, priv)) > 0
}

func PrivIn(privCode int, privList []int) bool {
    return goslice.InSliceInt(privCode, privList)
}

const (
    DEPLOY_APPLY      = 1001 // 填写上线单
    DEPLOY_VIEW       = 1002 // 查看上线单
    DEPLOY_AUDIT      = 1003 // 审核上线单
    DEPLOY_DEPLOY     = 1004 // 上线操作
    DEPLOY_DROP       = 1005 // 废弃上线单
    DEPLOY_EDIT       = 1006 // 编辑上线单

    PROJECT_SPACE_VIEW  = 2001 // 查看空间
    PROJECT_SPACE_NEW   = 2002 // 新增空间
    PROJECT_SPACE_EDIT  = 2003 // 编辑空间
    PROJECT_SPACE_DEL   = 2004 // 删除空间
    PROJECT_USER_VIEW = 2100 // 查看成员
    PROJECT_USER_NEW  = 2101 // 新增成员
    PROJECT_USER_DEL  = 2102 // 删除成员
    PROJECT_VIEW  = 2201 // 查看项目
    PROJECT_NEW   = 2202 // 新增项目
    PROJECT_EDIT  = 2203 // 编辑项目
    PROJECT_DEL   = 2204 // 删除项目
    PROJECT_AUDIT = 2205 // 启用项目
    PROJECT_BUILD = 2206 // 构建设置
    PROJECT_HOOK  = 2207 // Hook设置

    USER_ROLE_VIEW = 3001 // 查看角色
    USER_ROLE_NEW  = 3002 // 新增角色
    USER_ROLE_EDIT = 3003 // 编辑角色
    USER_ROLE_DEL  = 3004 // 删除角色
    USER_VIEW = 3101 // 查看用户
    USER_NEW  = 3102 // 新增用户
    USER_EDIT = 3103 // 编辑用户
    USER_DEL  = 3104 // 删除用户

    SERVER_GROUP_VIEW = 4001 // 查看集群
    SERVER_GROUP_NEW  = 4002 // 新增集群
    SERVER_GROUP_EDIT = 4003 // 编辑集群
    SERVER_GROUP_DEL  = 4004 // 删除集群
    SERVER_VIEW  = 4101 // 查看服务器
    SERVER_NEW   = 4102 // 新增服务器
    SERVER_EDIT  = 4103 // 编辑服务器
    SERVER_DEL   = 4104 // 删除服务器
)

var privToApiMap = map[int][]string{
    SERVER_GROUP_VIEW: []string{
        reqApi.SERVER_GROUP_LIST,
    },
    SERVER_GROUP_NEW: []string{
        reqApi.SERVER_GROUP_ADD,
    },
    SERVER_GROUP_EDIT: []string{
        reqApi.SERVER_GROUP_DETAIL,
        reqApi.SERVER_GROUP_UPDATE,
    },
    SERVER_GROUP_DEL: []string{
        reqApi.SERVER_GROUP_DELETE,
    },
    SERVER_VIEW: []string{
        reqApi.SERVER_LIST,
    },
    SERVER_NEW: []string{
        reqApi.SERVER_GROUP_LIST,
        reqApi.SERVER_ADD,
    },
    SERVER_EDIT: []string{
        reqApi.SERVER_DETAIL,
        reqApi.SERVER_UPDATE,
    },
    SERVER_DEL: []string{
        reqApi.SERVER_DELETE,
    },
    USER_ROLE_VIEW: []string{
        reqApi.USER_ROLE_LIST,
    },
    USER_ROLE_NEW: []string{
        reqApi.USER_ROLE_ADD,
        reqApi.USER_ROLE_PRIV_LIST,
    },
    USER_ROLE_EDIT: []string{
        reqApi.USER_ROLE_DETAIL,
        reqApi.USER_ROLE_UPDATE,
        reqApi.USER_ROLE_PRIV_LIST,
    },
    USER_ROLE_DEL: []string{
        reqApi.USER_ROLE_DELETE,
    },
    USER_VIEW: []string{
        reqApi.USER_LIST,
    },
    USER_NEW: []string{
        reqApi.USER_ADD,
        reqApi.USER_ROLE_LIST,
        reqApi.USER_EXISTS,
    },
    USER_EDIT: []string{
        reqApi.USER_ROLE_LIST,
        reqApi.USER_DETAIL,
        reqApi.USER_UPDATE,
        reqApi.USER_EXISTS,
    },
    USER_DEL: []string{
        reqApi.USER_DELETE,
    },
    PROJECT_SPACE_VIEW: []string{
        reqApi.PROJECT_SPACE_LIST,
    },
    PROJECT_SPACE_NEW: []string{
        reqApi.PROJECT_SPACE_ADD,
    },
    PROJECT_SPACE_EDIT: []string{
        reqApi.PROJECT_SPACE_DETAIL,
        reqApi.PROJECT_SPACE_UPDATE,
    },
    PROJECT_SPACE_DEL: []string{
        reqApi.PROJECT_SPACE_DELETE,
    },
    PROJECT_USER_VIEW: []string{
        reqApi.PROJECT_SPACE_DETAIL,
        reqApi.PROJECT_SPACE_LIST,
        reqApi.PROJECT_MEMBER_LIST,
    },
    PROJECT_USER_NEW: []string{
        reqApi.PROJECT_MEMBER_ADD,
        reqApi.PROJECT_MEMBER_SEARCH,
    },
    PROJECT_USER_DEL: []string{
        reqApi.PROJECT_MEMBER_REMOVE,
    },
    PROJECT_VIEW: []string{
        reqApi.SERVER_GROUP_LIST,
        reqApi.PROJECT_SPACE_LIST,
        reqApi.PROJECT_SPACE_DETAIL,
        reqApi.PROJECT_LIST,
        reqApi.PROJECT_DETAIL,
    },
    PROJECT_NEW: []string{
        reqApi.SERVER_GROUP_LIST,
        reqApi.PROJECT_ADD,
    },
    PROJECT_EDIT: []string{
        reqApi.SERVER_GROUP_LIST,
        reqApi.PROJECT_DETAIL,
        reqApi.PROJECT_UPDATE,
    },
    PROJECT_DEL: []string{
        reqApi.PROJECT_DELETE,
    },
    PROJECT_AUDIT: []string{
        reqApi.PROJECT_SWITCHSTATUS,
    },
    PROJECT_BUILD: []string{
        reqApi.PROJECT_DETAIL,
        reqApi.PROJECT_BUILDSCRIPT,
    },
    PROJECT_HOOK: []string{
        reqApi.PROJECT_DETAIL,
        reqApi.PROJECT_HOOKSCRIPT,
    },
    DEPLOY_APPLY: []string{
        reqApi.PROJECT_SPACE_LIST,
        reqApi.PROJECT_LIST,
        reqApi.DEPLOY_APPLY_PROJECT_DETAIL,
        reqApi.DEPLOY_APPLY_SUBMIT,
        reqApi.DEPLOY_APPLY_ROLLBACK,
    },
    DEPLOY_VIEW: []string{
        reqApi.DEPLOY_APPLY_PROJECT_ALL,
        reqApi.DEPLOY_APPLY_LIST,
        reqApi.DEPLOY_APPLY_DETAIL,
        reqApi.DEPLOY_APPLY_PROJECT_DETAIL,
    },
    DEPLOY_AUDIT: []string{
        reqApi.DEPLOY_APPLY_DETAIL,
        reqApi.DEPLOY_APPLY_PROJECT_DETAIL,
        reqApi.DEPLOY_APPLY_AUDIT,
    },
    DEPLOY_EDIT: []string{
        reqApi.DEPLOY_APPLY_DETAIL,
        reqApi.DEPLOY_APPLY_PROJECT_DETAIL,
        reqApi.DEPLOY_APPLY_UPDATE,
    },
    DEPLOY_DROP: []string{
        reqApi.DEPLOY_APPLY_DROP,
    },
    DEPLOY_DEPLOY: []string{
        reqApi.DEPLOY_APPLY_DETAIL,
        reqApi.DEPLOY_APPLY_PROJECT_DETAIL,
        reqApi.DEPLOY_BUILD_START,
        reqApi.DEPLOY_BUILD_STATUS,
        reqApi.DEPLOY_BUILD_STOP,
        reqApi.DEPLOY_DEPLOY_START,
        reqApi.DEPLOY_DEPLOY_STATUS,
        reqApi.DEPLOY_DEPLOY_STOP,
        reqApi.DEPLOY_DEPLOY_ROLLBACK,
    },
}

type PrivItem struct {
    Label   string  `json:"label"`
    Value   int     `json:"value"`
}

type PrivGroup struct {
    Label   string      `json:"label"`
    Items   []PrivItem  `json:"items"`
}

var PrivList = []PrivGroup {
    privProject, privUser, privServer, privDeploy,
}

var privProject = PrivGroup {
    Label: "项目",
    Items: []PrivItem {
        PrivItem{
            Label: "空间-查看",
            Value: PROJECT_SPACE_VIEW,
        },
        PrivItem{
            Label: "空间-新增",
            Value: PROJECT_SPACE_NEW,
        },
        PrivItem{
            Label: "空间-编辑",
            Value: PROJECT_SPACE_EDIT,
        },
        PrivItem{
            Label: "空间-删除",
            Value: PROJECT_SPACE_DEL,
        },
        PrivItem{
            Label: "成员-查看",
            Value: PROJECT_USER_VIEW,
        },
        PrivItem{
            Label: "成员-新增",
            Value: PROJECT_USER_NEW,
        },
        PrivItem{
            Label: "成员-移除",
            Value: PROJECT_USER_DEL,
        },
        PrivItem{
            Label: "项目-查看",
            Value: PROJECT_VIEW,
        },
        PrivItem{
            Label: "项目-新增",
            Value: PROJECT_NEW,
        },
        PrivItem{
            Label: "项目-编辑",
            Value: PROJECT_EDIT,
        },
        PrivItem{
            Label: "项目-删除",
            Value: PROJECT_DEL,
        },
        PrivItem{
            Label: "项目-启用",
            Value: PROJECT_AUDIT,
        },
        PrivItem{
            Label: "项目-构建设置",
            Value: PROJECT_BUILD,
        },
        PrivItem{
            Label: "项目-Hook设置",
            Value: PROJECT_HOOK,
        },
    },
}

var privUser = PrivGroup {
    Label: "用户",
    Items: []PrivItem {
        PrivItem{
            Label: "角色-查看",
            Value: USER_ROLE_VIEW,
        },
        PrivItem{
            Label: "角色-新增",
            Value: USER_ROLE_NEW,
        },
        PrivItem{
            Label: "角色-删除",
            Value: USER_ROLE_DEL,
        },
        PrivItem{
            Label: "角色-编辑",
            Value: USER_ROLE_EDIT,
        },
        PrivItem{
            Label: "用户-查看",
            Value: USER_VIEW,
        },
        PrivItem{
            Label: "用户-新增",
            Value: USER_NEW,
        },
        PrivItem{
            Label: "用户-编辑",
            Value: USER_EDIT,
        },
        PrivItem{
            Label: "用户-删除",
            Value: USER_DEL,
        },
    },
}

var privServer = PrivGroup {
    Label: "服务器",
    Items: []PrivItem {
        PrivItem{
            Label: "集群-查看",
            Value: SERVER_GROUP_VIEW,
        },
        PrivItem{
            Label: "集群-新增",
            Value: SERVER_GROUP_NEW,
        },
        PrivItem{
            Label: "集群-编辑",
            Value: SERVER_GROUP_EDIT,
        },
        PrivItem{
            Label: "集群-删除",
            Value: SERVER_GROUP_DEL,
        },
        PrivItem{
            Label: "服务器-查看",
            Value: SERVER_VIEW,
        },
        PrivItem{
            Label: "服务器-新增",
            Value: SERVER_NEW,
        },
        PrivItem{
            Label: "服务器-编辑",
            Value: SERVER_EDIT,
        },
        PrivItem{
            Label: "服务器-删除",
            Value: SERVER_DEL,
        },
    },
}

var privDeploy = PrivGroup {
    Label: "发布",
    Items: []PrivItem {
        PrivItem{
            Label: "上线单-提交申请",
            Value: DEPLOY_APPLY,
        },
        PrivItem{
            Label: "上线单-查看",
            Value: DEPLOY_VIEW,
        },
        PrivItem{
            Label: "上线单-编辑",
            Value: DEPLOY_EDIT,
        },
        PrivItem{
            Label: "上线单-审核",
            Value: DEPLOY_AUDIT,
        },
        PrivItem{
            Label: "上线单-部署",
            Value: DEPLOY_DEPLOY,
        },
        PrivItem{
            Label: "上线单-废弃",
            Value: DEPLOY_DROP,
        },
    },
}
