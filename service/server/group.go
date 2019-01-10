
// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "errors"
    "fmt"

    "github.com/tinystack/goutil/gostring"
    "github.com/tinystack/goutil/gois"
    "github.com/dreamans/syncd/model"
    serverGroupModel "github.com/dreamans/syncd/model/server_group"
)

type Group struct {
    ID      int         `json:"id"`
    Name    string      `json:"name"`
}

func GroupListByIds(ids []int) ([]Group, error) {
    list, ok := serverGroupModel.List(model.QueryParam{
        Fields: "id, name",
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Tag: "IN",
                Prepare: ids,
            },
        },
    })
    if !ok {
        return nil, errors.New("get group list failed")
    }
    var groupList []Group
    for _, l := range list {
        groupList = append(groupList, Group{
            ID: l.ID,
            Name: l.Name,
        })
    }
    return groupList, nil
}

func (g *Group) CreateOrUpdate() error {
    var ok bool
    serverGroup := serverGroupModel.ServerGroup{
        Name: g.Name,
    }
    if g.ID > 0 {
        ok = serverGroupModel.Update(g.ID, serverGroup)
    } else {
        ok = serverGroupModel.Create(&serverGroup)
    }
    if !ok {
        return errors.New("server group data update or create failed")
    }
    return nil
}

func (g *Group) Detail() error {
    if g.ID == 0 {
        return errors.New("id can not be empty")
    }
    detail, ok := serverGroupModel.Get(g.ID)
    if !ok {
        return errors.New("get server group detail data failed")
    }
    if detail.ID == 0 {
        return errors.New("server group detail not exists")
    }
    g.ID = detail.ID
    g.Name = detail.Name
    return nil
}

func (g *Group) List(keyword string, offset, limit int) ([]Group, int, error){
    var (
        ok bool
        groupId, total int
        where []model.WhereParam
    )
    if keyword != "" {
        if gois.IsInteger(keyword) {
            groupId = gostring.Str2Int(keyword)
            if groupId > 0 {
                where = append(where, model.WhereParam{
                    Field: "id",
                    Prepare: groupId,
                })
            }
        } else {
            where = append(where, model.WhereParam{
                Field: "name",
                Tag: "LIKE",
                Prepare: fmt.Sprintf("%%%s%%", keyword),
            })
        }
    }
    list, ok := serverGroupModel.List(model.QueryParam{
        Fields: "id, name",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get server group list data failed")
    }
    total, ok = serverGroupModel.Total(model.QueryParam{
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get server group total count failed")
    }
    var nlist []Group
    for _, l := range list {
        nlist = append(nlist, Group{
            ID: l.ID,
            Name: l.Name,
        })
    }
    return nlist, total, nil
}

func (g *Group) Delete() error {
    if g.ID == 0 {
        return errors.New("id can not be empty")
    }
    if err := g.Detail(); err != nil {
        return err
    }
    ok := serverGroupModel.Delete(g.ID)
    if !ok {
        return errors.New("delete server group data failed")
    }
    return nil
}
