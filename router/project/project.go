// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/router/common"
    "github.com/dreamans/syncd/module/project"
    "github.com/dreamans/syncd/util/gostring"
    "github.com/dreamans/syncd/util/goslice"
)

type ProjectFormBind struct {
    ID                  int     `form:"id"`
    SpaceId             int     `form:"space_id"`
    Name                string  `form:"name" binding:"required"`
    Description         string  `form:"description"`
    NeedAudit           int     `form:"need_audit"`
    RepoUrl             string  `form:"repo_url" binding:"required"`
    RepoBranch          string  `form:"repo_branch"`
    DeployMode          int     `form:"deploy_mode" binding:"required"`
    OnlineCluster       []int   `form:"online_cluster" binding:"required"`
    DeployUser          string  `form:"deploy_user" binding:"required"`
    DeployPath          string  `form:"deploy_path" binding:"required"`
    PreDeployCmd        string  `form:"pre_deploy_cmd"`
    AfterDeployCmd      string  `form:"after_deploy_cmd"`
    AuditNotice         string  `form:"audit_notice"`
    DeployNotice        string  `form:"deploy_notice"`
}

type ProjectBuildScriptBind struct {
    ID                  int     `form:"id" binding:"required"`
    BuildScript         string  `form:"build_script" binding:"required"`
}

type ProjectHookScriptBind struct {
    ID                      int     `form:"id" binding:"required"`
    BuildHookScript         string  `form:"build_hook_script"`
    DeployHookScript         string  `form:"deploy_hook_script"`
}

type QueryBind struct {
    SpaceId     int     `form:"space_id"`
    Keyword     string  `form:"keyword"`
    Offset      int     `form:"offset"`
    Limit       int     `form:"limit" binding:"required,gte=1,lte=999"`
}

func ProjectBuildScript(c *gin.Context) {
    var form ProjectBuildScriptBind
    if err := c.ShouldBind(&form); err != nil {
        render.ParamError(c, err.Error())
        return
    }

    p := &project.Project{
        ID: form.ID,
    }
    if err := p.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    if !common.InSpaceCheck(c, p.SpaceId) {
        return
    }

    proj := &project.Project{
        ID: form.ID,
        BuildScript: form.BuildScript,
    }
    if err := proj.UpdateBuildScript(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.Success(c)
}

func ProjectHookScript(c *gin.Context) {
    var form ProjectHookScriptBind
    if err := c.ShouldBind(&form); err != nil {
        render.ParamError(c, err.Error())
        return
    }

    p := &project.Project{
        ID: form.ID,
    }
    if err := p.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    if !common.InSpaceCheck(c, p.SpaceId) {
        return
    }

    proj := &project.Project{
        ID: form.ID,
        BuildHookScript: form.BuildHookScript,
        DeployHookScript: form.DeployHookScript,
    }
    if err := proj.UpdateHookScript(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.Success(c)
}

func ProjectDelete(c *gin.Context) {
    id := gostring.Str2Int(c.PostForm("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    proj := &project.Project{
        ID: id,
    }

    if err := proj.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    if !common.InSpaceCheck(c, proj.SpaceId) {
        return
    }

    if err := proj.Delete(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.Success(c)
}

func ProjectDetail(c *gin.Context) {
    id := gostring.Str2Int(c.Query("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    proj := &project.Project{
        ID: id,
    }
    if err := proj.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    if !common.InSpaceCheck(c, proj.SpaceId) {
        return
    }

    render.JSON(c, proj)
}

func ProjectSwitchStatus(c *gin.Context) {
    id, status := gostring.Str2Int(c.PostForm("id")), gostring.Str2Int(c.PostForm("status"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }

    p := &project.Project{
        ID: id,
    }
    if err := p.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    if !common.InSpaceCheck(c, p.SpaceId) {
        return
    }

    if status !=0 {
        status = 1
    }
    proj := &project.Project{
        ID: id,
        Status: status,
    }
    if err := proj.UpdateStatus(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.Success(c)
}

func ProjectList(c *gin.Context) {
    var query QueryBind
    if err := c.ShouldBind(&query); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    if query.SpaceId == 0 {
        render.ParamError(c, "space_id cannot be empty")
        return
    }

    if !common.InSpaceCheck(c, query.SpaceId) {
        return
    }

    proj := &project.Project{}
    list, err := proj.List(query.Keyword, query.SpaceId, query.Offset, query.Limit)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    total, err := proj.Total(query.Keyword, query.SpaceId)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    projList := []map[string]interface{}{}
    for _, l := range list {
        projList = append(projList, map[string]interface{}{
            "id": l.ID,
            "name": l.Name,
            "need_audit": l.NeedAudit,
            "status": l.Status,
        })
    }

    render.JSON(c, gin.H{
        "list": projList,
        "total": total,
    })
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

    if !common.InSpaceCheck(c, projectForm.SpaceId) {
        return
    }

    onlineCluster := goslice.FilterSliceInt(projectForm.OnlineCluster)
    if len(onlineCluster) == 0 {
        render.ParamError(c, "online_cluster cannot be empty")
        return
    }
    repoBranch := projectForm.RepoBranch
    if projectForm.DeployMode == 2 {
        repoBranch = ""
    }
    proj := &project.Project{
        ID: projectForm.ID,
        SpaceId: projectForm.SpaceId,
        Name: projectForm.Name,
        Description: projectForm.Description,
        NeedAudit: projectForm.NeedAudit,
        RepoUrl: projectForm.RepoUrl,
        DeployMode: projectForm.DeployMode,
        RepoBranch: repoBranch,
        OnlineCluster: onlineCluster,
        DeployUser: projectForm.DeployUser,
        DeployPath: projectForm.DeployPath,
        PreDeployCmd: projectForm.PreDeployCmd,
        AfterDeployCmd: projectForm.AfterDeployCmd,
        AuditNotice: projectForm.AuditNotice,
        DeployNotice: projectForm.DeployNotice,
    }
    if err := proj.CreateOrUpdate(); err != nil {
        render.AppError(c, err.Error())
        return 
    }
    render.Success(c)
}