// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

import(
    "time"
)

type DeployApply struct {
    ID                  int     `gorm:"primary_key"`
    SpaceId             int     `gorm:"type:int(11);not null;default:0"`
    ProjectId           int     `gorm:"type:int(11);not null;default:0"`
    Name                string  `gorm:"type:varchar(100);not null;default:''"`
    Description         string  `gorm:"type:varchar(500);not null;default:''"`
    BranchName          string  `gorm:"type:varchar(100);not null;default:''"`
    CommitVersion       string  `gorm:"type:varchar(50);not null;default:''"`
    AuditStatus         int     `gorm:"type:int(11);not null;default:0"`
    AuditRefusalReasion string  `gorm:"type:varchar(500);not null;default:''"`
    Status              int     `gorm:"type:int(11);not null;default:0"`
    UserId              int     `gorm:"type:int(11);not null;default:0"`
    RollbackId          int     `gorm:"type:int(11);not null;default:0"`
    RollbackApplyId     int     `gorm:"type:int(11);not null;default:0"`
    IsRollbackApply     int     `gorm:"type:int(11);not null;default:0"`
    Ctime               int     `gorm:"type:int(11);not null;default:0"`
}

func (m *DeployApply) TableName() string {
    return "syd_deploy_apply"
}

func (m *DeployApply) Create() bool {
    m.Ctime = int(time.Now().Unix())
    return Create(m)
}

func (m *DeployApply) Update() bool {
    return UpdateByPk(m)
}

func (m *DeployApply) UpdateByFields(data map[string]interface{}, query QueryParam) bool {
    return Update(m, data, query)
}

func (m *DeployApply) List(query QueryParam) ([]DeployApply, bool) {
    var data []DeployApply
    ok := GetMulti(&data, query)
    return data, ok
}

func (m *DeployApply) Count(query QueryParam) (int, bool) {
    var count int
    ok := Count(m, &count, query)
    return count, ok
}

func (m *DeployApply) Delete() bool {
    return DeleteByPk(m)
}

func (m *DeployApply) Get(id int) bool {
    return GetByPk(m, id)
}
