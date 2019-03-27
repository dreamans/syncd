// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "errors"
    "fmt"

    "github.com/dreamans/syncd/util/gostring"
    "github.com/dreamans/syncd/util/gois"
    "github.com/dreamans/syncd/model"
)

type Role struct {
    ID          int         `json:"id"`
    Name        string      `json:"name"`
    Privilege   []int       `json:"privilege"`
    Ctime       int         `json:"ctime"`
}

func RoleGetMapByIds(ids []int) (map[int]Role, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    role := &model.UserRole{}
    list, ok := role.List(model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Tag: "IN",
                Prepare: ids,
            },
        },
    })
    if !ok {
        return nil, errors.New("get user role maps failed")
    }
    roleMap := make(map[int]Role)
    for _, l := range list {
        roleMap[l.ID] = Role{
            ID: l.ID,
            Name: l.Name,
        }
    }
    return roleMap, nil
}

func (r *Role) Detail() error {
    role := &model.UserRole{}
    if ok := role.Get(r.ID); !ok {
        return errors.New("get user role detail failed")
    }
    if role.ID == 0 {
        return errors.New("user role not exists")
    }

    r.ID = role.ID
    r.Name = role.Name
    r.Privilege = gostring.StrSplit2IntSlice(role.Privilege, ",")
    r.Ctime = role.Ctime

    return nil
}

func (r *Role) CreateOrUpdate() error {
    role := &model.UserRole{
        ID: r.ID,
        Name: r.Name,
        Privilege: gostring.JoinIntSlice2String(r.Privilege, ","),
    }
    if role.ID == 0 {
        if ok := role.Create(); !ok {
            return errors.New("create user role data failed")
        }
    } else {
        if ok := role.Update(); !ok {
            return errors.New("update user role failed")
        }
    }
    return nil
}

func (r *Role) List(keyword string, offset, limit int) ([]Role, error) {
    role := &model.UserRole{}
    list, ok := role.List(model.QueryParam{
        Fields: "id, name, ctime",
        Offset: offset,
        Limit: limit,
        Order: "id ASC",
        Where: r.parseWhereConds(keyword),
    })
    if !ok {
        return nil, errors.New("get user role list failed")
    }

    var roleList []Role
    for _, l := range list {
        roleList = append(roleList, Role{
            ID: l.ID,
            Name: l.Name,
            Ctime: l.Ctime,
        })
    }
    return roleList, nil
}

func (r *Role) Total(keyword string) (int, error) {
    role := &model.UserRole{}
    total, ok := role.Count(model.QueryParam{
        Where: r.parseWhereConds(keyword),
    })
    if !ok {
        return 0, errors.New("get user role count failed")
    }
    return total, nil
}

func (r *Role) Delete() error {
    role := &model.UserRole{
        ID: r.ID,
    }
    if ok := role.Delete(); !ok {
        return errors.New("delete user role failed")
    }
    return nil
}

func (r *Role) parseWhereConds(keyword string) []model.WhereParam {
    var where []model.WhereParam
    if keyword != "" {
        if gois.IsInteger(keyword) {
            where = append(where, model.WhereParam{
                Field: "id",
                Prepare: keyword,
            })
        } else {
            where = append(where, model.WhereParam{
                Field: "name",
                Tag: "LIKE",
                Prepare: fmt.Sprintf("%%%s%%", keyword),
            })
        }
    }
    return where
}
