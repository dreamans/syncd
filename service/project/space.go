// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "errors"
    "fmt"

    baseModel "github.com/tinystack/syncd/model"
    projectSpaceModel "github.com/tinystack/syncd/model/project_space"
    projectUserModel "github.com/tinystack/syncd/model/project_user"
)

type Space struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Ctime       int     `json:"ctime"`
}

func SpaceGetByPk(id int) (*Space, error) {
    space := &Space{
        ID: id,
    }
    if err := space.Detail(); err != nil {
        return nil, err
    }
    return space, nil
}

func SpaceGetMapByIds(ids []int) (map[int]Space, error) {
    list, err := SpaceGetListByIds(ids)
    if err != nil {
        return nil, err
    }
    maps := map[int]Space{}
    for _, l := range list {
        maps[l.ID] = l
    }
    return maps, nil
}

func SpaceGetListByIds(ids []int) ([]Space, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    list, ok := projectSpaceModel.List(baseModel.QueryParam{
        Fields: "id, name",
        Order: "id DESC",
        Where: []baseModel.WhereParam{
            baseModel.WhereParam{
                Field: "id",
                Tag: "IN",
                Prepare: ids,
            },
        },
    })
    if !ok {
        return nil, errors.New("get project space list failed")
    }
    var spaceList []Space
    for _, l := range list {
        spaceList = append(spaceList, Space{
            ID: l.ID,
            Name: l.Name,
        })
    }
    return spaceList, nil
}

func SpaceGetIdListByUserId(userId int) ([]int, error) {
    spaceList, err := SpaceGetListByUserId(userId)
    if err != nil {
        return nil, err
    }
    var idList []int
    for _, l := range spaceList {
        idList = append(idList, l.ID)
    }
    return idList, nil
}

func SpaceGetListByUserId(userId int) ([]Space, error) {
    spaceIds, ok := projectUserModel.GetSpaceIdsByUserId(userId)
    if !ok {
        return nil, errors.New("get space ids failed")
    }
    list, ok := projectSpaceModel.List(baseModel.QueryParam{
        Where: []baseModel.WhereParam{
            baseModel.WhereParam{
                Field: "id",
                Tag: "IN",
                Prepare: spaceIds,
            },
        },
    })

    var spaceList []Space
    for _, l := range list {
        spaceList = append(spaceList, Space{
            ID: l.ID,
            Name: l.Name,
        })
    }

    return spaceList, nil
}

func (s *Space) CreateOrUpdate() error {
    var ok bool
    space := projectSpaceModel.ProjectSpace{
        Name: s.Name,
        Description: s.Description,
    }
    if s.ID > 0 {
        ok = projectSpaceModel.Update(s.ID, space)
    } else {
        ok = projectSpaceModel.Create(&space)
    }
    if !ok {
        return errors.New("project space update or create failed")
    }
    return nil
}

func (s *Space) List(keyword string, offset, limit int) ([]Space, int, error) {
    var where []baseModel.WhereParam
    if keyword != "" {
        where = append(where, baseModel.WhereParam{
            Field: "name",
            Tag: "LIKE",
            Prepare: fmt.Sprintf("%%%s%%", keyword),
        })
    }
    list, ok := projectSpaceModel.List(baseModel.QueryParam{
        Fields: "id, name, description, ctime",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get project space list failed")
    }
    total, ok := projectSpaceModel.Total(baseModel.QueryParam{
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get project space total count failed")
    }
    var nlist []Space
    for _, l := range list {
        nlist = append(nlist, Space{
            ID: l.ID,
            Name: l.Name,
            Description: l.Description,
            Ctime: l.Ctime,
        })
    }
    return nlist, total, nil
}

func (s *Space) Detail() error {
    if s.ID == 0 {
        return errors.New("space id not exists")
    }
    detail, ok := projectSpaceModel.Get(s.ID)
    if !ok {
        return errors.New("get project space failed")
    }
    if detail.ID == 0 {
        return errors.New("space not exists")
    }
    s.Name = detail.Name
    s.Description = detail.Description
    s.Ctime = detail.Ctime
    return nil
}

func (s *Space) Delete() error {
    if s.ID == 0 {
        return errors.New("id can not be empty")
    }
    ok := projectSpaceModel.Delete(s.ID)
    if !ok {
        return errors.New("project space delete failed")
    }
    return nil
}

func (s *Space) CheckExists() (bool, error) {
    var where []baseModel.WhereParam
    if s.Name != "" {
        where = append(where, baseModel.WhereParam{
            Field: "name",
            Prepare: s.Name,
        })
    }
    if s.ID > 0 {
        where = append(where, baseModel.WhereParam{
            Field: "id",
            Tag: "!=",
            Prepare: s.ID,
        })
    }
    detail, ok := projectSpaceModel.GetOne(baseModel.QueryParam{
        Where: where,
    })
    if !ok {
        return false, errors.New("get project space one data failed")
    }
    return detail.ID > 0, nil
}

