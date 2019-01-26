// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

import(
    "time"
)

type ProjectMember struct {
    ID              int     `gorm:"primary_key"`
    SpaceId         int     `gorm:"type:int(11);not null;default:0"`
    UserId          int     `gorm:"type:int(11);not null;default:0"`
    Ctime           int     `gorm:"type:int(11);not null;default:0"`
}

func (m *ProjectMember) TableName() string {
    return "syd_project_member"
}

func (m *ProjectMember) Create() bool {
    m.Ctime = int(time.Now().Unix())
    return Create(m)
}

func (m *ProjectMember) Update() bool {
    return UpdateByPk(m)
}

func (m *ProjectMember) List(query QueryParam) ([]ProjectMember, bool) {
    var data []ProjectMember
    ok := GetMulti(&data, query)
    return data, ok
}

func (m *ProjectMember) Count(query QueryParam) (int, bool) {
    var count int
    ok := Count(m, &count, query)
    return count, ok
}

func (m *ProjectMember) Delete() bool {
    return DeleteByPk(m)
}

func (m *ProjectMember) Get(id int) bool {
    return GetByPk(m, id)
}
