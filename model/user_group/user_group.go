// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package group

import (
    "time"

    "github.com/dreamans/syncd/model"
)

type UserGroup struct {
    ID      int         `gorm:"primary_key"`
    Name    string      `gorm:"type:varchar(100);not null;default:''"`
    Priv    string      `gorm:"type:varchar(10000);not null;default:''"`
    Ctime   int         `gorm:"type:int(11);not null;default:0"`
}

const (
    TableName = "user_group"
)

func Create(data *UserGroup) bool {
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

func List(query model.QueryParam) ([]UserGroup, bool) {
    var data []UserGroup
    ok := model.GetMulti(TableName, &data, query)
    return data, ok
}

func Total(query model.QueryParam) (int, bool) {
    var count int
    ok := model.Count(TableName, &count, query)
    return count, ok
}

func Get(id int) (UserGroup, bool){
    var data UserGroup
    ok := model.GetByPk(TableName, &data, id)
    return data, ok
}

func Delete(id int) bool {
    ok := model.DeleteByPk(TableName, UserGroup{ID: id})
    return ok
}

func GetOne(query model.QueryParam) (UserGroup, bool) {
    var data UserGroup
    ok := model.GetOne(TableName, &data, query)
    return data, ok
}