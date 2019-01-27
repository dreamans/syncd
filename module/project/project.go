// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "errors"
    //"fmt"

    "github.com/dreamans/syncd/model"
    "github.com/dreamans/syncd/util/gostring"
)

type Project struct {
    ID                  int     `json:"id"`
    Name                string  `json:"name"`
    Description         string  `json:"description"`
    NeedAudit           int     `json:"need_audit"`
    RepoUrl             string  `json:"repo_url"`
    RepoBranch          string  `json:"repo_branch"`
    PreReleaseCluster   int     `json:"pre_release_cluster"`
    OnlineCluster       []int   `json:"online_cluster"`
    DeployUser          string  `json:"deploy_user"`
    DeployPath          string  `json:"deploy_path"`
    PreDeployCmd        string  `json:"pre_deploy_cmd"`
    AfterDeployCmd      string  `json:"after_deploy_cmd"`
    DeployTimeout       int     `json:"deploy_timeout"`
    Ctime               int     `json:"ctime"`
}

func (p *Project) CreateOrUpdate() error {
    project := &model.Project{
        ID: p.ID,
        Name: p.Name,
        Description: p.Description,
        NeedAudit: p.NeedAudit,
        RepoUrl: p.RepoUrl,
        RepoBranch: p.RepoBranch,
        PreReleaseCluster: p.PreReleaseCluster,
        OnlineCluster: gostring.JoinIntSlice2String(p.OnlineCluster, ","),
        DeployUser: p.DeployUser,
        DeployPath: p.DeployPath,
        PreDeployCmd: p.PreDeployCmd,
        AfterDeployCmd: p.AfterDeployCmd,
        DeployTimeout: p.DeployTimeout,
    }
    if project.ID > 0 {
        if ok := project.Update(); !ok {
            return errors.New("create project failed")
        }
    } else {
        if ok := project.Create(); !ok {
            return errors.New("update project failed")
        }
    }
    return nil
}
