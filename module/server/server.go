// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "errors"
    "fmt"

    "github.com/dreamans/syncd/util/gois"
    "github.com/dreamans/syncd/model"
)

type Server struct {
    ID          int     `json:"id"`
    GroupId     int     `json:"group_id"`
    GroupName   string  `json:"group_name"`
    Name        string  `json:"name"`
    Ip          string  `json:"ip"`
    SSHPort     int     `json:"ssh_port"`
    Ctime       int     `json:"ctime"`
}

func ServerGetListByGroupIds(groupIds []int) ([]Server, error){
    server := &model.Server{}
    list, ok := server.List(model.QueryParam{
        Fields: "id, group_id, name, ip, ssh_port, ctime",
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "group_id",
                Tag: "IN",
                Prepare: groupIds,
            },
        },
    })
    if !ok {
        return nil, errors.New("get server list failed")
    }
    serList := []Server{}
    for _, s := range list {
        serList = append(serList, Server{
            ID: s.ID,
            GroupId: s.GroupId,
            Name: s.Name,
            Ip: s.Ip,
            SSHPort: s.SSHPort,
        })
    }
    return serList, nil
}

func (s *Server) CreateOrUpdate() error {
    server := &model.Server{
        ID: s.ID,
        GroupId: s.GroupId,
        Name: s.Name,
        Ip: s.Ip,
        SSHPort: s.SSHPort,
    }
    if server.ID == 0 {
        if ok := server.Create(); !ok {
            return errors.New("create server failed")
        }
    } else {
        if ok := server.Update(); !ok {
            return errors.New("update server failed")
        }
    }
    return nil
}

func (s *Server) List(keyword string, offset, limit int) ([]Server, error) {
    server := &model.Server{}
    list, ok := server.List(model.QueryParam{
        Fields: "id, group_id, name, ip, ssh_port, ctime",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: s.parseWhereConds(keyword),
    })
    if !ok {
        return nil, errors.New("get server list failed")
    }

    var (
        serverList []Server
        groupIds []int
    )
    for _, l := range list {
        serverList = append(serverList, Server{
            ID: l.ID,
            GroupId: l.GroupId,
            Name: l.Name,
            Ip: l.Ip,
            SSHPort: l.SSHPort,
            Ctime: l.Ctime,
        })
        groupIds = append(groupIds, l.GroupId)
    }

    groupMap, err := GroupGetMapByIds(groupIds)
    if err != nil {
        return nil, errors.New("get server group map failed")
    }
    for k, l := range serverList {
        if g, exists := groupMap[l.GroupId]; exists {
            serverList[k].GroupName = g.Name
        }
    }

    return serverList, nil
}

func (s *Server) Total(keyword string) (int, error) {
    server := model.Server{}
    total, ok := server.Count(model.QueryParam{
        Where: s.parseWhereConds(keyword),
    })
    if !ok {
        return 0, errors.New("get server count failed")
    }
    return total, nil
}

func (s *Server) Delete() error {
    server := &model.Server{
        ID: s.ID,
    }
    if ok := server.Delete(); !ok {
        return errors.New("delete server failed")
    }
    return nil
}

func (s *Server) Detail() error {
    server := &model.Server{}
    if ok := server.Get(s.ID); !ok {
        return errors.New("get server detail failed")
    }
    if server.ID == 0 {
        return errors.New("server not exists")
    }

    s.ID = server.ID
    s.GroupId = server.GroupId
    s.Name = server.Name
    s.Ip = server.Ip
    s.SSHPort = server.SSHPort
    s.Ctime = server.Ctime

    return nil
}

func (s *Server) parseWhereConds(keyword string) []model.WhereParam {
    var where []model.WhereParam
    if keyword != "" {
        if gois.IsInteger(keyword) {
            where = append(where, model.WhereParam{
                Field: "id",
                Prepare: keyword,
            })
        } else {
            if gois.IsIp(keyword) {
                where = append(where, model.WhereParam{
                    Field: "ip",
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
    }
    return where
}
