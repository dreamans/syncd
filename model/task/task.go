// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
    "time"

    "github.com/tinystack/syncd/model"
)

type Task struct {
    ID              int     `gorm:"primary_key"`
    Name            string  `gorm:"type:varchar(100);not null;default:''"`
    Key             string  `gorm:"type:char(32);not null;default:''"`
    Cmd             string  `gorm:"type:varchar(2000);not null;default:''"`
    Status          int     `gorm:"type:int(11);not null;default:1"`
    Ctime           int     `gorm:"type:int(11);not null;default:0"`
    Ftime           int     `gorm:"type:int(11);not null;default:0"`
}

const (
    TableName = "task"
)

func Create(data *Task) bool {
    data.Ctime = int(time.Now().Unix())
    return model.Create(TableName, data)
}

func UpdateByKey(key string, data map[string]interface{}) bool {
    ok := model.Update(TableName, data, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "`key`",
                Prepare: key,
            },
        },
    })
    return ok
}

