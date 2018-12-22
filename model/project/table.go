// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

type Project struct {
    ID              int     `gorm:"primary_key"`
    Name            string  `gorm:"type:varchar(100);not null;default:''"`
    Description     string  `gorm:"type:varchar(100);not null;default:''"`
    SpaceId         int     `gorm:"type:int(11);not null;default:0"`
    Space           string  `gorm:"type:varchar(100);not null;default:''"`
    Repo            string  `gorm:"type:varchar(20);not null;default:''"`
    RepoUrl         string  `gorm:"type:varchar(200);not null;default:''"`
    DeployServer    string  `gorm:"type:varchar(2000);not null;default:''"`
    DeployUser      string  `gorm:"type:varchar(20);not null;default:''"`
    DeployPath      string  `gorm:"type:varchar(100);not null;default:''"`
    DeployHistory   int     `gorm:"type:int(11);not null;default:0"`
    PreDeployCmd    string  `gorm:"type:varchar(2000);not null;default:''"`
    PostDeployCmd   string  `gorm:"type:varchar(2000);not null;default:''"`
    NeedAudit       int     `gorm:"type:int(11);not null;default:0"`
    Status          int     `gorm:"type:int(11);not null;default:0"`
    RepoUser        string  `gorm:"type:varchar(100);not null;default:''"`
    RepoPass        string  `gorm:"type:varchar(100);not null;default:''"`
    RepoMode        int     `gorm:"type:int(11);not null;default:0"`
    BuildScript     string  `gorm:"type:text;not null"`
    Utime           int     `gorm:"type:int(11);not null;default:0"`
}

const (
    TableName = "project"
)
