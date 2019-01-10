// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package token

import (
    "time"

    "github.com/dreamans/syncd/model"
)

type UserToken struct {
    ID              int     `gorm:"primary_key"`
    UserId          int     `gorm:"type:int(11);unique;not null;default:0"`
    Token           string  `gorm:"type:varchar(40);not null;default:''"`
    ExpireTime      int     `gorm:"type:int(11);not null;default:0"`
    Ctime           int     `gorm:"type:int(11);not null;default:0"`
}

const (
    TableName = "user_token"
)

func Create(data *UserToken) bool {
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

func Get(id int) (UserToken, bool) {
    var data UserToken
    ok := model.GetByPk(TableName, &data, id)
    return data, ok
}

func GetOne(query model.QueryParam) (UserToken, bool) {
    var data UserToken
    ok := model.GetOne(TableName, &data, query)
    return data, ok
}

func DeleteByUserId(id int) bool {
    ok := model.Delete(TableName, UserToken{}, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "user_id",
                Prepare: id,
            },
        },
    })
    return ok
}
