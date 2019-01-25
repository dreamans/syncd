// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

import(
    "time"
)

type ProjectSpace struct {
    ID              int     `gorm:"primary_key"`
    Name            string  `gorm:"type:varchar(100);not null;default:''"`
    Description     string  `gorm:"type:varchar(2000);not null;default:''"`
    Ctime           int     `gorm:"type:int(11);not null;default:0"`
}

func (m *ProjectSpace) TableName() string {
    return "syd_project_space"
}

func (m *ProjectSpace) Create() bool {
    m.Ctime = int(time.Now().Unix())
    return Create(m)
}

func (m *ProjectSpace) Update() bool {
    return UpdateByPk(m)
}

func (m *ProjectSpace) List(query QueryParam) ([]ProjectSpace, bool) {
    var data []ProjectSpace
    ok := GetMulti(&data, query)
    return data, ok
}

func (m *ProjectSpace) Count(query QueryParam) (int, bool) {
    var count int
    ok := Count(m, &count, query)
    return count, ok
}

func (m *ProjectSpace) Delete() bool {
    return DeleteByPk(m)
}

func (m *ProjectSpace) Get(id int) bool {
    return GetByPk(m, id)
}
