
// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server_group

import (
    "time"

    "github.com/tinystack/syncd/model"
)

type ServerGroup struct {
    ID      int         `gorm:"primary_key"`
    Name    string      `gorm:"type:varchar(100);not null;default:''"`
    Ctime   int         `gorm:"type:int(11);not null;default:0"`
}

const (
    TableName = "server_group"
)

func Create(data *ServerGroup) bool {
    data.Ctime = int(time.Now().Unix())
    return model.Create(TableName, data)
}

func Update(id int, data ServerGroup) bool {
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

func List(query model.QueryParam) ([]ServerGroup, bool) {
    var data []ServerGroup
    ok := model.GetMulti(TableName, &data, query)
    return data, ok
}

func Total(query model.QueryParam) (int, bool) {
    var count int
    ok := model.Count(TableName, &count, query)
    return count, ok
}

func Get(id int) (ServerGroup, bool){
    var data ServerGroup
    ok := model.GetByPk(TableName, &data, id)
    return data, ok
}

func Delete(id int) bool {
    ok := model.DeleteByPk(TableName, ServerGroup{ID: id})
    return ok
}

