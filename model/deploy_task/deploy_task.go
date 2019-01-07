// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy_task

import (
    "github.com/tinystack/syncd/model"
)

type DeployTask struct {
    ID              int     `gorm:"primary_key"`
    ApplyId         int     `gorm:"type:int(11);not null;default:0"`
    Level           int     `gorm:"type:int(11);not null;default:0"`
    Cmd             string  `gorm:"type:text;not null"`
    Status          int     `gorm:"type:int(11);not null;default:1"`
    Output          string  `gorm:"type:mediumtext;not null"`
    Name            string  `gorm:"type:varchar(100);not null;default:''"`
}

const (
    TableName = "deploy_task"
)

func Create(data *DeployTask) bool {
    return model.Create(TableName, data)
}

func UpdateByPk(id int, data map[string]interface{}) bool {
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

func List(query model.QueryParam) ([]DeployTask, bool) {
    var data []DeployTask
    ok := model.GetMulti(TableName, &data, query)
    return data, ok
}

func DeleteByApplyId(id int) bool {
    ok := model.Delete(TableName, DeployTask{}, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "apply_id",
                Prepare: id,
            },
        },
    })
    return ok
}

