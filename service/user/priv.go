// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/tinystack/goutil/goslice"
    "github.com/dreamans/syncd"
)

var privToApiMap = map[int][]string{
    SERVER_GROUP_VIEW: []string{
        syncd.API_SERVER_GROUP_LIST,
        syncd.API_SERVER_GROUP_DETAIL,
    },
    SERVER_GROUP_NEW: []string{
        syncd.API_SERVER_GROUP_NEW,
    },
    SERVER_GROUP_EDIT: []string{
        syncd.API_SERVER_GROUP_UPDATE,
        syncd.API_SERVER_GROUP_DETAIL,
    },
    SERVER_GROUP_DEL: []string{
        syncd.API_SERVER_GROUP_DELETE,
    },
    SERVER_VIEW: []string{
        syncd.API_SERVER_LIST,
        syncd.API_SERVER_DETAIL,
    },
    SERVER_NEW: []string{
        syncd.API_SERVER_NEW,
    },
    SERVER_EDIT: []string{
        syncd.API_SERVER_UPDATE,
        syncd.API_SERVER_DETAIL,
    },
    SERVER_DEL: []string{
        syncd.API_SERVER_DELETE,
    },
    USER_ROLE_VIEW: []string{
        syncd.API_USER_GROUP_LIST,
    },
    USER_ROLE_NEW: []string{
        syncd.API_USER_GROUP_NEW,
        syncd.API_USER_GROUP_PRIV,
        syncd.API_USER_GROUP_EXISTS,
    },
    USER_ROLE_EDIT: []string{
        syncd.API_USER_GROUP_PRIV,
        syncd.API_USER_GROUP_DETAIL,
        syncd.API_USER_GROUP_UPDATE,
        syncd.API_USER_GROUP_EXISTS,
    },
    USER_ROLE_DEL: []string{
        syncd.API_USER_GROUP_DELETE,
    },
    USER_VIEW: []string{
        syncd.API_USER_LIST,
    },
    USER_NEW: []string{
        syncd.API_USER_NEW,
        syncd.API_USER_EXISTS,
    },
    USER_EDIT: []string{
        syncd.API_USER_EXISTS,
        syncd.API_USER_UPDATE,
        syncd.API_USER_DETAIL,
    },
    USER_DEL: []string{
        syncd.API_USER_DELETE,
    },
    PROJECT_SPACE_VIEW: []string{
        syncd.API_PROJECT_SPACE_LIST,
    },
    PROJECT_SPACE_NEW: []string{
        syncd.API_PROJECT_SPACE_NEW,
        syncd.API_PROJECT_SPACE_EXISTS,
    },
    PROJECT_SPACE_EDIT: []string{
        syncd.API_PROJECT_SPACE_DETAIL,
        syncd.API_PROJECT_SPACE_EXISTS,
        syncd.API_PROJECT_SPACE_UPDATE,
    },
    PROJECT_SPACE_DEL: []string{
        syncd.API_PROJECT_SPACE_DELETE,
    },
    PROJECT_USER_VIEW: []string{
        syncd.API_PROJECT_USER_LIST,
    },
    PROJECT_USER_NEW: []string{
        syncd.API_PROJECT_USER_ADD,
        syncd.API_PROJECT_USER_SEARCH,
    },
    PROJECT_USER_DEL: []string{
        syncd.API_PROJECT_USER_REMOVE,
    },
    PROJECT_VIEW: []string{
        syncd.API_PROJECT_SPACE_DETAIL,
        syncd.API_PROJECT_LIST,
        syncd.API_PROJECT_DETAIL,
    },
    PROJECT_NEW: []string{
        syncd.API_PROJECT_NEW,
        syncd.API_SERVER_GROUP_LIST,
        syncd.API_PROJECT_EXISTS,
    },
    PROJECT_EDIT: []string{
        syncd.API_PROJECT_DETAIL,
        syncd.API_SERVER_GROUP_LIST,
        syncd.API_PROJECT_UPDATE,
        syncd.API_PROJECT_EXISTS,
    },
    PROJECT_DEL: []string{
        syncd.API_PROJECT_DELETE,
    },
    PROJECT_AUDIT: []string{
        syncd.API_PROJECT_STATUS_CHANGE,
    },
    PROJECT_REPO: []string{
        syncd.API_PROJECT_REPO_RESET,
    },
    PROJECT_CHECK: []string{
        syncd.API_PROJECT_SERVER_CHECK,
    },
    DEPLOY_APPLY: []string{
        syncd.API_DEPLOY_APPLY_SPACE_LIST,
        syncd.API_DEPLOY_APPLY_PROJECT_LIST,
        syncd.API_DEPLOY_APPLY_PROJECT_DETAIL,
        syncd.API_DEPLOY_APPLY_TAGLIST,
        syncd.API_DEPLOY_APPLY_COMMITLIST,
        syncd.API_DEPLOY_APPLY_SUBMIT,
    },
    DEPLOY_VIEW_MY: []string{
        syncd.API_DEPLOY_APPLY_LIST,
        syncd.API_DEPLOY_APPLY_DETAIL,
        syncd.API_DEPLOY_APPLY_PROJECT_ALL,
        syncd.API_DEPLOY_APPLY_LOG,
    },
    DEPLOY_VIEW_ALL: []string{
        syncd.API_DEPLOY_APPLY_LIST,
        syncd.API_DEPLOY_APPLY_DETAIL,
        syncd.API_DEPLOY_APPLY_PROJECT_ALL,
        syncd.API_DEPLOY_APPLY_LOG,
    },
    DEPLOY_AUDIT_MY: []string{
        syncd.API_DEPLOY_APPLY_AUDIT,
        syncd.API_DEPLOY_APPLY_UNAUDIT,
    },
    DEPLOY_AUDIT_ALL: []string{
        syncd.API_DEPLOY_APPLY_AUDIT,
        syncd.API_DEPLOY_APPLY_UNAUDIT,
    },
    DEPLOY_DROP_MY: []string{
        syncd.API_DEPLOY_APPLY_DISCARD,
    },
    DEPLOY_DROP_ALL: []string{
        syncd.API_DEPLOY_APPLY_DISCARD,
    },
    DEPLOY_EDIT_MY: []string{
        syncd.API_DEPLOY_APPLY_PROJECT_DETAIL,
        syncd.API_DEPLOY_APPLY_TAGLIST,
        syncd.API_DEPLOY_APPLY_COMMITLIST,
        syncd.API_DEPLOY_APPLY_UPDATE,
    },
    DEPLOY_DEPLOY_MY: []string{
        syncd.API_DEPLOY_DEPLOY_START,
        syncd.API_DEPLOY_APPLY_DETAIL,
        syncd.API_DEPLOY_DEPLOY_STATUS,
        syncd.API_DEPLOY_DEPLOY_STOP,
    },
    DEPLOY_DEPLOY_ALL: []string{
        syncd.API_DEPLOY_DEPLOY_START,
        syncd.API_DEPLOY_APPLY_DETAIL,
        syncd.API_DEPLOY_DEPLOY_STATUS,
        syncd.API_DEPLOY_DEPLOY_STOP,
    },
}

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

type PrivItem struct {
    Label   string  `json:"label"`
    Value   int     `json:"value"`
}

type PrivGroup struct {
    Label   string      `json:"label"`
    Items   []PrivItem  `json:"items"`
}

const (
    DEPLOY_APPLY      = 1001 // 填写上线单
    DEPLOY_VIEW_MY    = 1002 // 查看上线单(自己)
    DEPLOY_VIEW_ALL   = 1003 // 查看上线单(全部)
    DEPLOY_AUDIT_MY   = 1004 // 审核上线单(自己)
    DEPLOY_AUDIT_ALL  = 1005 // 审核上线单(全部)
    DEPLOY_DEPLOY_MY  = 1006 // 上线操作(自己)
    DEPLOY_DEPLOY_ALL = 1007 // 上线操作(全部)
    DEPLOY_DROP_MY    = 1008 // 废弃上线单(自己)
    DEPLOY_DROP_ALL   = 1009 // 废弃上线单(全部)
    DEPLOY_EDIT_MY    = 1010 // 编辑上线单(自己)

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
    PROJECT_REPO  = 2206 // 项目仓库重置
    PROJECT_CHECK = 2207 // 集群检测

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

var PrivList = []PrivGroup{
    privDeploy, privProject, privUser, privServer,
}

var privDeploy = PrivGroup{
    Label: "发布",
    Items: []PrivItem{
        PrivItem{ Label: "上线单-申请", Value: DEPLOY_APPLY },
        PrivItem{ Label: "上线单-查看", Value: DEPLOY_VIEW_MY },
        PrivItem{ Label: "上线单-编辑", Value: DEPLOY_EDIT_MY},
        PrivItem{ Label: "上线单-审核", Value: DEPLOY_AUDIT_MY },
        PrivItem{ Label: "上线单-上线", Value: DEPLOY_DEPLOY_MY },
        PrivItem{ Label: "上线单-废弃", Value: DEPLOY_DROP_MY },
        PrivItem{ Label: "上线单-查看全部", Value: DEPLOY_VIEW_ALL },
        PrivItem{ Label: "上线单-审核全部", Value: DEPLOY_AUDIT_ALL },
        PrivItem{ Label: "上线单-上线全部", Value: DEPLOY_DEPLOY_ALL },
        PrivItem{ Label: "上线单-废弃全部", Value: DEPLOY_DROP_ALL },
    },
}

var privProject = PrivGroup{
    Label: "项目",
    Items: []PrivItem{
        PrivItem{ Label: "空间-查看", Value: PROJECT_SPACE_VIEW },
        PrivItem{ Label: "空间-新增", Value: PROJECT_SPACE_NEW },
        PrivItem{ Label: "空间-编辑", Value: PROJECT_SPACE_EDIT },
        PrivItem{ Label: "空间-删除", Value: PROJECT_SPACE_DEL },
        PrivItem{ Label: "成员-查看", Value: PROJECT_USER_VIEW },
        PrivItem{ Label: "成员-新增", Value: PROJECT_USER_NEW },
        PrivItem{ Label: "成员-删除", Value: PROJECT_USER_DEL },
        PrivItem{ Label: "项目-查看", Value: PROJECT_VIEW },
        PrivItem{ Label: "项目-新增", Value: PROJECT_NEW },
        PrivItem{ Label: "项目-编辑", Value: PROJECT_EDIT },
        PrivItem{ Label: "项目-删除", Value: PROJECT_DEL },
        PrivItem{ Label: "项目-启用", Value: PROJECT_AUDIT },
        PrivItem{ Label: "项目-仓库重置", Value: PROJECT_REPO },
        PrivItem{ Label: "项目-集群检测", Value: PROJECT_CHECK},
    },
}

var privUser = PrivGroup{
    Label: "用户",
    Items: []PrivItem{
        PrivItem{ Label: "角色-查看", Value: USER_ROLE_VIEW },
        PrivItem{ Label: "角色-新增", Value: USER_ROLE_NEW },
        PrivItem{ Label: "角色-编辑", Value: USER_ROLE_EDIT },
        PrivItem{ Label: "角色-删除", Value: USER_ROLE_DEL },
        PrivItem{ Label: "用户-查看", Value: USER_VIEW },
        PrivItem{ Label: "用户-新增", Value: USER_NEW },
        PrivItem{ Label: "用户-编辑", Value: USER_EDIT },
        PrivItem{ Label: "用户-删除", Value: USER_DEL },
    },
}

var privServer = PrivGroup{
    Label: "服务器",
    Items: []PrivItem{
        PrivItem{ Label: "集群-查看", Value: SERVER_GROUP_VIEW },
        PrivItem{ Label: "集群-新增", Value: SERVER_GROUP_NEW },
        PrivItem{ Label: "集群-编辑", Value: SERVER_GROUP_EDIT },
        PrivItem{ Label: "集群-删除", Value: SERVER_GROUP_DEL },
        PrivItem{ Label: "服务器-查看", Value: SERVER_VIEW },
        PrivItem{ Label: "服务器-新增", Value: SERVER_NEW },
        PrivItem{ Label: "服务器-编辑", Value: SERVER_EDIT },
        PrivItem{ Label: "服务器-删除", Value: SERVER_DEL },
    },
}
