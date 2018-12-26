// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "time"

    "github.com/tinystack/syncd/model"
)

type Server struct {
    ID      int         `gorm:"primary_key"`
    GroupId int         `gorm:"type:int(11);not null;default:0"`
    Name    string      `gorm:"type:varchar(100);not null;default:''"`
    Ip      string      `gorm:"type:varchar(15);not null;default:''"`
    SshPort int         `gorm:"type:int(11);not null;default:22"`
    Utime   int         `gorm:"type:int(11);not null;default:0"`
}

const (
    TableName = "server"
)

func Create(data *Server) bool {
    data.Utime = int(time.Now().Unix())
    return model.Create(TableName, data)
}

func Update(id int, data Server) bool {
    data.Utime = int(time.Now().Unix())
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

func List(query model.QueryParam) ([]Server, bool) {
    var data []Server
    ok := model.GetMulti(TableName, &data, query)
    return data, ok
}

func Total(query model.QueryParam) (int, bool) {
    var count int
    ok := model.Count(TableName, &count, query)
    return count, ok
}

func Get(id int) (Server, bool) {
    var data Server
    ok := model.GetByPk(TableName, &data, id)
    return data, ok
}

func Delete(id int) bool {
    ok := model.DeleteByPk(TableName, Server{ID: id})
    return ok
}
