// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

import(
    "time"
)

type User struct {
    ID              int         `gorm:"primary_key"`
    RoleId          int         `gorm:"type:int(11);not null;default:0"`
    Username	    string      `gorm:"type:varchar(20);not null;default:''"`
    Password        string      `gorm:"type:char(32);not null;default:''"`
    Salt            string      `gorm:"type:char(10);not null;default:''"`
    Truename        string      `gorm:"type:varchar(20);not null;default:''"`
    Mobile          string      `gorm:"type:varchar(20);not null;default:''"`
    Email           string      `gorm:"type:varchar(500);not null;default:''"`
    Status          int         `gorm:"type:int(11);not null;default:0"`
    LastLoginTime   int         `gorm:"type:int(11);not null;default:0"`
    LastLoginIp     string      `gorm:"type:varchar(50);not null;default:''"`
    Ctime           int         `gorm:"type:int(11);not null;default:0"`
}

func (m *User) TableName() string {
    return "syd_user"
}

func (m *User) Create() bool {
    m.Ctime = int(time.Now().Unix())
    return Create(m)
}

func (m *User) UpdateByFields(data map[string]interface{}, query QueryParam) bool {
    return Update(m, data, query)
}

func (m *User) List(query QueryParam) ([]User, bool) {
    var data []User
    ok := GetMulti(&data, query)
    return data, ok
}

func (m *User) Count(query QueryParam) (int, bool) {
    var count int
    ok := Count(m, &count, query)
    return count, ok
}

func (m *User) Delete() bool {
    return DeleteByPk(m)
}

func (m *User) Get(id int) bool {
    return GetByPk(m, id)
}

func (m *User) GetOne(query QueryParam) bool {
	return GetOne(m, query)
}
