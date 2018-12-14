// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "time"

    "github.com/tinystack/syncd/model"
)

func Create(data *Project) bool {
    data.Utime = int(time.Now().Unix())
    return model.Create(TableName, data)
}

func List(fields string, offset, limit int) ([]ProjectList, bool) {
    var p []ProjectList
    ok := model.GetMulti(TableName, &p, model.QueryParam{
        Offset: offset,
        Limit: limit,
        Order: "id desc",
        Fields: fields,
    })
    return p, ok
}

func Total() (int, bool) {
    var count int
    ok := model.Count(TableName, &count, model.QueryParam{})
    return count, ok
}

func Get(id int) (Project, bool){
    var p Project
    ok := model.GetOne(TableName, &p, model.QueryParam{
        Plain: "id = ?",
        Prepare: []interface{}{id},
    })
    return p, ok
}

func Update(id int, data Project) bool {
    data.Utime = int(time.Now().Unix())
    ok := model.Update(TableName, data, model.QueryParam{
        Plain: "id = ?",
        Prepare: []interface{}{id},
    })
    return ok
}
