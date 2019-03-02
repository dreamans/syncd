// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/module/project"
    "github.com/dreamans/syncd/module/deploy"
    "github.com/dreamans/syncd/integrate"
    "github.com/dreamans/syncd/util/command"
    "github.com/dreamans/syncd/util/gostring"
)

func BuildStop(c *gin.Context) {
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
    integrate.StopBuild(id)
    render.JSON(c, nil)
}

func BuildStatus(c *gin.Context) {
    id := gostring.Str2Int(c.Query("id"))
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
    if err := build.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    var output []*command.TaskResult
    gostring.JsonDecode([]byte(build.Output), &output)

    render.JSON(c, map[string]interface{}{
        "apply_id": id,
        "status": build.Status,
        "start_time": build.StartTime,
        "finish_time": build.FinishTime,
        "tar": build.Tar,
        "output": output,
        "errmsg": build.Errmsg,
        "ctime": build.Ctime,
    })
}

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

    //apply must audit passed
    if apply.AuditStatus != deploy.AUDIT_STATUS_OK {
        render.AppError(c, "apply audit_status must passed")
        return
    }

    // apply status checked
    if apply.Status != deploy.APPLY_STATUS_NONE && apply.Status != deploy.APPLY_STATUS_FAILED {
        render.AppError(c, "apply status must none or failed")
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
        if (build.Status == deploy.BUILD_STATUS_START) {
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
    integrate.NewBuild(build.ApplyId, workSpace, packFile, p.BuildScript, repo, func(applyId int, err error, tar string, output string) {
        status := deploy.BUILD_STATUS_SUCCESS
        errmsg := ""
        if err != nil {
            status = deploy.BUILD_STATUS_FAILED
            errmsg = err.Error()
        }
        b := deploy.Build{
            ApplyId: applyId,
            Status: status,
            Tar: tar,
            Output: output,
            Errmsg: errmsg,
        }
        b.Finish()
    })

    render.JSON(c, nil)
}


