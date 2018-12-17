// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "errors"
    "fmt"

    "github.com/tinystack/goutil"
    baseModel "github.com/tinystack/syncd/model"
    serverGroupModel "github.com/tinystack/syncd/model/server/group"
)

type Group struct {
    ID      int         `json:"id"`
    Name    string      `json:"name"`
}

func (g *Group) Detail() error {
    if g.ID == 0 {
        return errors.New("id can not be empty")
    }
    detail, ok := serverGroupModel.Get(g.ID)
    if !ok {
        return errors.New("get server group detail data failed")
    }

    g.ID = detail.ID
    g.Name = detail.Name

    return nil
}

func (g *Group) CreateOrUpdate() error {
    var ok bool
    group := serverGroupModel.ServerGroup{
        Name: g.Name,
    }
    if g.ID > 0 {
        ok = serverGroupModel.Update(g.ID, group)
    } else {
        ok = serverGroupModel.Create(&group)
    }
    if !ok {
        return errors.New("server group data update or create failed")
    }
    return nil
}

func (g *Group) GetMultiById(ids []int) ([]Group, error) {
    var where []baseModel.WhereParam
    if len(ids) > 0 {
        where = append(where, baseModel.WhereParam{
            Field: "id",
            Tag: "IN",
            Prepare: ids,
        })
    }
    list, ok := serverGroupModel.List(baseModel.QueryParam{
        Fields: "id, name",
        Where: where,
    })
    if !ok {
        return nil, errors.New("get server group list data failed")
    }
    var nlist []Group
    for _, l := range list {
        nlist = append(nlist, Group{
            ID: l.ID,
            Name: l.Name,
        })
    }
    return nlist, nil
}

func (g *Group) List(keyword string, offset, limit int) ([]Group, int, error){
    var (
        ok bool
        groupId, total int
        where []baseModel.WhereParam
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
    list, ok := serverGroupModel.List(baseModel.QueryParam{
        Fields: "id, name",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get server group list data failed")
    }
    total, ok = serverGroupModel.Total(baseModel.QueryParam{
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
    ok := serverGroupModel.Delete(g.ID)
    if !ok {
        return errors.New("delete server group data failed")
    }
    return nil
}
