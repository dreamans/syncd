// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package group

import (
    "time"

    "github.com/tinystack/syncd/model"
)

func Create(data *Group) bool {
    data.Utime = int(time.Now().Unix())
    return model.Create(TableName, data)
}

func Update(id int, data Group) bool {
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

func List(query model.QueryParam) ([]Group, bool) {
    var data []Group
    ok := model.GetMulti(TableName, &data, query)
    return data, ok
}

func Total(query model.QueryParam) (int, bool) {
    var count int
    ok := model.Count(TableName, &count, query)
    return count, ok
}

func Get(id int) (Group, bool){
    var data Group
    ok := model.GetByPk(TableName, &data, id)
    return data, ok
}

func Delete(id int) bool {
    ok := model.DeleteByPk(TableName, Group{ID: id})
    return ok
}
