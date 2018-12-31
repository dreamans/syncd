// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task_log

import (
    "time"

    "github.com/tinystack/syncd/model"
)

type TaskLog struct {
    ID              int     `gorm:"primary_key"`
    Key             string  `gorm:"type:char(32);not null;default:''"`
    Content         string  `gorm:"type:mediumtext;not null"`
    Ctime           int     `gorm:"type:int(11);not null;default:0"`
}

const (
    TableName = "task_log"
)

func Create(data *TaskLog) bool {
    data.Ctime = int(time.Now().Unix())
    return model.Create(TableName, data)
}

