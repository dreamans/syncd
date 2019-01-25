// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "errors"
    "fmt"

    "github.com/dreamans/syncd/model"
)

type Space struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Ctime       int     `json:"ctime"`
}

func (s *Space) Delete() error {
    space := &model.ProjectSpace{
        ID: s.ID,
    }
    if ok := space.Delete(); !ok {
        return errors.New("delete project space failed")
    }
    return nil
}

func (s *Space) Detail() error {
    space := &model.ProjectSpace{}
    if ok := space.Get(s.ID); !ok {
        return errors.New("get project space detail failed")
    }
    if space.ID == 0 {
        return errors.New("project space detail not exists")
    }

    s.ID = space.ID
    s.Name = space.Name
    s.Description = space.Description
    s.Ctime = space.Ctime

    return nil
}

func (s *Space) CreateOrUpdate() error {
    space := &model.ProjectSpace{
        ID: s.ID,
        Name: s.Name,
        Description: s.Description,
    }
    if space.ID == 0 {
        if ok := space.Create(); !ok {
            return errors.New("project space create failed")
        }
    } else {
        if ok := space.Update(); !ok {
            return errors.New("project space update failed")
        }
    }
    return nil
}

func (s *Space) List(keyword string, offset, limit int) ([]Space, error) {
    space := &model.ProjectSpace{}
    list, ok := space.List(model.QueryParam{
        Fields: "id, name, description, ctime",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: s.parseWhereConds(keyword),
    })
    if !ok {
        return nil, errors.New("get project space list failed")
    }

    var spaceList []Space
    for _, l := range list {
        spaceList = append(spaceList, Space{
            ID: l.ID,
            Name: l.Name,
            Description: l.Description,
            Ctime: l.Ctime,
        })
    }
    return spaceList, nil
}

func (s *Space) Total(keyword string) (int, error) {
    space := &model.ProjectSpace{}
    total, ok := space.Count(model.QueryParam{
        Where: s.parseWhereConds(keyword),
    })
    if !ok {
        return 0, errors.New("get project space count failed")
    }
    return total, nil
}

func (s *Space) parseWhereConds(keyword string) []model.WhereParam {
    var where []model.WhereParam
    if keyword != "" {
        where = append(where, model.WhereParam{
            Field: "name",
            Tag: "LIKE",
            Prepare: fmt.Sprintf("%%%s%%", keyword),
        })
    }
    return where
}

/*
func GroupGetMapByIds(ids []int) (map[int]Group, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    group := &model.ServerGroup{}
    groupList, ok := group.List(model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Tag: "IN",
                Prepare: ids,
            },
        },
    })
    if !ok {
        return nil, errors.New("get server group list failed")
    }
    groupMap := make(map[int]Group)
    for _, l := range groupList{
        groupMap[l.ID] = Group{
            ID: l.ID,
            Name: l.Name,
            Ctime: l.Ctime,
        }
    }
    return groupMap, nil
}

func (g *Group) Create() error {
    serverGroup := model.ServerGroup{
        Name: g.Name,
    }
    if ok := serverGroup.Create(); !ok {
        return errors.New("create server group data failed")
    }
    return nil
}

func (g *Group) Update() error {
    serverGroup := model.ServerGroup{
        ID: g.ID,
        Name: g.Name,
    }
    if ok := serverGroup.Update(); !ok {
        return errors.New("update server group data failed")
    }
    return nil
}

func (g *Group) List(keyword string, offset, limit int) ([]Group, error) {
    group := model.ServerGroup{}
    list, ok := group.List(model.QueryParam{
        Fields: "id, name, ctime",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: g.parseWhereConds(keyword),
    })
    if !ok {
        return nil, errors.New("get server group list failed")
    }

    var groupList []Group
    for _, l := range list {
        groupList = append(groupList, Group{
            ID: l.ID,
            Name: l.Name,
            Ctime: l.Ctime,
        })
    }
    return groupList, nil
}



func (g *Group) Delete() error {
    group := &model.ServerGroup{
        ID: g.ID,
    }
    if ok := group.Delete(); !ok {
        return errors.New("delete server group failed")
    }
    return nil
}

func (g *Group) Detail() error {
    group := model.ServerGroup{}
    if ok := group.Get(g.ID); !ok {
        return errors.New("get server group detail failed")
    }
    if group.ID == 0 {
        return errors.New("server group not exists")
    }

    g.ID = group.ID
    g.Name = group.Name
    g.Ctime = group.Ctime

    return nil
}


*/