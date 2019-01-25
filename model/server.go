// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

import(
    "time"
)

type Server struct {
    ID          int     `gorm:"primary_key"`
    GroupId     int     `gorm:"type:int(11);not null;default:0"`
    Name        string  `gorm:"type:varchar(100);not null;default:''"`
    Ip          string  `gorm:"type:varchar(100);not null;default:''"`
    SSHPort     int     `gorm:"type:int(11);not null;default:0"`
    Ctime       int     `gorm:"type:int(11);not null;default:0"`
}

func (m *Server) TableName() string {
    return "syd_server"
}

func (m *Server) Create() bool {
    m.Ctime = int(time.Now().Unix())
    return Create(m)
}

func (m *Server) Update() bool {
    return UpdateByPk(m)
}

func (m *Server) List(query QueryParam) ([]Server, bool) {
    var data []Server
    ok := GetMulti(&data, query)
    return data, ok
}

func (m *Server) Count(query QueryParam) (int, bool) {
    var count int
    ok := Count(m, &count, query)
    return count, ok
}

func (m *Server) Delete() bool {
    return DeleteByPk(m)
}

func (m *Server) Get(id int) bool {
    return GetByPk(m, id)
}
