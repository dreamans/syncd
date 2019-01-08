// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "strings"

    "github.com/tinystack/goutil/gostring"
    "github.com/tinystack/goweb"
    "github.com/tinystack/govalidate"
    "github.com/tinystack/syncd"
    projectService "github.com/tinystack/syncd/service/project"
    serverService "github.com/tinystack/syncd/service/server"
)

type ProjectParamValid struct {
    Name            string      `valid:"required" errmsg:"required=project name cannot be empty"`
    Description     string      `valid:"require" errmsg:"required=project description cannot be empty"`
    SpaceId         int         `valid:"int_min=1" errmsg:"int_min=space_id cannot be empty"`
    RepoMode        int         `valid:"int_min=1" errmsg:"int_min=repo_mode cannot be empty"`
    RepoUrl         string      `valid:"require" errmsg:"required=repo remote addr cannot be empty"`
    DeployServer    []string    `valid:"require" errmsg:"required=deploy server cannot be empty"`
    DeployUser      string      `valid:"require" errmsg:"required=deploy user cannot be epmty"`
    DeployPath      string      `valid:"require" errmsg:"required=deploy path cannot be epmty"`
    DeployTimeout   int         `valid:"int_min=1" errmsg:"required=int_min cannot be empty"`
}

func ProjectNew(c *goweb.Context) error {
    return projectUpdate(c, 0)
}

func ProjectEdit(c *goweb.Context) error {
    id := c.PostFormInt("id")
    if id == 0 {
        return syncd.RenderParamError("id can not empty")
    }
    return projectUpdate(c, id)
}

func projectUpdate(c *goweb.Context, id int) error {
    params := ProjectParamValid{
        Name: c.PostForm("name"),
        Description: c.PostForm("description"),
        SpaceId: c.PostFormInt("space_id"),
        RepoMode: c.PostFormInt("repo_mode"),
        RepoUrl: c.PostForm("repo_url"),
        DeployServer: c.PostFormArray("deploy_server"),
        DeployUser: c.PostForm("deploy_user"),
        DeployPath: c.PostForm("deploy_path"),
        DeployTimeout: c.PostFormInt("deploy_timeout"),
    }
    if valid := govalidate.NewValidate(&params); !valid.Pass() {
        return syncd.RenderParamError(valid.LastFailed().Msg)
    }
    repoBranch := c.PostForm("repo_branch")
    if params.RepoMode == 1 && repoBranch == "" {
        return syncd.RenderParamError("repo_branch can not be empty")
    }
    var (
        needAudit int
        exists bool
        err error
    )
    projExists := &projectService.Project{
        ID: id,
        SpaceId: params.SpaceId,
        Name: params.Name,
    }
    exists, err = projExists.CheckProjectExists()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if exists {
        return syncd.RenderAppError("project update failed, project name have exists")
    }
    deployServer := gostring.StrSlice2IntSlice(params.DeployServer)
    if c.PostFormInt("need_audit") != 0 {
        needAudit = 1
    }

    excludeFiles := c.PostForm("exclude_files")
    if excludeFiles != "" {
        exFileList := strings.Split(excludeFiles, "\n")
        exFileList = gostring.StrFilterSliceEmpty(exFileList)
        if len(exFileList) > 0 {
            excludeFiles = gostring.JoinSepStrings("\n", exFileList...)
        }
    }

    preDeployCmd, postDeployCmd := c.PostForm("pre_deploy_cmd"), c.PostForm("post_deploy_cmd")
    preDeployCmd = gostring.JoinSepStrings("\n", gostring.StrFilterSliceEmpty(strings.Split(preDeployCmd, "\n"))...)
    postDeployCmd = gostring.JoinSepStrings("\n", gostring.StrFilterSliceEmpty(strings.Split(postDeployCmd, "\n"))...)

    auditEmail, deployEmail := c.PostForm("audit_notice_email"), c.PostForm("deploy_notice_email")
    auditEmail = gostring.JoinSepStrings(",", gostring.StrFilterSliceEmpty(strings.Split(auditEmail, ","))...)
    deployEmail = gostring.JoinSepStrings(",", gostring.StrFilterSliceEmpty(strings.Split(deployEmail, ","))...)

    project := &projectService.Project{
        ID: id,
        Name: params.Name,
        Description: params.Description,
        SpaceId: params.SpaceId,
        RepoUrl: params.RepoUrl,
        RepoMode: params.RepoMode,
        RepoBranch: repoBranch,
        ExcludeFiles: excludeFiles,
        DeployServer: deployServer,
        DeployUser: params.DeployUser,
        DeployPath: params.DeployPath,
        DeployTimeout: params.DeployTimeout,
        PreDeployCmd: preDeployCmd,
        PostDeployCmd: postDeployCmd,
        NeedAudit: needAudit,
        AuditNoticeEmail: auditEmail,
        DeployNoticeEmail: deployEmail,
    }
    if err = project.CreateOrUpdate(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

func ProjectList(c *goweb.Context) error {
    offset, limit, keyword, spaceId, status := c.QueryInt("offset"), c.GetInt("limit"), c.Query("keyword"), c.QueryInt("space_id"), c.QueryInt("status")
    project := &projectService.Project{
        SpaceId: spaceId,
        Status: status,
    }
    list, total, err := project.List(keyword, offset, limit)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "list": list,
        "total": total,
    })
}

func ProjectDetail(c *goweb.Context) error {
    project := &projectService.Project{
        ID: c.QueryInt("id"),
    }
    if err := project.Detail(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    groupList, err := serverService.GroupListByIds(project.DeployServer)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    var deploySers []map[string]interface{}
    for _, l := range groupList {
        deploySers = append(deploySers, map[string]interface{}{
            "id": l.ID,
            "name": l.Name,
        })
    }
    project.DeployServers = deploySers
    return syncd.RenderJson(c, project)
}

func ProjectDelete(c *goweb.Context) error {
    project := &projectService.Project{
        ID: c.PostFormInt("id"),
    }
    if err := project.Detail(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if project.Status != 0 {
        return syncd.RenderAppError("project delete falied, project status must be unavailable")
    }
    if err := project.Delete(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

func ProjectExists(c *goweb.Context) error {
    id, spaceId, keyword := c.QueryInt("id"), c.QueryInt("space_id"), c.Query("keyword")
    if spaceId == 0 || keyword == "" {
        return syncd.RenderParamError("params error")
    }
    project := &projectService.Project{
        ID: id,
        SpaceId: spaceId,
        Name: keyword,
    }
    exists, err := project.CheckProjectExists()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "exists": exists,
    })
}

func ProjectChangeStatus(c *goweb.Context) error {
    id, status := c.PostFormInt("id"), c.PostFormInt("status")
    project := &projectService.Project{
        ID: id,
        Status: status,
    }
    if err := project.ChangeStatus(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}
