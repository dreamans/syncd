// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

type Project struct {
    ID              int     `gorm:"primary_key" json:"id"`
    Name            string  `gorm:"type:varchar(100);not null;default:''" json:"name"`
    Description     string  `gorm:"type:varchar(100);not null;default:''" json:"description"`
    Space           string  `gorm:"type:varchar(100);not null;default:''" json:"space"`
    Repo            string  `gorm:"type:varchar(20);not null;default:''" json:"repo"`
    RepoUrl         string  `gorm:"type:varchar(200);not null;default:''" json:"repo_url"`
    DeployServer    string  `gorm:"type:varchar(2000);not null;default:''" json:"deploy_server"`
    DeployUser      string  `gorm:"type:varchar(20);not null;default:''" json:"deploy_user"`
    DeployPath      string  `gorm:"type:varchar(100);not null;default:''" json:"deploy_path"`
    DeployHistory   int     `gorm:"type:int(11);not null;default:0" json:"deploy_history"`
    PreDeployCmd    string  `gorm:"type:varchar(2000);not null;default:''" json:"pre_deploy_cmd"`
    PostDeployCmd   string  `gorm:"type:varchar(2000);not null;default:''" json:"post_deploy_cmd"`
    NeedAudit       int     `gorm:"type:int(11);not null;default:0" json:"need_audit"`
    Status          int     `gorm:"type:int(11);not null;default:0" json:"status"`
    RepoUser        string  `gorm:"type:varchar(100);not null;default:''" json:"repo_user"`
    RepoPass        string  `gorm:"type:varchar(100);not null;default:''" json:"repo_pass"`
    RepoMode        int     `gorm:"type:int(11);not null;default:0" json:"repo_mode"`
    BuildScript     string  `gorm:"type:text;not null" json:"build_script"`
    Utime           int     `gorm:"type:int(11);not null;default:0" json:"utime"`
}

type ProjectList struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    RepoMode    int     `json:"repo_mode"`
    NeedAudit   int     `json:"need_audit"`
    Status      int     `json:"status"`
}

const (
    TableName = "syd_project"
)
