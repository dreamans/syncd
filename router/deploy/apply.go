// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/util/gostring"
    "github.com/dreamans/syncd/module/project"
    "github.com/dreamans/syncd/module/deploy"
)

type ApplyFormBind struct {
    ProjectId       int     `form:"project_id" binding:"required"`
    SpaceId         int     `form:"space_id" binding:"required"`
    Name            string  `form:"name" binding:"required"`
    BranchName      string  `form:"branch_name"`
    CommitVersion   string  `form:"commit_version"`
    Description     string  `form:"description"`
}

func ApplyProjectAll(c *gin.Context) {
    member := &project.Member{
        UserId: c.GetInt("user_id"),
    }
    spaceIds, err := member.SpaceIdsByUserId()
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    space := &project.Space{}
    spaceList, err := space.List(spaceIds, "", 0, 999)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    spaceMap := map[int]project.Space{}
    for _, l := range spaceList {
        spaceMap[l.ID] = l
    }

    projList, err := project.ProjectAllBySpaceIds(spaceIds)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    list := []map[string]interface{}{}
    for _, l := range projList {
        var (
            spaceId int
            spaceName string
        )
        if space, exists := spaceMap[l.SpaceId]; exists {
            spaceId, spaceName = space.ID, space.Name
        }
        list = append(list, map[string]interface{}{
            "space_id": spaceId,
            "project_id": l.ID,
            "project_name": l.Name,
            "space_name": spaceName,
        })
    }

    render.JSON(c, list)
}

func ApplySubmit(c *gin.Context) {
    var form ApplyFormBind
    if err := c.ShouldBind(&form); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    proj := &project.Project{
        ID: form.ProjectId,
    }
    if err := proj.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    branchName, commitVersion := form.BranchName, form.CommitVersion
    if proj.DeployMode == 1 {
        if proj.RepoBranch != "" {
            branchName = proj.RepoBranch
        }
    } else {
        commitVersion = ""
    }
    if branchName == "" {
        render.ParamError(c, "branch_name cannot be empty")
        return
    }

    apply := &deploy.Apply{
        SpaceId: form.SpaceId,
        ProjectId: form.ProjectId,
        Name: form.Name,
        Description: form.Description,
        BranchName: branchName,
        CommitVersion: commitVersion,
    }
    if err := apply.Create(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.Success(c)
}

func ApplyProjectDetail(c *gin.Context) {
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

    restProj := map[string]interface{}{
        "id": proj.ID,
        "name": proj.Name,
        "deploy_mode": proj.DeployMode,
        "repo_branch": proj.RepoBranch,
    }

    render.JSON(c, restProj)
}
