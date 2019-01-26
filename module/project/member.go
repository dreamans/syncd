// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "errors"
    //"fmt"

    "github.com/dreamans/syncd/model"
)

type Member struct {
    ID          int     `json:"id"`
    UserId      int     `json:"user_id"`
    SpaceId     int     `json:"space_id"`
    Username    string  `json:"username"`
    Email       string  `json:"email"`
    RoleName    string  `json:"role_name"`
    Status      int     `json:"status"`
    Ctime       int     `json:"ctime"`
}

func (m *Member) List(spaceId, offset, limit int) ([]Member, error) {
    member := &model.ProjectMember{}
    list, ok := member.List(model.QueryParam{
        Fields: "id, space_id, user_id, ctime",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
    })
    if !ok {
        return nil, errors.New("get project member list failed")
    }

    var (
        memberList []Member
        userIdList []int
    )
    for _, l := range list {
        userIdList = append(userIdList, l.UserId)
        memberList = append(memberList, Member{
            ID: l.ID,
            UserId: l.UserId,
            SpaceId: l.SpaceId,
            Ctime: l.Ctime,
        })
    }
    user := &model.User{}
    userList, ok := user.List(model.QueryParam{
        Fields: "id, role_id, username, email, status",
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Tag: "IN",
                Prepare: userIdList,
            },
        },
    })
    if !ok {
        return nil, errors.New("get project user detail list failed")
    }
    userMap := make(map[int]model.User)
    var roleIdList []int
    for _, l := range userList {
        userMap[l.ID] = l
        roleIdList = append(roleIdList, l.RoleId)
    }
    role := &model.UserRole{}
    roleList, ok := role.List(model.QueryParam{
        Fields: "id, role_id, username, email, status",
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Tag: "IN",
                Prepare: userIdList,
            },
        },
    })
    if !ok {
        return nil, errors.New("get project user role list failed")
    }

    return memberList, nil
}

func (m *Member) Exists() (bool, error) {
    where := []model.WhereParam{
        model.WhereParam{
            Field: "user_id",
            Prepare: m.UserId,
        },
        model.WhereParam{
            Field: "space_id",
            Prepare: m.SpaceId,
        },
    }
	member := &model.ProjectMember{}
    count, ok := member.Count(model.QueryParam{
        Where: where,
    })
    if !ok {
        return false, errors.New("check member exists failed")
    }
    return count > 0, nil
}

func (m *Member) Create() error {
    member := &model.ProjectMember{
        SpaceId: m.SpaceId,
        UserId: m.UserId,
    }
    if ok := member.Create(); !ok {
        return errors.New("create project member failed")
    }
    return nil
}