// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

import(
    "time"
)

type DeployBuild struct {
    ID                  int     `gorm:"primary_key"`
    ApplyId             int     `gorm:"type:int(11);not null;default:0"`
    StartTime           int     `gorm:"type:int(11);not null;default:0"`
    FinishTime          int     `gorm:"type:int(11);not null;default:0"`
    Status              int     `gorm:"type:int(11);not null;default:0"`
    Tar                 string  `gorm:"type:varchar(2000);not null;default:''"`
    Output              string  `gorm:"type:text;not null"`
    Errmsg              string  `gorm:"type:text;not null"`
    Ctime               int     `gorm:"type:int(11);not null;default:0"`
}

func (m *DeployBuild) TableName() string {
    return "syd_deploy_build"
}

func (m *DeployBuild) GetByApplyId(id int) bool {
    return GetOne(m, QueryParam{
        Where: []WhereParam{
            WhereParam{
                Field: "apply_id",
                Prepare: id,
            },
        },
    })
}

func (m *DeployBuild) UpdateByFields(data map[string]interface{}, query QueryParam) bool {
    return Update(m, data, query)
}

func (m *DeployBuild) Create() bool {
    m.Ctime = int(time.Now().Unix())
    return Create(m)
}

func (m *DeployBuild) Delete() bool {
    return DeleteByPk(m)
}
