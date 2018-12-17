// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "errors"
    "fmt"

    "github.com/tinystack/goutil"
    baseModel "github.com/tinystack/syncd/model"
    serverModel "github.com/tinystack/syncd/model/server"
)

type Server struct {
    ID      int     `json:"id"`
    GroupId int     `json:"group_id"`
    Name    string  `json:"name"`
    Ip      string  `json:"ip"`
    SshPort int     `json:"ssh_port"`
}

func (s *Server) CreateOrUpdate() error {
    smodel := serverModel.Server{
        GroupId: s.GroupId,
        Name: s.Name,
        Ip: s.Ip,
        SshPort: s.SshPort,
    }
    var ok bool
    if s.ID > 0 {
        ok = serverModel.Update(s.ID, smodel)
    } else {
        ok = serverModel.Create(&smodel)
    }
    if !ok {
        return errors.New("server data update failed")
    }
    return nil
}

func (g *Server) List(keyword string, groupId, offset, limit int) ([]Server, int, error){
    var (
        serverId int
        where []baseModel.WhereParam
    )
    if keyword != "" {
        if goutil.IsInteger(keyword) {
            serverId = goutil.Str2Int(keyword)
            if serverId > 0 {
                where = append(where, baseModel.WhereParam{
                    Field: "id",
                    Prepare: serverId,
                })
            }
        } else {
            if goutil.IsIp(keyword) {
                where = append(where, baseModel.WhereParam{
                    Field: "ip",
                    Prepare: keyword,
                })
            } else {
                where = append(where, baseModel.WhereParam{
                    Field: "name",
                    Tag: "LIKE",
                    Prepare: fmt.Sprintf("%%%s%%", keyword),
                })
            }
        }
    }

    if groupId > 0 {
        where = append(where, baseModel.WhereParam{
            Field: "group_id",
            Prepare: groupId,
        })
    }

    list, ok := serverModel.List(baseModel.QueryParam{
        Fields: "id, group_id, name, ip, ssh_port",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get server list data failed")
    }

    total, ok := serverModel.Total(baseModel.QueryParam{
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
    ok := serverModel.Delete(s.ID)
    if !ok {
        return errors.New("delete server data failed")
    }
    return nil
}
