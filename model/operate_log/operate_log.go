// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package operate_log

import (
    "time"

    "github.com/tinystack/syncd/model"
)

type OperateLog struct {
    ID              int     `gorm:"primary_key"`
    DataId          int     `gorm:"type:int(11);not null;default:0"`
    OpType          string  `gorm:"type:varchar(10);not null;default:''"`
    OpName          string  `gorm:"type:varchar(100);not null;default:''"`
    OpContent       string  `gorm:"type:varchar(1000);not null;default:''"`
    UserId          int     `gorm:"type:int(11);not null;default:0"`
    UserName        string  `gorm:"type:varchar(100);unique;not null;default:''"`
    Ctime           int     `gorm:"type:int(11);not null;default:0"`
}

const (
    TableName = "operate_log"
)

func Create(data *OperateLog) bool {
    data.Ctime = int(time.Now().Unix())
    return model.Create(TableName, data)
}

func List(query model.QueryParam) ([]OperateLog, bool) {
    var data []OperateLog
    ok := model.GetMulti(TableName, &data, query)
    return data, ok
}

