// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "strings"
    "errors"
    "fmt"

    "github.com/tinystack/goutil/gostring"
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    projectService "github.com/tinystack/syncd/service/project"
    deployService "github.com/tinystack/syncd/service/deploy"
    userService "github.com/tinystack/syncd/service/user"
    repoService "github.com/tinystack/syncd/service/repo"
    serverService "github.com/tinystack/syncd/service/server"
    taskService "github.com/tinystack/syncd/service/task"
)

func DeployStop(c *goweb.Context) error {
    id := c.PostFormInt("id")
    apply, err := deployService.ApplyGetByPk(id)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if err := deployCheckUserPriv(c, apply); err != nil {
        return err
    }
    if apply.Status != deployService.APPLY_STATUS_DEPLOY_ING {
        return syncd.RenderAppError("apply status wrong")
    }
    apply.Status = deployService.APPLY_STATUS_DEPLOY_FAILED
    if err := apply.UpdateStatus(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

func DeployStatus(c *goweb.Context) error {
    id := c.QueryInt("id")
    apply, err := deployService.ApplyGetByPk(id)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if err := deployCheckUserPriv(c, apply); err != nil {
        return err
    }

    deploy := &deployService.DeployTask{
        ApplyId: apply.ID,
    }
    taskList, err := deploy.GetTaskItem()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    newTaskList := []map[string]interface{}{}
    for _, l := range taskList {
        newTaskList = append(newTaskList, map[string]interface{}{
            "id": l.ID,
            "name": l.Name,
            "status": l.Status,
            "level": l.Level,
            "output": l.Output,
        })
    }

    return syncd.RenderJson(c, goweb.JSON{
        "apply_status": apply.Status,
        "deploy_list": newTaskList,
    })
}

func DeployStart(c *goweb.Context) error {
    id := c.PostFormInt("id")
    apply, err := deployService.ApplyGetByPk(id)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if err := deployCheckUserPriv(c, apply); err != nil {
        return err
    }
    if apply.Status != deployService.APPLY_STATUS_AUDIT_PASS &&
    apply.Status != deployService.APPLY_STATUS_DEPLOY_FAILED {
        return syncd.RenderAppError("apply status wrong")
    }

    project, err := projectService.ProjectGetByPk(apply.ProjectId)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }

    repo := &repoService.Repo{
        ID: apply.ProjectId,
        ApplyId: apply.ID,
        Url: apply.RepoData.RepoUrl,
    }
    if repo, err = repoService.RepoNew(repo); err != nil {
        return syncd.RenderAppError(err.Error())
    }

    var updateRepoCmd, update2CommitCmd string
    if apply.RepoData.RepoMode == 1 {
        updateRepoCmd, err = repo.UpdateRepo(apply.RepoData.RepoBranch)
        if err != nil {
            return syncd.RenderAppError(err.Error())
        }
        tmpCommit := strings.Split(apply.RepoData.Commit, " - ")
        if len(tmpCommit) == 0 {
            return syncd.RenderAppError("git commit version  wrong")
        }
        commit := tmpCommit[0]
        update2CommitCmd = repo.Update2CommitRepo(apply.RepoData.RepoBranch, commit)
    } else {
        updateRepoCmd, err = repo.UpdateRepo("")
        if err != nil {
            return syncd.RenderAppError(err.Error())
        }
        update2CommitCmd = repo.Update2CommitRepo(apply.RepoData.Tag, "")
    }

    deployTasks := []*deployService.DeployTask{
        &deployService.DeployTask{
            ApplyId: apply.ID,
            Level: deployService.DEPLOY_LEVEL_UPDATE_REPO,
            Cmd: updateRepoCmd,
            Name: "fetch_repo_code_files",
        },
        &deployService.DeployTask{
            ApplyId: apply.ID,
            Level: deployService.DEPLOY_LEVEL_UPDATE_COMMIT,
            Cmd: update2CommitCmd,
            Name: "reset_repo_to_version",
        },
    }

    //tar zcvf 
    exFiles := gostring.StrFilterSliceEmpty(strings.Split(project.ExcludeFiles, "\n"))
    packRepoCmd := repo.PackRepo(exFiles)
    deployTasks = append(deployTasks, &deployService.DeployTask{
        ApplyId: apply.ID,
        Level: deployService.DEPLOY_LEVEL_PACK_REPO,
        Cmd: packRepoCmd,
        Name: "pack_repo_code_files",
    })

    srvList, err := serverService.ServerGetListByGroupIds(project.DeployServer)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    for _, srv := range srvList {
        deployCmds := repo.DeployRepo(gostring.Int2Str(srv.SshPort), srv.Ip, project.DeployUser, project.DeployPath, project.PreDeployCmd, project.PostDeployCmd)
        deployCmd, _ := gostring.JsonEncode(deployCmds)
        deployTasks = append(deployTasks, &deployService.DeployTask{
            ApplyId: apply.ID,
            Level: deployService.DEPLOY_LEVEL_DEPLOY,
            Cmd: deployCmd,
            Name: gostring.JoinStrings(srv.Name,"-",srv.Ip),
        })
    }

    depTask := &deployService.DeployTask{
        ApplyId: apply.ID,
    }
    if err := depTask.Flush(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    for _, dt := range deployTasks {
        if _, err := dt.Create(); err != nil {
            return syncd.RenderAppError(err.Error())
        }
    }

    deployRecordInfoLog(fmt.Sprintf("deploy_start, apply_id[%d]", apply.ID))

    go func(id int) {
        if err := deployRunTask(id); err != nil {
            deployRecordErrorLog(fmt.Sprintf("deploy_error, apply_id[%d], errmsg[%s]", apply.ID, err.Error()))
        } else {
            deployRecordInfoLog(fmt.Sprintf("deploy_finish, apply_id[%d]", apply.ID))
        }
    }(apply.ID)

    return syncd.RenderJson(c, nil)
}

func deployCheckUserPriv(c *goweb.Context, apply *deployService.Apply) error {
    if havePriv := userService.PrivIn(userService.DEPLOY_DEPLOY_ALL, c.GetIntSlice("priv")); !havePriv {
        if apply.UserId != c.GetInt("user_id") {
            return syncd.RenderCustomerError(syncd.CODE_ERR_NO_PRIV, "no priv")
        }
    }
    if err := applyCheckUserInSpace(apply.SpaceId, c.GetInt("user_id")); err != nil {
        return err
    }
    return nil
}

func deployRunTask(id int) error {
    apply := &deployService.Apply{
        ID: id,
        Status: deployService.APPLY_STATUS_DEPLOY_ING,
    }
    apply.UpdateStatus()
    deployTask := &deployService.DeployTask{
        ApplyId: id,
        Status: deployService.DEPLOY_STATUS_INIT,
    }
    taskList, err := deployTask.GetTaskItem()
    if err != nil {
        return err
    }

    var taskError error
    for _, tk := range taskList {
        if err := deployRunTaskItem(&tk); err != nil {
            taskError = err
            tk.Status = deployService.DEPOLY_STATUS_PANIC
            tk.UpdateStatus()
            deployRecordErrorLog(fmt.Sprintf("deploy_task_run_panic, apply_id[%d], level[%d], cmd[%s], name[%s], errmsg[%s]", apply.ID, tk.Level, tk.Cmd, tk.Name, err.Error()))
            break
        } else {
            deployRecordInfoLog(fmt.Sprintf("deploy_task_run_ok, apply_id[%d], level[%d], cmd[%s], name[%s]", apply.ID, tk.Level, tk.Cmd, tk.Name))
        }
    }

    if taskError != nil {
        apply.Status = deployService.APPLY_STATUS_DEPLOY_FAILED
        apply.ErrorLog = taskError.Error()
    } else {
        apply.Status = deployService.APPLY_STATUS_DEPLOY_SUCCESS
    }
    apply.UpdateStatus()

    return nil
}

func deployRunTaskItem(dt *deployService.DeployTask) error {
    // check task whether needs to be terminated
    applyDetail, err := deployService.ApplyGetByPk(dt.ApplyId)
    if err != nil {
        return err
    }
    if applyDetail.Status == deployService.APPLY_STATUS_DEPLOY_FAILED {
        dt.Status = deployService.DEPLOY_STATUS_STOP
        dt.UpdateStatus()
        return errors.New("deploy task is terminated by user")
    }

    dt.Status = deployService.DEPLOY_STATUS_START
    dt.UpdateStatus()

    var task *taskService.Task
    var taskTimeout = 60
    if dt.Level != deployService.DEPLOY_LEVEL_DEPLOY {
        task = taskService.TaskCreate(taskService.TASK_REPO_DEPLOY, []string{
            dt.Cmd,
        }, taskTimeout)
    } else {
        var cmds []string
        gostring.JsonDecode(dt.Cmd, &cmds)
        task = taskService.TaskCreate(taskService.TASK_REPO_DEPLOY, cmds, taskTimeout)
    }

    task.TaskRun()

    if err := task.LastError(); err != nil {
        var errMsg []string
        if s := err.Error(); s != "" {
            errMsg = append(errMsg, s)
        }
        if s := task.Stdout(); s != "" {
            errMsg = append(errMsg, s)
        }
        if s := task.Stderr(); s != "" {
            errMsg = append(errMsg, s)
        }
        errMessage := gostring.JoinSepStrings("\n", errMsg...)
        dt.Output = errMessage
        dt.Status = deployService.DEPOLY_STATUS_PANIC
        dt.UpdateStatus()
        return errors.New(errMessage)
    }
    dt.Output = task.Stdout()
    dt.Status = deployService.DEPLOY_STATUS_END
    dt.UpdateStatus()

    return nil
}

func deployRecordErrorLog(msg string) {
    syncd.Logger.Error("DEPLOY_TASK_ERROR, %s", msg)
}

func deployRecordInfoLog(msg string) {
    syncd.Logger.Info("DEPLOY_TASK_INFO, %s", msg)
}

