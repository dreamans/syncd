// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
    "time"
    "strings"

    "github.com/tinystack/syncd/model"
)

func Create(data *Server) bool {
    data.Utime = int(time.Now().Unix())
    return model.Create(TableName, data)
}

func Update(id int, data Server) bool {
    data.Utime = int(time.Now().Unix())
    ok := model.Update(TableName, data, model.QueryParam{
        Plain: "id = ?",
        Prepare: []interface{}{id},
    })
    return ok
}

func List(fields string, offset, limit int) ([]Server, bool) {
    var data []Server
    ok := model.GetMulti(TableName, &data, model.QueryParam{
        Offset: offset,
        Limit: limit,
        Order: "id desc",
        Fields: fields,
    })
    return data, ok
}

func Multi(fields string, groupId int) ([]Server, bool) {
    q := model.QueryParam{
        Fields: fields,
    }
    var (
        plain   []string
        prepare []interface{}
    )
    if groupId > 0 {
        plain = append(plain, "group_id = ?")
        prepare = append(prepare, groupId)
    }
    if len(plain) > 0 {
        q.Plain = strings.Join(plain, " AND ")
        q.Prepare = prepare
    }
    var data []Server
    ok := model.GetMulti(TableName, &data, q)
    return data, ok
}

func Total() (int, bool) {
    var count int
    ok := model.Count(TableName, &count, model.QueryParam{})
    return count, ok
}

func Get(id int) (Server, bool){
    var data Server
    ok := model.GetOne(TableName, &data, model.QueryParam{
        Plain: "id = ?",
        Prepare: []interface{}{id},
    })
    return data, ok
}

func Delete(id int) bool {
    ok := model.Delete(TableName, Server{}, model.QueryParam{
        Plain: "id = ?",
        Prepare: []interface{}{id},
    })
    return ok
}
