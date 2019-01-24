// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

import(
    "time"
)

type UserRole struct {
    ID          int         `gorm:"primary_key"`
    Name        string      `gorm:"type:varchar(100);not null;default:''"`
    Privilege   string	    `gorm:"type:varchar(2000);not null;default:''"`
    Ctime       int         `gorm:"type:int(11);not null;default:0"`
}

func (m *UserRole) TableName() string {
    return "syd_user_role"
}

func (m *UserRole) Create() bool {
    m.Ctime = int(time.Now().Unix())
    return Create(m)
}

func (m *UserRole) Update() bool {
    return UpdateByPk(m)
}

func (m *UserRole) List(query QueryParam) ([]UserRole, bool) {
    var data []UserRole
    ok := GetMulti(&data, query)
    return data, ok
}

func (m *UserRole) Count(query QueryParam) (int, bool) {
    var count int
    ok := Count(m, &count, query)
    return count, ok
}

func (m *UserRole) Delete() bool {
    return DeleteByPk(m)
}

func (m *UserRole) Get(id int) bool {
    return GetByPk(m, id)
}
