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

func List(query model.QueryParam) ([]ProjectList, bool) {
    var p []ProjectList
    ok := model.GetMulti(TableName, &p, query)
    return p, ok
}

func Total(query model.QueryParam) (int, bool) {
    var count int
    ok := model.Count(TableName, &count, query)
    return count, ok
}

func Get(id int) (Project, bool){
    var data Project
    ok := model.GetByPk(TableName, &data, id)
    return data, ok
}

func Update(id int, data Project) bool {
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
