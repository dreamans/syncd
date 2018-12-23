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
    userGroupModel "github.com/tinystack/syncd/model/user/group"
)

type Group struct {
    ID      int     `json:"id"`
    Name    string  `json:"name"`
    Priv    []int   `json:"priv"`
    Utime   int     `json:"utime"`
}

type GroupItem struct {
    ID      int     `json:"id"`
    Name    string  `json:"name"`
}

func (g *Group) Get() error {
    if g.ID == 0 {
        return errors.New("id can not be empty")
    }
    detail, ok := userGroupModel.Get(g.ID)
    if !ok {
        return errors.New("get user group detail data failed")
    }

    privList := []int{}
    if detail.Priv != "" {
        strPrivList := goutil.StrFilterSliceEmpty(strings.Split(detail.Priv, ","))
        privList = goutil.StrSlice2IntSlice(strPrivList)
    }

    g.ID = detail.ID
    g.Name = detail.Name
    g.Priv = privList
    g.Utime = detail.Utime

    return nil
}

func (g *Group) CreateOrUpdate() error {
    var ok bool
    group := userGroupModel.UserGroup{
        Name: g.Name,
        Priv: strings.Join(goutil.IntSlice2StrSlice(g.Priv), ","),
    }
    if g.ID > 0 {
        ok = userGroupModel.Update(g.ID, map[string]interface{}{
            "name": group.Name,
            "priv": group.Priv,
        })
    } else {
        ok = userGroupModel.Create(&group)
    }
    if !ok {
        return errors.New("user group data update failed")
    }
    return nil
}

func (g *Group) List(keyword string, offset, limit int) ([]GroupItem, int, error) {
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
    list, ok := userGroupModel.List(baseModel.QueryParam{
        Fields: "id, name",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get user group list data failed")
    }
    total, ok = userGroupModel.Total(baseModel.QueryParam{
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

func (g *Group) Delete() error {
    if g.ID == 0 {
        return errors.New("id can not be empty")
    }
    ok := userGroupModel.Delete(g.ID)
    if !ok {
        return errors.New("user group delete failed")
    }
    return nil
}

func (g *Group) GetNameByIds(ids []int) (map[int]string, error){
    list, ok := userGroupModel.List(baseModel.QueryParam{
        Fields: "id, name",
        Where: []baseModel.WhereParam{
            baseModel.WhereParam{
                Field: "id",
                Tag: "IN",
                Prepare: ids,
            },
        },
    })
    if !ok {
        return nil, errors.New("get user group list failed")
    }
    groupNameList := make(map[int]string)
    for _, g := range list {
        groupNameList[g.ID] = g.Name
    }
    return groupNameList, nil
}
