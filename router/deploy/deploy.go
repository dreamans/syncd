// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/util/gostring"
    "github.com/dreamans/syncd/module/deploy"
    "github.com/dreamans/syncd/module/project"
    "github.com/dreamans/syncd/module/server"
    dep "github.com/dreamans/syncd/deploy"
)

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
    d := &deploy.Deploy{
        ApplyId: id,
    }
    taskList, err := d.TaskList()
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    render.JSON(c, map[string]interface{}{
        "status": apply.Status,
        "task_list": taskList,
    })
}

func DeployStart(c *gin.Context) {
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

    if apply.Status != deploy.APPLY_STATUS_NONE && apply.Status != deploy.APPLY_STATUS_FAILED {
        render.AppError(c, "deploy apply status wrong")
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
        render.NoDataError(c, err.Error())
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

    deployGroups := []*dep.Deploy{}
    for _, gid := range proj.OnlineCluster {
        srvs := []*dep.Server{}
        for _, srv := range serverList {
            if srv.GroupId == gid {
                srvs = append(srvs, dep.NewServer(srv.ID, srv.SSHPort, srv.Ip, proj.DeployUser, ""))
            }
        }
        depGroup, err := dep.NewDepoly(&dep.Deploy{
            ID: id,
            Srvs: srvs,
            PreCmd: proj.PreDeployCmd,
            PostCmd: proj.AfterDeployCmd,
            DeployPath: proj.DeployPath,
            PackFile: build.Tar,
            CallbackFn: func(id int, srvId int, srvStatus *dep.ServerStatus) {
                d := &deploy.Deploy{
                    ApplyId: id,
                    ServerId: srvId,
                    Status: srvStatus.Status,
                    Output: string(gostring.JsonEncode(srvStatus.TaskResult)),
                    FinishTime: int(time.Now().Unix()),
                }
                if srvStatus.Error != nil {
                    d.Errmsg = srvStatus.Error.Error()
                }
                d.UpdateResult()
            },
            StartCallbackFn: func(id int, srvId int, srvStatus *dep.ServerStatus) {
                d := &deploy.Deploy{
                    ApplyId: id,
                    ServerId: srvId,
                    Status: deploy.DEPLOY_STATUS_START,
                    StartTime: int(time.Now().Unix()),
                }
                d.UpdateResult()
            },
        })
        if err != nil {
            render.AppError(c, err.Error())
            return 
        }
        deployGroups = append(deployGroups, depGroup)
    }

    // Write task to DB
    for _, s := range serverList {
        d := &deploy.Deploy{
            ApplyId: id,
            GroupId: s.GroupId,
            ServerId: s.ID,
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

    // deployGroups
    dep.DeployGroups(id, deployGroups, func(id int, haveError bool) {
        apply := &deploy.Apply{
            ID: id,
        }
        if haveError {
            apply.Status = deploy.APPLY_STATUS_FAILED
        } else {
            apply.Status = deploy.APPLY_STATUS_SUCCESS
        }
        apply.UpdateStatus()
    })

    render.JSON(c, nil)
}
