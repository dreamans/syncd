// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

import (
    "time"
)

type ServerGroup struct {
    ID      int         `gorm:"primary_key"`
    Name    string      `gorm:"type:varchar(100);not null;default:''"`
    Ctime   int         `gorm:"type:int(11);not null;default:0"`
}

func (m *ServerGroup) TableName() string {
    return "syd_server_group"
}

func (m *ServerGroup) Create() bool {
    m.Ctime = int(time.Now().Unix())
    return Create(m)
}

func (m *ServerGroup) Update() bool {
    return UpdateByPk(m)
}

func (m *ServerGroup) List(query QueryParam) ([]ServerGroup, bool) {
    var data []ServerGroup
    ok := GetMulti(&data, query)
    return data, ok
}

func (m *ServerGroup) Count(query QueryParam) (int, bool) {
    var count int
    ok := Count(m, &count, query)
    return count, ok
}

func (m *ServerGroup) Delete() bool {
    return DeleteByPk(m)
}

func (m *ServerGroup) Get(id int) bool {
    return GetByPk(m, id)
}
