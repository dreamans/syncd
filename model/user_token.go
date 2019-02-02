// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

import(
    "time"
)

type UserToken struct {
    ID              int         `gorm:"primary_key"`
    UserId          int         `gorm:"type:int(11);not null;default:0"`
    Token           string      `gorm:"type:varchar(100);not null;default:''"`
    Expire          int         `gorm:"type:int(11);not null;default:0"`
    Ctime           int         `gorm:"type:int(11);not null;default:0"`
}

func (m *UserToken) TableName() string {
    return "syd_user_token"
}

func (m *UserToken) Create() bool {
    m.Ctime = int(time.Now().Unix())
    return Create(m)
}

func (m *UserToken) UpdateByFields(data map[string]interface{}, query QueryParam) bool {
    return Update(m, data, query)
}

func (m *UserToken) Update() bool {
    return UpdateByPk(m)
}

func (m *UserToken) List(query QueryParam) ([]UserToken, bool) {
    var data []UserToken
    ok := GetMulti(&data, query)
    return data, ok
}

func (m *UserToken) Count(query QueryParam) (int, bool) {
    var count int
    ok := Count(m, &count, query)
    return count, ok
}

func (m *UserToken) Delete() bool {
    return DeleteByPk(m)
}

func (m *UserToken) DeleteByFields(query QueryParam) bool {
    return Delete(m, query)
}

func (m *UserToken) Get(id int) bool {
    return GetByPk(m, id)
}

func (m *UserToken) GetOne(query QueryParam) bool {
    return GetOne(m, query)
}
