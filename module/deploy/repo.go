// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    deployService "github.com/tinystack/syncd/service/deploy"
    projectService "github.com/tinystack/syncd/service/project"
)

func RepoReset(c *goweb.Context) error {
    repo, err := deployServiceRepo(c.PostFormInt("project_id"))
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if err := repo.ResetRepo(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

func RepoTagList(c *goweb.Context) error {
    repo, err := deployServiceRepo(c.QueryInt("id"))
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }

    list, err := repo.TagListRepo()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "list": list,
    })
}

func RepoCommitList(c *goweb.Context) error {
    repo, err := deployServiceRepo(c.QueryInt("id"))
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    list, err := repo.CommitListRepo()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "list": list,
    })
}

func deployServiceRepo(id int) (*deployService.Repo, error) {
    project := &projectService.Project{
        ID: id,
    }
    if err := project.Get(); err != nil {
        return nil, err
    }
    repo, err := deployService.NewRepo(&deployService.Repo{
        ID: project.ID,
        Repo: project.Repo,
        Url: project.RepoUrl,
        User: project.RepoUser,
        Pass: project.RepoPass,
    })
    if err != nil {
        return nil, err
    }
    return repo, nil
}
