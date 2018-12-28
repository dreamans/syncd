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
        syncd.API_PROJECT_SPACE_LIST,
        syncd.API_PROJECT_LIST,
        syncd.API_PROJECT_DETAIL,
        syncd.API_DEPLOY_REPO_COMMITLIST,
        syncd.API_DEPLOY_REPO_TAGLIST,
        syncd.API_DEPLOY_APPLY_SUBMIT,
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
