// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    repoService "github.com/tinystack/syncd/service/repo"
    projectService "github.com/tinystack/syncd/service/project"
    taskService "github.com/tinystack/syncd/service/task"
)

func RepoReset(c *goweb.Context) error {
    id := c.PostFormInt("id")
    if id == 0 {
        return syncd.RenderParamError("id can not empty")
    }
    project := &projectService.Project{
        ID: id,
    }
    if err := project.Detail(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    repo := &repoService.Repo{
        ID: id,
        Url: project.RepoUrl,
    }
    var err error
    if repo, err = repoService.RepoNew(repo); err != nil {
        return syncd.RenderAppError(err.Error())
    }

    taskTimeout := 60
    task := taskService.TaskCreate(taskService.TASK_REPO_RESET, []string{
        "pwd",
        repo.ResetRepo(),
    }, taskTimeout)
    c.CloseCallback(func() {
        task.Terminate()
    }, taskTimeout)

    task.TaskRun()
    if task.LastError() != nil {
        errMsg := task.Stdout() + task.Stderr()
        return syncd.RenderTaskError(errMsg)
    }

    return syncd.RenderJson(c, nil)
}

