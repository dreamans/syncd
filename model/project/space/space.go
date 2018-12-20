// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package space

import (
    "time"

    "github.com/tinystack/syncd/model"
)

func Create(data *ProjectSpace) bool {
    data.Ctime = int(time.Now().Unix())
    return model.Create(TableName, data)
}

func Update(id int, data ProjectSpace) bool {
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

func List(query model.QueryParam) ([]ProjectSpace, bool) {
    var data []ProjectSpace
    ok := model.GetMulti(TableName, &data, query)
    return data, ok
}

func Total(query model.QueryParam) (int, bool) {
    var count int
    ok := model.Count(TableName, &count, query)
    return count, ok
}

func Get(id int) (ProjectSpace, bool){
    var data ProjectSpace
    ok := model.GetByPk(TableName, &data, id)
    return data, ok
}

func GetOne(query model.QueryParam) (ProjectSpace, bool) {
    var data ProjectSpace
    ok := model.GetOne(TableName, &data, query)
    return data, ok
}

func Delete(id int) bool {
    ok := model.DeleteByPk(TableName, ProjectSpace{ID: id})
    return ok
}
