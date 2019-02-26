// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/util/gostring"
    "github.com/dreamans/syncd/module/project"
    "github.com/dreamans/syncd/module/deploy"
    "github.com/dreamans/syncd/integrate"
)

func BuildStart(c *gin.Context) {
    id := gostring.Str2Int(c.PostForm("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    apply := &deploy.Apply{
        ID: id,
    }
    if err := apply.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    m := &project.Member{
        UserId: c.GetInt("user_id"),
        SpaceId: apply.SpaceId,
    }
    if in := m.MemberInSpace(); !in {
        render.CustomerError(c, render.CODE_ERR_NO_PRIV, "user is not in the project space")
        return
    }

    build := &deploy.Build{
        ApplyId: id,
    }
    var (
        exists bool
        err error
    )
    if exists, err = build.Exists(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    if exists {
        if (build.Status != deploy.BUILD_STATUS_FAILED) {
            render.AppError(c, "deploy build status wrong")
            return
        }
        if err := build.Delete(); err != nil {
            render.AppError(c, err.Error())
            return
        }
    }

    p := &project.Project{
        ID: apply.ProjectId,
    }
    if err := p.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    build.ApplyId = id
    build.Status = deploy.BUILD_STATUS_START
    build.Cmd = p.BuildScript
    if err := build.Create(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    // build
    workSpace := gostring.JoinStrings(syncd.App.LocalTmpSpace, "/", gostring.Int2Str(build.ApplyId))
    packFile := gostring.JoinStrings(syncd.App.LocalTarSpace, "/", gostring.Int2Str(build.ApplyId), ".tgz")
    repo := &integrate.Repo{
        Url: p.RepoUrl,
        Branch: apply.BranchName,
        Commit: apply.CommitVersion,
        Local: workSpace,
    }
    b := integrate.NewBuild(workSpace, packFile, p.BuildScript, repo)
    b.Build()

    render.JSON(c, nil)
}


