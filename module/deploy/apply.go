// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "github.com/tinystack/govalidate"
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    projectService "github.com/tinystack/syncd/service/project"
    deployService "github.com/tinystack/syncd/service/deploy"
)

type ApplyParamValid struct {
    ProjectId   int     `valid:"int_min=1" errmsg:"required=project_id cannot be empty"`
    Name        string  `valid:"required" errmsg:"required=name cannot be empty"`
    Description string  `valid:"required" errmsg:"required=name cannot be empty"`
}

func ApplySubmit(c *goweb.Context) error {
    offset, limit, keyword := c.QueryInt("offset"), c.QueryInt("limit"), c.Query("keyword")

    apply := deployService.Apply{}
    list, total, err := apply.List(keyword, offset, limit)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    var projectIds, spaceIds []int
    for _, l := range list {
        projectIds = append(projectIds, l.ProjectId)
        spaceIds = append(spaceIds, l.SpaceId)
    }
    projMaps, err := projectService.ProjectGetMapByIds(projectIds)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    spaceMaps, err := projectService.SpaceGetMapByIds(spaceIds)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }

    var newList []map[string]interface{}
    for _, l := range list {
        var projectName, spaceName string
        if proj, exists := projMaps[l.ProjectId]; exists {
            projectName = proj.Name
        }
        if space, exists := spaceMaps[l.SpaceId]; exists {
            spaceName = space.Name
        }
        newList = append(newList, map[string]interface{}{
            "id": l.ID,
            "name": l.Name,
            "project_name": projectName,
            "space_name": spaceName,
            "status": l.Status,
            "ctime": l.Ctime,
        })
    }

    return syncd.RenderJson(c, goweb.JSON{
        "list": newList,
        "total": total,
    })
}

func ApplyList(c *goweb.Context) error {
    params := ApplyParamValid{
        ProjectId: c.PostFormInt("project_id"),
        Name: c.PostForm("name"),
        Description: c.PostForm("description"),
    }
    if valid := govalidate.NewValidate(&params); !valid.Pass() {
        return syncd.RenderParamError(valid.LastFailed().Msg)
    }
    tag, commit := c.PostForm("tag"), c.PostForm("commit")

    project, err := projectService.ProjectGetByPk(params.ProjectId)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if project.Status != 1 {
        return syncd.RenderParamError("roject not enabled")
    }
    if project.RepoMode == 1 && commit == "" {
        return syncd.RenderParamError("commit can not be empty")
    }
    if project.RepoMode == 2 && tag == "" {
        return syncd.RenderParamError("tag can not be empty")
    }
    var status int
    if project.NeedAudit == 0 {
        status = 2
    }
    apply := &deployService.Apply{
        ProjectId: project.ID,
        SpaceId: project.SpaceId,
        Name: params.Name,
        Description: params.Description,
        Status: status,
        RepoData: deployService.ApplyRepoData{
            Repo: project.Repo,
            RepoUrl: project.RepoUrl,
            RepoUser: project.RepoUser,
            RepoPass: project.RepoPass,
            RepoMode: project.RepoMode,
            RepoBranch: project.RepoBranch,
            Tag: tag,
            Commit: commit,
        },
    }
    if err := apply.Create(); err != nil {
        return syncd.RenderAppError(err.Error())
    }

    return syncd.RenderJson(c, nil)
}
