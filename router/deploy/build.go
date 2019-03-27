// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/router/common"
    "github.com/dreamans/syncd/module/project"
    "github.com/dreamans/syncd/module/deploy"
    buiTask "github.com/dreamans/syncd/build"
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
    buiTask.StopTask(id)
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
    if build.Status == deploy.BUILD_STATUS_START {
        _, output, _ = buiTask.StatusTask(id)
    } else {
        gostring.JsonDecode([]byte(build.Output), &output)
    }

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

    // apply must audit passed
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

    // create build task
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

    repo := buiTask.NewRepo(p.RepoUrl, workSpace)
    if apply.BranchName != "" {
        repo.SetBranch(apply.BranchName)
    }
    if apply.CommitVersion != "" {
        repo.SetCommit(apply.CommitVersion)
    }

    bui, err := buiTask.NewBuild(repo, workSpace, syncd.App.LocalTmpSpace, packFile, p.BuildScript)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    buiTask.NewTask(id, bui, func(id int, packFile string, result *buiTask.Result, taskResult []*command.TaskResult) {
        status := deploy.BUILD_STATUS_SUCCESS
        errmsg := ""
        if err := result.GetError(); err != nil {
            status = deploy.BUILD_STATUS_FAILED
            errmsg = err.Error()
            packFile = ""
        }
        output := string(gostring.JsonEncode(taskResult))
        b := deploy.Build{
            ApplyId: id,
            Status: status,
            Tar: packFile,
            Output: output,
            Errmsg: errmsg,
        }
        b.Finish()

        // run hook script
        common.HookBuild(id)
    })
    render.JSON(c, nil)
}