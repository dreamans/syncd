// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "strings"

    "github.com/tinystack/goweb"
    "github.com/tinystack/govalidate"
    "github.com/tinystack/syncd/route"
    "github.com/tinystack/syncd"
    projectModel "github.com/tinystack/syncd/model/project"
)

func init() {
    route.Register(route.API_PROJECT_UPDATE, updateProject)
    route.Register(route.API_PROJECT_LIST, listProject)
    route.Register(route.API_PROJECT_GET, getProject)
}

type ProjectParamValid struct {
    Name            string      `valid:"required" errmsg:"required=Project name cannot be empty"`
    Description     string      `valid:"require" errmsg:"required=Project description cannot be empty"`
    Space           string      `valid:"require" errmsg:"required=Project space cannot be empty"`
    Repo            string      `valid:"require" errmsg:"required=Repo type cannot be empty"`
    RepoUrl         string      `valid:"require" errmsg:"required=Repo remote addr cannot be empty"`
    DeployServer    []string    `valid:"require" errmsg:"required=Deploy server cannot be empty"`
    DeployUser      string      `valid:"require" errmsg:"required=Deploy user cannot be epmty"`
    DeployPath      string      `valid:"require" errmsg:"required=Deploy path cannot be epmty"`
    DeployHistory   int         `valid:"int_min=3" errmsg:"int_min=Deploy history at least 3"`
}

func updateProject(c *goweb.Context) error {
    params := ProjectParamValid{
        Name: c.PostForm("name"),
        Description: c.PostForm("description"),
        Space: c.PostForm("space"),
        Repo: c.PostForm("repo"),
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

    deployServer := strings.Join(params.DeployServer, ",")
    needAudit := 0
    if c.PostFormInt("needAudit") != 0 {
        needAudit = 1
    }
    status := 0
    if c.PostFormInt("status") != 0 {
        status = 1
    }
    projectId := c.PostFormInt("id")
    p := projectModel.Project{
        Name: params.Name,
        Description: params.Description,
        Space: params.Space,
        Repo: params.Repo,
        RepoUrl: params.RepoUrl,
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
        RepoMode: c.PostFormInt("repoMode"),
        BuildScript: c.PostForm("buildScript"),
    }
    if projectId > 0 {
        p.ID = projectId
        if ok := projectModel.Update(projectId, p); !ok {
            syncd.RenderAppError(c, "project data update failed")
            return nil
        }
    } else {
        if ok := projectModel.Create(&p); !ok {
            syncd.RenderAppError(c, "project data create failed")
            return nil
        }
    }

    syncd.RenderJson(c, nil)
    return nil
}

func listProject(c *goweb.Context) error {
    var (
        ok      bool
        total   int
        offset  int
        limit   int
    )
    offset, limit = c.QueryInt("offset"), c.QueryInt("limit")
    projList, ok := projectModel.List("id, name, repo_mode, need_audit, status", offset, limit)
    if !ok {
        syncd.RenderAppError(c, "get project list data failed")
        return nil
    }

    total, ok = projectModel.Total()
    if !ok {
        syncd.RenderAppError(c, "get project total count failed")
        return nil
    }
    syncd.RenderJson(c, goweb.JSON{
        "list": projList,
        "total": total,
    })
    return nil
}

func getProject(c *goweb.Context) error {
    id := c.QueryInt("id")
    if id == 0 {
        syncd.RenderParamError(c, "id can not be empty")
        return nil
    }

    p, ok := projectModel.Get(id)
    if !ok {
        syncd.RenderAppError(c, "get project detail data failed")
        return nil
    }
    syncd.RenderJson(c, goweb.JSON{
        "detail": p,
    })
    return nil
}
