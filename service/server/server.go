// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "errors"
    "fmt"

    "github.com/tinystack/goutil/gois"
    "github.com/tinystack/goutil/gostring"
    "github.com/dreamans/syncd/model"
    serverModel "github.com/dreamans/syncd/model/server"
)

type Server struct {
    ID      int     `json:"id"`
    GroupId int     `json:"group_id"`
    Name    string  `json:"name"`
    Ip      string  `json:"ip"`
    SshPort int     `json:"ssh_port"`
}

func ServerGetListByGroupIds(ids []int) ([]Server, error){
    list, ok := serverModel.List(model.QueryParam{
        Fields: "id, group_id, name, ip, ssh_port",
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "group_id",
                Tag: "IN",
                Prepare: ids,
            },
        },
    })
    if !ok {
        return nil, errors.New("get server list data failed")
    }
    var serList []Server
    for _, l := range list {
        serList = append(serList, Server{
            ID: l.ID,
            GroupId: l.GroupId,
            Ip: l.Ip,
            SshPort: l.SshPort,
            Name: l.Name,
        })
    }
    return serList, nil
}

func (s *Server) CreateOrUpdate() error {
    server := serverModel.Server{
        GroupId: s.GroupId,
        Name: s.Name,
        Ip: s.Ip,
        SshPort: s.SshPort,
    }
    var ok bool
    if s.ID > 0 {
        ok = serverModel.Update(s.ID, server)
    } else {
        ok = serverModel.Create(&server)
    }
    if !ok {
        return errors.New("server data update failed")
    }
    return nil
}

func (g *Server) List(keyword string, groupId, offset, limit int) ([]Server, int, error){
    var (
        serverId int
        where []model.WhereParam
    )
    if keyword != "" {
        if gois.IsInteger(keyword) {
            serverId = gostring.Str2Int(keyword)
            if serverId > 0 {
                where = append(where, model.WhereParam{
                    Field: "id",
                    Prepare: serverId,
                })
            }
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

    if groupId > 0 {
        where = append(where, model.WhereParam{
            Field: "group_id",
            Prepare: groupId,
        })
    }

    list, ok := serverModel.List(model.QueryParam{
        Fields: "id, group_id, name, ip, ssh_port",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get server list data failed")
    }

    total, ok := serverModel.Total(model.QueryParam{
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get server total count failed")
    }

    var nlist []Server
    for _, l := range list {
        nlist = append(nlist, Server{
            ID: l.ID,
            GroupId: l.GroupId,
            Name: l.Name,
            Ip: l.Ip,
            SshPort: l.SshPort,
        })
    }
    return nlist, total, nil
}

func (s *Server) Get() error {
    if s.ID == 0 {
        return errors.New("id can not be empty")
    }
    detail, ok := serverModel.Get(s.ID)
    if !ok {
        return errors.New("get server detail data failed")
    }
    if detail.ID == 0 {
        return errors.New("server detail data not exists")
    }
    s.GroupId = detail.GroupId
    s.Name = detail.Name
    s.Ip = detail.Ip
    s.SshPort = detail.SshPort
    return nil
}

func (s *Server) Delete() error {
    if s.ID == 0 {
        return errors.New("id can not be empty")
    }
    if err := s.Get(); err != nil {
        return err
    }
    ok := serverModel.Delete(s.ID)
    if !ok {
        return errors.New("delete server data failed")
    }
    return nil
}
