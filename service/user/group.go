// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "errors"
    "strings"
    "fmt"

    "github.com/tinystack/goutil"
    baseModel "github.com/tinystack/syncd/model"
    groupModel "github.com/tinystack/syncd/model/user/group"
)

type GroupDetail struct {
    ID      int     `json:"id"`
    Name    string  `json:"name"`
    Priv    []int   `json:"priv"`
    Utime   int     `json:"utime"`
}

type GroupItem struct {
    ID      int     `json:"id"`
    Name    string  `json:"name"`
}

func GetUserGroupDetail(id int) (*GroupDetail, error) {
    if id == 0 {
        return nil, errors.New("id can not be empty")
    }

    detail, ok := groupModel.Get(id)
    if !ok {
        return nil, errors.New("get user group detail data failed")
    }

    privList := []int{}
    if detail.Priv != "" {
        strPrivList := goutil.StrFilterSliceEmpty(strings.Split(detail.Priv, ","))
        privList = goutil.StrSlice2IntSlice(strPrivList)
    }

    return &GroupDetail{
        ID: detail.ID,
        Name: detail.Name,
        Priv: privList,
        Utime: detail.Utime,
    }, nil
}

func UpdateUserGroup(id int, name string, priv []string) error {
    var ok bool
    priv = goutil.StrFilterSliceEmpty(priv)
    g := groupModel.UserGroup{
        Name: name,
        Priv: strings.Join(priv, ","),
    }
    if id > 0 {
        ok = groupModel.Update(id, map[string]interface{}{
            "name": g.Name,
            "priv": g.Priv,
        })
    } else {
        ok = groupModel.Create(&g)
    }
    if !ok {
        return errors.New("user group data update failed")
    }

    return nil
}

func GetUserGroupList(keyword string, offset, limit int) ([]GroupItem, int, error) {
    var (
        ok bool
        groupId, total int
        where []baseModel.WhereParam
        groupList []GroupItem
    )
    if keyword != "" {
        if goutil.IsInteger(keyword) {
            groupId = goutil.Str2Int(keyword)
            if groupId > 0 {
                where = append(where, baseModel.WhereParam{
                    Field: "id",
                    Prepare: groupId,
                })
            }
        } else {
            where = append(where, baseModel.WhereParam{
                Field: "name",
                Tag: "LIKE",
                Prepare: fmt.Sprintf("%%%s%%", keyword),
            })
        }
    }
    list, ok := groupModel.List(baseModel.QueryParam{
        Fields: "id, name",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get user group list data failed")
    }
    total, ok = groupModel.Total(baseModel.QueryParam{
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get user group total count failed")
    }

    for _, g := range list {
        groupList = append(groupList, GroupItem{
            ID: g.ID,
            Name: g.Name,
        })
    }
    return groupList, total, nil
}
