// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "time"

    "github.com/tinystack/syncd/model"
)

type User struct {
    ID              int     `gorm:"primary_key"`
    GroupId         int     `gorm:"type:int(11);not null;default:0"`
    Name            string  `gorm:"type:varchar(100);unique;not null;default:''"`
    Password        string  `gorm:"type:char(32);not null;default:''"`
    Email           string  `gorm:"type:varchar(100);unique;not null;default:''"`
    TrueName        string  `gorm:"type:varchar(20);not null;default:''"`
    Mobile          string  `gorm:"type:varchar(20);not null;default:''"`
    Salt            string  `gorm:"type:varchar(10);not null;default:''"`
    LockStatus      int     `gorm:"type:int(11);not null;default:0"`
    LastLoginIp     string  `gorm:"type:varchar(45);not null;default:''"`
    LastLoginTime   int     `gorm:"type:int(11);not null;default:0"`
    Ctime           int     `gorm:"type:int(11);not null;default:0"`
}

const (
    TableName = "user"
)

func Create(data *User) bool {
    data.Ctime = int(time.Now().Unix())
    return model.Create(TableName, data)
}

func Update(id int, data map[string]interface{}) bool {
    ok := model.Update(TableName, data, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Prepare: id,
            },
        },
    })
    return ok
}

func List(query model.QueryParam) ([]User, bool) {
    var data []User
    ok := model.GetMulti(TableName, &data, query)
    return data, ok
}

func Total(query model.QueryParam) (int, bool) {
    var count int
    ok := model.Count(TableName, &count, query)
    return count, ok
}

func Get(id int) (User, bool) {
    var data User
    ok := model.GetByPk(TableName, &data, id)
    return data, ok
}

func GetOne(query model.QueryParam) (User, bool) {
    var data User
    ok := model.GetOne(TableName, &data, query)
    return data, ok
}

func Delete(id int) bool {
    ok := model.DeleteByPk(TableName, User{ID: id})
    return ok
}

