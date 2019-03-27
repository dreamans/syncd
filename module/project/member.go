// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "errors"

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

func (m *Member) Detail() error {
    member := model.ProjectMember{}
    if ok := member.Get(m.ID); !ok {
        return errors.New("get project member detail failed")
    }
    m.ID = member.ID
    m.UserId = member.UserId
    m.SpaceId = member.SpaceId
    return nil
}

func (m *Member) MemberInSpace() bool {
    exists , _ := m.Exists()
    return exists
}

func (m *Member) Delete() error {
    member := &model.ProjectMember{
        ID: m.ID,
    }
    if ok := member.Delete(); !ok {
        return errors.New("remove project member failed")
    }
    return nil
}

func (m *Member) Total(spaceId int) (int, error) {
    member := &model.ProjectMember{}
    total, ok := member.Count(model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "space_id",
                Prepare: spaceId,
            },
        },
    })
    if !ok {
        return 0, errors.New("get project member count failed")
    }
    return total, nil
}

func (m *Member) SpaceIdsByUserId() ([]int, error) {
    member := &model.ProjectMember{}
    list, ok := member.List(model.QueryParam{
        Fields: "space_id",
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "user_id",
                Prepare: m.UserId,
            },
        },
    })
    if !ok {
        return nil, errors.New("get project member list failed")
    }
    var spaceIds []int
    for _, l := range list {
        spaceIds = append(spaceIds, l.SpaceId)
    }

    return spaceIds, nil
}

func (m *Member) List(spaceId, offset, limit int) ([]Member, error) {
    member := &model.ProjectMember{}
    list, ok := member.List(model.QueryParam{
        Fields: "id, space_id, user_id, ctime",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "space_id",
                Prepare: spaceId,
            },
        },
    })
    if !ok {
        return nil, errors.New("get project member list failed")
    }

    var (
        memberList []Member
        userIdList, roleIdList []int
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
    for _, l := range userList {
        userMap[l.ID] = l
        roleIdList = append(roleIdList, l.RoleId)
    }
    role := &model.UserRole{}
    roleList, ok := role.List(model.QueryParam{
        Fields: "id, name",
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Tag: "IN",
                Prepare: roleIdList,
            },
        },
    })
    if !ok {
        return nil, errors.New("get project user role list failed")
    }
    roleMap := make(map[int]model.UserRole)
    for _, l := range roleList {
        roleMap[l.ID] = l
    }

    for k, m := range memberList {
        if u, ok := userMap[m.UserId]; ok {
            memberList[k].Username = u.Username
            memberList[k].Email = u.Email
            memberList[k].Status = u.Status
            if r, ok := roleMap[u.RoleId]; ok {
                memberList[k].RoleName = r.Name
            }
        }

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
