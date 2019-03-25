// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "fmt"

    "github.com/dreamans/syncd"
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/router/common"
    "github.com/dreamans/syncd/util/gostring"
    "github.com/dreamans/syncd/module/deploy"
    "github.com/dreamans/syncd/module/project"
    "github.com/dreamans/syncd/module/server"
    "github.com/dreamans/syncd/module/user"
    depTask "github.com/dreamans/syncd/deploy"
)

func DeployRollback(c *gin.Context) {
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

    if apply.RollbackId == 0 {
        render.NoDataError(c, "rollback apply not exists")
        return
    }

    oldApply := &deploy.Apply{
        ID: apply.RollbackId,
    }
    if err := oldApply.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    if oldApply.ProjectId != apply.ProjectId {
        render.NoDataError(c, "rollback apply not exists")
        return
    }

    //old build order
    oldBuild := &deploy.Build{
        ApplyId: oldApply.ID,
    }
    if err := oldBuild.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    //create rollback apply order
    rollbackApply := &deploy.Apply{
        SpaceId: oldApply.SpaceId,
        ProjectId: oldApply.ProjectId,
        Name: fmt.Sprintf("rollback: (ID:%d) %s", apply.ID, apply.Name),
        Description: fmt.Sprintf("rollback \"(ID:%d) %s\" to \"(ID:%d) %s\"", apply.ID, apply.Name, oldApply.ID, oldApply.Name),
        BranchName: oldApply.BranchName,
        CommitVersion: oldApply.CommitVersion,
        Status: deploy.APPLY_STATUS_NONE,
        UserId: c.GetInt("user_id"),
        AuditStatus: deploy.AUDIT_STATUS_OK,
        IsRollbackApply: 1,
        RollbackApplyId: apply.ID,
    }
    if err := rollbackApply.Create(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    //create rollback build order
    rollbackBuild := &deploy.Build{
        ApplyId: rollbackApply.ID, 
        StartTime: oldBuild.StartTime, 
        FinishTime: oldBuild.FinishTime, 
        Status: oldBuild.Status,
        Tar: oldBuild.Tar,
        Output: oldBuild.Output, 
        Errmsg: oldBuild.Errmsg, 
    }
    if err := rollbackBuild.CreateFull(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    // update current apply rollback_apply_id
    apply.RollbackApplyId = rollbackApply.ID
    if err := apply.UpdateRollback(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    render.JSON(c, rollbackApply)
}

func DeployStop(c *gin.Context) {
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

    if !depTask.ExistsTask(id) {
        apply := &deploy.Apply{
            ID: id,
            Status: deploy.APPLY_STATUS_FAILED,
        }
        apply.UpdateStatus()
    } else {
        depTask.StopTask(id)
    }
    render.JSON(c, nil)
}

func DeployStatus(c *gin.Context) {
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

    taskRest := []map[string]interface{}{}
    depResults := depTask.StatusTask(id)
    if depResults != nil  {
        for _, depResult := range depResults {
            groupStatus := deploy.DEPLOY_STATUS_NONE
            switch depResult.Status {
            case depTask.STATUS_INIT:
                groupStatus = deploy.DEPLOY_STATUS_NONE
            case depTask.STATUS_ING:
                groupStatus = deploy.DEPLOY_STATUS_START
            case depTask.STATUS_DONE:
                groupStatus = deploy.DEPLOY_STATUS_SUCCESS
            case depTask.STATUS_FAILED:
                groupStatus = deploy.DEPLOY_STATUS_FAILED
            }

            srvRest := []map[string]interface{}{}
            for _, r := range depResult.ServerRest {
                var err string
                if e := r.Error; e != nil {
                    err = e.Error()
                }
                srvRest = append(srvRest, map[string]interface{}{
                    "id": r.ID,
                    "task": r.TaskResult,
                    "status": r.Status,
                    "error": err, 
                })
            }
            groupRest := map[string]interface{}{
                "group_id": depResult.ID,
                "status": groupStatus,
                "content": srvRest,
            }
            taskRest = append(taskRest, groupRest)
        }
    } else {
        d := &deploy.Deploy{
            ApplyId: id,
        }
        taskList, err := d.TaskList()
        if err != nil {
            render.AppError(c, err.Error())
            return
        }
        for _, l := range taskList {
            var obj interface{}
            gostring.JsonDecode([]byte(l.Content), &obj)

            groupRest := map[string]interface{}{
                "group_id": l.GroupId,
                "status": l.Status,
                "content": obj,
            }
            taskRest = append(taskRest, groupRest)
        }
    }
    render.JSON(c, map[string]interface{}{
        "status": apply.Status,
        "task_list": taskRest,
    })
}

func DeployStart(c *gin.Context) {
    id := gostring.Str2Int(c.PostForm("id"))
    isParallel := gostring.Str2Int(c.PostForm("is_parallel"))

    if id == 0 {
        render.ParamError(c, "apply id cannot be empty")
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

    if apply.Status != deploy.APPLY_STATUS_NONE && apply.Status != deploy.APPLY_STATUS_FAILED {
        render.AppError(c, "deploy apply have deployed success")
        return
    }

    //check project have deploying apply
    if canDeploy, err := apply.CheckHaveDeploying(); !canDeploy || err != nil {
        if err != nil {
            render.AppError(c, err.Error())
            return
        }
        render.RepeatError(c, "project have deploying apply within 24 hours")
        return
    }

    //apply must audit passed
    if apply.AuditStatus != deploy.AUDIT_STATUS_OK {
        render.AppError(c, "apply audit_status must passed")
        return
    }

    // check is building
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
    if !exists ||  build.Status != deploy.BUILD_STATUS_SUCCESS {
        render.NoDataError(c, "have no build tar package file")
        return
    }

    proj := &project.Project{
        ID: apply.ProjectId,
    }
    if err := proj.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    // cluster servers
    serverList, err := server.ServerGetListByGroupIds(proj.OnlineCluster)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    d := &deploy.Deploy{
        ApplyId: id,
    }
    if err := d.DeleteByApplyId(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    groupSrvs := map[int][]server.Server{}
    for _, srv := range serverList {
        groupSrvs[srv.GroupId] = append(groupSrvs[srv.GroupId], srv)
    }

    deploys := []*depTask.Deploy{}
    for _, gid := range proj.OnlineCluster {
        gsrv, exists := groupSrvs[gid]
        if !exists {
            continue
        }
        dep := &depTask.Deploy{
            ID: gid,
            User: proj.DeployUser,
            PreCmd: proj.PreDeployCmd,
            PostCmd: proj.AfterDeployCmd,
            DeployPath: proj.DeployPath,
            DeployTmpPath: syncd.App.RemoteSpace,
            PackFile: build.Tar,
        }
        for _, srv := range gsrv {
            dep.AddServer(srv.ID, srv.Ip, srv.SSHPort)
        }
        deploys = append(deploys, dep)

        // Write task init to DB
        d := &deploy.Deploy{
            ApplyId: id,
            GroupId: gid,
            Status: deploy.DEPLOY_STATUS_NONE,
        }
        if err := d.Create(); err != nil {
            render.AppError(c, err.Error())
            return
        }
    }

    apply.Status = deploy.APPLY_STATUS_ING
    if err := apply.UpdateStatus(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    startDeployFn := func(id, gid, status int, serverResult []*depTask.ServerResult) {
        d := &deploy.Deploy{
            ApplyId: id,
            GroupId: gid,
            Status: deploy.DEPLOY_STATUS_START,
        }
        d.UpdateStatus()
    }
    finishDeployFn := func(id, gid, status int, serverResult []*depTask.ServerResult) {
        taskStatus := deploy.DEPLOY_STATUS_SUCCESS
        if status == depTask.STATUS_FAILED {
            taskStatus = deploy.DEPLOY_STATUS_FAILED
        }
        srvRest := []map[string]interface{}{}
        for _, r := range serverResult {
            err := ""
            if e := r.Error; e != nil {
                err = e.Error()
            }
            srvRest = append(srvRest, map[string]interface{}{
                "id": r.ID,
                "task": r.TaskResult,
                "status": r.Status,
                "error": err, 
            })
        }
        d := &deploy.Deploy{
            ApplyId: id,
            GroupId: gid,
            Status: taskStatus,
            Content: string(gostring.JsonEncode(srvRest)),
        }
        d.UpdateResult()
    }
    taskFn := func(id, status int) {
        apply := &deploy.Apply{
            ID: id,
        }
        var deployStatus int
        if status == depTask.STATUS_FAILED {
            apply.Status = deploy.APPLY_STATUS_FAILED
            deployStatus = MAIL_STATUS_FAILED
        } else {
            apply.Status = deploy.APPLY_STATUS_SUCCESS
            deployStatus = MAIL_STATUS_SUCCESS
        }
        apply.UpdateStatus()

        //send deploy email
        if err := apply.Detail(); err != nil {
            return
        }
        u := &user.User{
            ID: apply.UserId,
        }
        if err := u.Detail(); err != nil {
            return
        }
        proj := &project.Project{
            ID: apply.ProjectId,
        }
        if err := proj.Detail(); err != nil {
            return
        }
        mails := gostring.JoinStrings(proj.DeployNotice, ",", u.Email)
        MailSend(&MailMessage{
            Mail: mails,
            ApplyId: id,
            Mode: MAIL_MODE_DEPLOY,
            Status: deployStatus,
            Title: apply.Name,
        })

        // run hook script
        common.HookDeploy(id)
    }
    if err := depTask.NewTask(
        id, 
        isParallel, 
        deploys, 
        startDeployFn, 
        finishDeployFn, 
        taskFn,
    ); err != nil {
        render.AppError(c, err.Error())
        return
    }

    render.JSON(c, nil)
}