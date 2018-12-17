// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "github.com/tinystack/goutil"
    "github.com/tinystack/goweb"
    "github.com/tinystack/govalidate"
    "github.com/tinystack/syncd/route"
    "github.com/tinystack/syncd"
    projectService "github.com/tinystack/syncd/service/project"
)

func init() {
    route.Register(route.API_PROJECT_UPDATE, updateProject)
    route.Register(route.API_PROJECT_LIST, listProject)
    route.Register(route.API_PROJECT_DETAIL, detailProject)
    route.Register(route.API_PROJECT_DELETE, deleteProject)
}

type ProjectParamValid struct {
    Name            string      `valid:"required" errmsg:"required=project name cannot be empty"`
    Description     string      `valid:"require" errmsg:"required=project description cannot be empty"`
    Space           string      `valid:"require" errmsg:"required=project space cannot be empty"`
    Repo            string      `valid:"require" errmsg:"required=repo type cannot be empty"`
    RepoMode        int         `valid:"require" errmsg:"required=repo mode cannot be empty"`
    RepoUrl         string      `valid:"require" errmsg:"required=repo remote addr cannot be empty"`
    DeployServer    []string    `valid:"require" errmsg:"required=deploy server cannot be empty"`
    DeployUser      string      `valid:"require" errmsg:"required=deploy user cannot be epmty"`
    DeployPath      string      `valid:"require" errmsg:"required=deploy path cannot be epmty"`
    DeployHistory   int         `valid:"int_min=3" errmsg:"int_min=deploy history at least 3"`
}

func updateProject(c *goweb.Context) error {
    params := ProjectParamValid{
        Name: c.PostForm("name"),
        Description: c.PostForm("description"),
        Space: c.PostForm("space"),
        Repo: c.PostForm("repo"),
        RepoMode: c.PostFormInt("repoMode"),
        RepoUrl: c.PostForm("repoUrl"),
        DeployServer: c.PostFormArray("deployServer"),
        DeployUser: c.PostForm("deployUser"),
        DeployPath: c.PostForm("deployPath"),
        DeployHistory: c.PostFormInt("deployHistory"),
    }
    if valid := govalidate.NewValidate(&params); !valid.Pass() {
        syncd.RenderParamError(c, valid.LastFailed().Msg)
        return nil
    }

    deployServer := goutil.StrSlice2IntSlice(params.DeployServer)
    needAudit := 0
    if c.PostFormInt("needAudit") != 0 {
        needAudit = 1
    }
    status := 0
    if c.PostFormInt("status") != 0 {
        status = 1
    }

    project := &projectService.Project{
        ID: c.PostFormInt("id"),
        Name: params.Name,
        Description: params.Description,
        Space: params.Space,
        Repo: params.Repo,
        RepoUrl: params.RepoUrl,
        RepoMode: params.RepoMode,
        DeployServer: deployServer,
        DeployUser: params.DeployUser,
        DeployPath: params.DeployPath,
        DeployHistory: params.DeployHistory,
        PreDeployCmd: c.PostForm("preDeployCmd"),
        PostDeployCmd: c.PostForm("postDeployCmd"),
        NeedAudit: needAudit,
        Status: status,
        RepoUser: c.PostForm("repoUser"),
        RepoPass: c.PostForm("repoPass"),
        BuildScript: c.PostForm("buildScript"),
    }
    if err := project.CreateOrUpdate(); err != nil {
        syncd.RenderParamError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}

func listProject(c *goweb.Context) error {
    offset, limit := c.QueryInt("offset"), c.QueryInt("limit")
    keyword := c.Query("keyword")

    project := &projectService.Project{}
    list, total, err := project.List(keyword, offset, limit)
    if err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, goweb.JSON{
        "list": list,
        "total": total,
    })
    return nil
}

func detailProject(c *goweb.Context) error {
    project := &projectService.Project{
        ID: c.QueryInt("id"),
    }
    if err := project.Get(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, project)
    return nil
}

func deleteProject(c *goweb.Context) error {
    project := &projectService.Project{
        ID: c.PostFormInt("id"),
    }
    if err := project.Get(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }

    if project.Status != 0 {
        syncd.RenderAppError(c, "project delete falied, project status must be unavailable")
        return nil
    }

    if err := project.Delete(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }

    syncd.RenderJson(c, nil)
    return nil
}
