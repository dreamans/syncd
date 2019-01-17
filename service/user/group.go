// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "errors"
    "strings"
    "fmt"

    "github.com/tinystack/goutil/gostring"
    "github.com/tinystack/goutil/gois"
    baseModel "github.com/dreamans/syncd/model"
    userGroupModel "github.com/dreamans/syncd/model/user_group"
)

type Group struct {
    ID      int     `json:"id"`
    Name    string  `json:"name"`
    Priv    []int   `json:"priv"`
    Ctime   int     `json:"ctime"`
}

type GroupItem struct {
    ID      int     `json:"id"`
    Name    string  `json:"name"`
}

func GroupUserListFillGroupName(list []UserItem) ([]UserItem, error) {
    var groupIdList []int
    for _, l := range list {
        groupIdList = append(groupIdList, l.GroupId)
    }
    if len(groupIdList) > 0 {
        group := &Group{}
        groupNameList, err := group.GetNameByIds(groupIdList)
        if err != nil {
            return nil, err
        }
        for k, v := range list {
            if groupName, exists := groupNameList[v.GroupId]; exists {
                list[k].GroupName = groupName
            }
        }
    }
    return list, nil
}

func (g *Group) Detail() error {
    if g.ID == 0 {
        return errors.New("id can not be empty")
    }
    detail, ok := userGroupModel.Get(g.ID)
    if !ok {
        return errors.New("get user group detail data failed")
    }
    if detail.ID == 0 {
        return errors.New("user group not exists")
    }
    privList := []int{}
    if detail.Priv != "" {
        strPrivList := gostring.StrFilterSliceEmpty(strings.Split(detail.Priv, ","))
        privList = gostring.StrSlice2IntSlice(strPrivList)
    }
    g.ID = detail.ID
    g.Name = detail.Name
    g.Priv = privList
    g.Ctime = detail.Ctime
    return nil
}

func (g *Group) CreateOrUpdate() error {
    var ok bool
    group := userGroupModel.UserGroup{
        Name: g.Name,
        Priv: strings.Join(gostring.IntSlice2StrSlice(g.Priv), ","),
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
        if gois.IsInteger(keyword) {
            groupId = gostring.Str2Int(keyword)
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
    if err := g.Detail(); err != nil {
        return errors.New("user group not exists")
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

func (g *Group) CheckGroupExists() (bool, error){
    var where []baseModel.WhereParam
	if g.Name != "" {
		where = append(where, baseModel.WhereParam{
			Field: "name",
			Prepare: g.Name,
		})
	}
	if g.ID > 0 {
		where = append(where, baseModel.WhereParam{
			Field: "id",
			Tag: "!=",
			Prepare: g.ID,
		})
	}
    detail, ok := userGroupModel.GetOne(baseModel.QueryParam{
    	Where: where,
	})
	if !ok {
		return false, errors.New("get group one data failed")
	}
    return detail.ID > 0, nil
}