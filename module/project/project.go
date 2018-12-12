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
    "github.com/tinystack/syncd/model"
)

func init() {
    route.Register(route.API_PROJECT_NEW, newProject)
}

type NewProjectParam struct {
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

func newProject(c *goweb.Context) error {
    p := NewProjectParam{
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
    if valid := govalidate.NewValidate(&p); !valid.Pass() {
        syncd.RenderParamError(c, valid.LastFailed().Msg)
        return nil
    }

    deployServer := strings.Join(p.DeployServer, ",")
    needAudit := 0
    if c.PostFormInt("needAudit") != 0 {
        needAudit = 1
    }
    status := 0
    if c.PostFormInt("status") != 0 {
        status = 1
    }

    projectModel := model.Project{
        Name: p.Name,
        Description: p.Description,
        Space: p.Space,
        Repo: p.Repo,
        RepoUrl: p.RepoUrl,
        DeployServer: deployServer,
        DeployUser: p.DeployUser,
        DeployPath: p.DeployPath,
        DeployHistory: p.DeployHistory,
        PreDeployCmd: c.PostForm("preDeployCmd"),
        PostDeployCmd: c.PostForm("postDeployCmd"),
        NeedAudit: needAudit,
        Status: status,
        RepoUser: c.PostForm("repoUser"),
        RepoPass: c.PostForm("repoPass"),
        RepoMode: c.PostFormInt("repoMode"),
        BuildScript: c.PostForm("buildScript"),
    }
    id, ok := projectModel.Create()
    if !ok {
        syncd.RenderAppError(c, "项目新增失败")
        return nil
    }
    syncd.RenderJson(c, goweb.JSON{
        "id": id,
    })
    return nil
}
