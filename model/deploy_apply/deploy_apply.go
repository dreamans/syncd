// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy_apply

import (
    "time"

    "github.com/tinystack/syncd/model"
)

type DeployApply struct {
    ID              int     `gorm:"primary_key"`
    ProjectId       int     `gorm:"type:int(11);not null;default:0"`
    Name            string  `gorm:"type:varchar(100);not null;default:''"`
    Description     string  `gorm:"type:varchar(2000);not null;default:''"`
    SpaceId         int     `gorm:"type:int(11);not null;default:0"`
    RepoData        string  `gorm:"type:varchar(10000);not null;default:''"`
    Status          int     `gorm:"type:int(11);not null;default:1"`
    ErrorLog        string  `gorm:"type:mediumtext;not null"`
    UserId          int     `gorm:"type:int(11);not null;default:0"`
    Ctime           int     `gorm:"type:int(11);not null;default:0"`
}

const (
    TableName = "deploy_apply"
)

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

func Create(data *DeployApply) bool {
    data.Ctime = int(time.Now().Unix())
    return model.Create(TableName, data)
}

func List(query model.QueryParam) ([]DeployApply, bool) {
    var data []DeployApply
    ok := model.GetMulti(TableName, &data, query)
    return data, ok
}

func Total(query model.QueryParam) (int, bool) {
    var count int
    ok := model.Count(TableName, &count, query)
    return count, ok
}

func Get(id int) (DeployApply, bool) {
    var data DeployApply
    ok := model.GetByPk(TableName, &data, id)
    return data, ok
}

