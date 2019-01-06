// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "fmt"

    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    projectService "github.com/tinystack/syncd/service/project"
    serverService "github.com/tinystack/syncd/service/server"
    taskService "github.com/tinystack/syncd/service/task"
)

func ServerCheck(c *goweb.Context) error {
    id := c.QueryInt("id")
    project, err := projectService.ProjectGetByPk(id)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    serList, err := serverService.ServerGetListByGroupIds(project.DeployServer)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }

    var sshCmds []string
    for _, ser := range serList {
        sshCmds = append(sshCmds, fmt.Sprintf("/usr/bin/env ssh -p %d %s@%s 'mkdir -p %s; cd %s; echo \"%s\" $( pwd ) $( whoami ) $( date +\"%%Y-%%m-%%d %%H:%%M:%%S\" ) \"[检测成功]\"'", ser.SshPort, project.DeployUser, ser.Ip, project.DeployPath, project.DeployPath, ser.Ip))
    }

    task := taskService.TaskCreate(taskService.TASK_SERVER_CHECK, sshCmds)
    c.CloseCallback(func() {
        task.Terminate()
    }, len(sshCmds) * 10)

    task.TaskRun()
    if task.LastError() != nil {
        errMsg := task.Stderr() + task.LastError().Error()
        return syncd.RenderTaskError(errMsg)
    }

    return syncd.RenderJson(c, goweb.JSON{
        "srv_list": serList,
        "output": task.Stdout(),
    })
}
