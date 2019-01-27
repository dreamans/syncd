// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/module/project"
    "github.com/dreamans/syncd/util/gostring"
    "github.com/dreamans/syncd/util/goslice"
)

type ProjectFormBind struct {
    ID                  int     `form:"id"`
    Name                string  `form:"name" binding:"required"`
    Description         string  `form:"description"`
    NeedAudit           int     `form:"need_audit"`
    RepoUrl             string  `form:"repo_url" binding:"required"`
    RepoBranch          string  `form:"repo_branch"`
    PreReleaseCluster   int     `form:"pre_release_cluster"`
    OnlineCluster       []int   `form:"online_cluster" binding:"required"`
    DeployUser          string  `form:"deploy_user" binding:"required"`
    DeployPath          string  `form:"deploy_path" binding:"required"`
    PreDeployCmd        string  `form:"pre_deploy_cmd"`
    AfterDeployCmd      string  `form:"after_deploy_cmd"`
    DeployTimeout       int     `form:"deploy_timeout" binding:"required"`
}

func ProjectAdd(c *gin.Context) {
    projectCreateOrUpdate(c)
}

func ProjectUpdate(c *gin.Context) {
    id := gostring.Str2Int(c.PostForm("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    projectCreateOrUpdate(c)
}

func projectCreateOrUpdate(c *gin.Context) {
    var projectForm ProjectFormBind
    if err := c.ShouldBind(&projectForm); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    onlineCluster := goslice.FilterSliceInt(projectForm.OnlineCluster)
    if len(onlineCluster) == 0 {
        render.ParamError(c, "online_cluster cannot be empty")
        return
    }
    proj := &project.Project{
        ID: projectForm.ID,
        Name: projectForm.Name,
        Description: projectForm.Description,
        NeedAudit: projectForm.NeedAudit,
        RepoUrl: projectForm.RepoUrl,
        RepoBranch: projectForm.RepoBranch,
        PreReleaseCluster: projectForm.PreReleaseCluster,
        OnlineCluster: onlineCluster,
        DeployUser: projectForm.DeployUser,
        DeployPath: projectForm.DeployPath,
        PreDeployCmd: projectForm.PreDeployCmd,
        AfterDeployCmd: projectForm.AfterDeployCmd,
        DeployTimeout: projectForm.DeployTimeout,
    }
    if err := proj.CreateOrUpdate(); err != nil {
        render.AppError(c, err.Error())
        return 
    }
    render.Success(c)
}