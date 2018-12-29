// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/tinystack/syncd"
    "github.com/tinystack/goutil/goslice"
    userPrivModel "github.com/tinystack/syncd/model/user/priv"
)

var privToApiMap = map[int][]string{
    userPrivModel.DEPLOY_APPLY: []string{
        //syncd.API_PROJECT_SPACE_LIST,
        //syncd.API_PROJECT_LIST,
        //syncd.API_PROJECT_DETAIL,
        //syncd.API_DEPLOY_REPO_COMMITLIST,
        //syncd.API_DEPLOY_REPO_TAGLIST,
        //syncd.API_DEPLOY_APPLY_SUBMIT,
    },

    // project priv
    userPrivModel.PROJECT_SPACE_VIEW: []string{
        syncd.API_PROJECT_SPACE_LIST,
    },
    userPrivModel.PROJECT_SPACE_NEW: []string{
        syncd.API_PROJECT_SPACE_NEW,
        syncd.API_PROJECT_SPACE_EXISTS,
    },
    userPrivModel.PROJECT_SPACE_EDIT: []string{
        syncd.API_PROJECT_SPACE_DETAIL,
        syncd.API_PROJECT_SPACE_UPDATE,
    },
    userPrivModel.PROJECT_SPACE_DEL: []string{
        syncd.API_PROJECT_SPACE_DELETE,
    },
    userPrivModel.PROJECT_VIEW: []string{
        syncd.API_PROJECT_SPACE_DETAIL,
        syncd.API_PROJECT_LIST,
        syncd.API_PROJECT_DETAIL,
    },
    userPrivModel.PROJECT_NEW: []string{
        syncd.API_SERVER_GROUP_LIST,
        syncd.API_PROJECT_EXISTS,
        syncd.API_PROJECT_NEW,
    },
    userPrivModel.PROJECT_EDIT: []string{
        syncd.API_PROJECT_UPDATE,
        syncd.API_PROJECT_EXISTS,
        syncd.API_PROJECT_DETAIL,
    },
    userPrivModel.PROJECT_DEL: []string{
        syncd.API_PROJECT_DELETE,
    },
    userPrivModel.PROJECT_AUDIT: []string{
        syncd.API_PROJECT_STATUS_CHANGE,
    },
    userPrivModel.PROJECT_REPO: []string{
        syncd.API_DEPLOY_REPO_RESET,
    },
    userPrivModel.PROJECT_USER_VIEW: []string{
        syncd.
    },

    // user priv
    userPrivModel.USER_ROLE_VIEW: []string{
        syncd.API_USER_GROUP_LIST,
    },
    userPrivModel.USER_ROLE_NEW: []string{
        syncd.API_USER_GROUP_NEW,
        syncd.API_USER_PRIV_LIST,
    },
    userPrivModel.USER_ROLE_EDIT: []string{
        syncd.API_USER_GROUP_DETAIL,
        syncd.API_USER_GROUP_UPDATE,
    },
    userPrivModel.USER_ROLE_DEL: []string{
        syncd.API_USER_GROUP_DELETE,
    },
    userPrivModel.USER_VIEW: []string{
        syncd.API_USER_LIST,
    },
    userPrivModel.USER_NEW: []string{
        syncd.API_USER_NEW,
        syncd.API_USER_EXISTS,
    },
    userPrivModel.USER_EDIT: []string{
        syncd.API_USER_DETAIL,
        syncd.API_USER_UPDATE,
        syncd.API_USER_GROUP_LIST,
    },
    userPrivModel.USER_DEL: []string{
        syncd.API_USER_DELETE,
    },

    // server priv
    userPrivModel.SERVER_GROUP_VIEW: []string{
        syncd.API_SERVER_GROUP_LIST,
    },
    userPrivModel.SERVER_GROUP_NEW: []string{
        syncd.API_SERVER_GROUP_NEW,
    },
    userPrivModel.SERVER_GROUP_EDIT: []string{
        syncd.API_SERVER_GROUP_DETAIL,
        syncd.API_SERVER_GROUP_UPDATE,
    },
    userPrivModel.SERVER_GROUP_DEL: []string{
        syncd.API_SERVER_GROUP_DELETE,
    },
    userPrivModel.SERVER_VIEW: []string{
        syncd.API_SERVER_LIST,
    },
    userPrivModel.SERVER_NEW: []string{
        syncd.API_SERVER_NEW,
    },
    userPrivModel.SERVER_EDIT: []string{
        syncd.API_SERVER_DETAIL,
        syncd.API_SERVER_UPDATE,
    },
    userPrivModel.SERVER_DEL: []string{
        syncd.API_SERVER_DELETE,
    },
    userPrivModel.SERVER_CHECK: []string{

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
