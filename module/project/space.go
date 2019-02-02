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

func SpaceListByIds(spaceIds []int) ([]Space, error) {
    s := &Space{}
    return s.List(spaceIds, "", 0, 999)
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
        s.ID = space.ID
    } else {
        if ok := space.Update(); !ok {
            return errors.New("project space update failed")
        }
    }
    return nil
}

func (s *Space) List(spaceIds []int, keyword string, offset, limit int) ([]Space, error) {
    space := &model.ProjectSpace{}
    list, ok := space.List(model.QueryParam{
        Fields: "id, name, description, ctime",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: s.parseWhereConds(spaceIds, keyword),
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

func (s *Space) Total(spaceIds []int, keyword string) (int, error) {
    space := &model.ProjectSpace{}
    total, ok := space.Count(model.QueryParam{
        Where: s.parseWhereConds(spaceIds, keyword),
    })
    if !ok {
        return 0, errors.New("get project space count failed")
    }
    return total, nil
}

func (s *Space) parseWhereConds(spaceIds []int, keyword string) []model.WhereParam {
    var where []model.WhereParam
    if keyword != "" {
        where = append(where, model.WhereParam{
            Field: "name",
            Tag: "LIKE",
            Prepare: fmt.Sprintf("%%%s%%", keyword),
        })
    }
    where = append(where, model.WhereParam{
        Field: "id",
        Tag: "IN",
        Prepare: spaceIds,
    })
    return where
}
