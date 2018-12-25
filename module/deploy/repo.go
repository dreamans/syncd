// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    "github.com/tinystack/syncd/route"
    projectService "github.com/tinystack/syncd/service/project"
    deployService "github.com/tinystack/syncd/service/deploy"
)

func init() {
    route.Register(route.API_DEPLOY_REPO_TAGLIST, tagListRepo)
    route.Register(route.API_DEPLOY_REPO_RESET, resetRepo)
    route.Register(route.API_DEPLOY_REPO_COMMITLIST, commitListRepo)
}

func resetRepo(c *goweb.Context) error {
    project := &projectService.Project{
        ID: c.PostFormInt("project_id"),
    }
    if err := project.Get(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    repo, err := deployService.NewRepo(&deployService.Repo{
        ID: project.ID,
        Repo: project.Repo,
        Url: project.RepoUrl,
        User: project.RepoUser,
        Pass: project.RepoPass,
    })
    if err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }

    if err := repo.ResetRepo(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }

    syncd.RenderJson(c, nil)
    return nil
}

func tagListRepo(c *goweb.Context) error {
    project := &projectService.Project{
        ID: c.QueryInt("id"),
    }
    if err := project.Get(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }

    repo, err := deployService.NewRepo(&deployService.Repo{
        ID: project.ID,
        Repo: project.Repo,
        Url: project.RepoUrl,
        User: project.RepoUser,
        Pass: project.RepoPass,
    })
    if err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }

    list, err := repo.TagListRepo()
    if err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, goweb.JSON{
        "list": list,
    })
    return nil
}

func commitListRepo(c *goweb.Context) error {
    return nil
}
