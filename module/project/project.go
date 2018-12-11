// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
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
    PreDeployCmd    string
    PostDeployCmd   string
    NeedAudit       int
    Status          int
    RepoUser        string
    RepoPass        string
    RepoMode        int
    BuildScript     string
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
        PreDeployCmd: c.PostForm("preDeployCmd"),
        PostDeployCmd: c.PostForm("postDeployCmd"),
        NeedAudit: c.PostFormInt("needAudit"),
        Status: c.PostFormInt("status"),
        RepoUser: c.PostForm("repoUser"),
        RepoPass: c.PostForm("repoPass"),
        RepoMode: c.PostFormInt("repoMode"),
        BuildScript: c.PostForm("buildScript"),
    }
    if valid := govalidate.NewValidate(&p); !valid.Pass() {
        syncd.RenderParamError(c, valid.LastFailed().Msg)
        return nil
    }

    projectModel := model.Project{
        Name: p.Name,
        Description: p.Description,
        Space: p.Space,
        BuildScript: p.BuildScript,
    }
    projectModel.Create()

    syncd.RenderJson(c, p)
    return nil
}
