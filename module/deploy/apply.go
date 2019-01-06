// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "strings"
    "time"

    "github.com/tinystack/govalidate"
    "github.com/tinystack/goutil/gostring"
    "github.com/tinystack/goutil/goslice"
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    projectService "github.com/tinystack/syncd/service/project"
    deployService "github.com/tinystack/syncd/service/deploy"
    repoService "github.com/tinystack/syncd/service/repo"
    taskService "github.com/tinystack/syncd/service/task"
    userService "github.com/tinystack/syncd/service/user"
    logService "github.com/tinystack/syncd/service/operate_log"
)

type ApplyParamValid struct {
    ProjectId   int     `valid:"int_min=1" errmsg:"required=project_id cannot be empty"`
    Name        string  `valid:"required" errmsg:"required=name cannot be empty"`
    Description string  `valid:"required" errmsg:"required=name cannot be empty"`
}

func ApplySpaceList(c *goweb.Context) error {
    userId := c.GetInt("user_id")
    list, err := projectService.SpaceGetListByUserId(userId)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "list": list,
    })
}

func ApplyProjectList(c *goweb.Context) error {
    spaceId, userId := c.QueryInt("space_id"), c.GetInt("user_id")
    projectUser := &projectService.User{
        SpaceId: spaceId,
        UserId: userId,
    }
    if exists, err := projectUser.CheckUserInSpace(); err != nil || !exists {
        return syncd.RenderAppError("user have no privilege to access space")
    }
    list, err := projectService.ProjectGetListBySpaceId(spaceId)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "list": list,
    })
}

func ApplyProjectAll(c *goweb.Context) error {
    userId := c.GetInt("user_id")
    spaceList, err := projectService.SpaceGetListByUserId(userId)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }

    spaceIds := []int{}
    spaceNameMap := map[int]string{}
    for _, l := range spaceList {
        spaceIds = append(spaceIds, l.ID)
        spaceNameMap[l.ID] = l.Name
    }
    projectList, err := projectService.ProjectGetListBySpaceIds(spaceIds)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }

    projList := []map[string]interface{}{}
    for _, pjl := range projectList {
        var spaceName string
        vv, exists := spaceNameMap[pjl.SpaceId]
        if exists {
            spaceName = vv
        }
        projList = append(projList, map[string]interface{}{
            "id": pjl.ID,
            "name": pjl.Name,
            "space_name": spaceName,
            "status": pjl.Status,
        })
    }

    return syncd.RenderJson(c, goweb.JSON{
        "list": projList,
    })
}

func ApplyProjectDetail(c *goweb.Context) error {
    project, err := applyCheckAndGetProjectDetail(c.QueryInt("id"), c.GetInt("user_id"))
    if err != nil {
        return err
    }
    return syncd.RenderJson(c, goweb.JSON{
        "id": project.ID,
        "name": project.Name,
        "repo_mode": project.RepoMode,
        "repo_branch": project.RepoBranch,
    })
}

func ApplyRepoTagList(c *goweb.Context) error {
    project, err := applyCheckAndGetProjectDetail(c.QueryInt("id"), c.GetInt("user_id"))
    if err != nil {
        return err
    }
    repo := &repoService.Repo{
        ID: project.ID,
        Url: project.RepoUrl,
    }
    if repo, err = repoService.RepoNew(repo); err != nil {
        return syncd.RenderAppError(err.Error())
    }

    updateRepoCmd, err := repo.UpdateRepo("")
    if err != nil {
        return syncd.RenderTaskError(err.Error())
    }
    taskUpdateRepo := taskService.TaskCreate(taskService.TASK_REPO_UPDATE, []string{
        updateRepoCmd,
    })

    taskTagListRepo := taskService.TaskCreate(taskService.TASK_REPO_TAG_LIST, []string{
        repo.TagListRepo(),
    })

    c.CloseCallback(func() {
        taskUpdateRepo.Terminate()
        taskTagListRepo.Terminate()
    }, 60)

    taskUpdateRepo.TaskRun()
    if taskUpdateRepo.LastError() != nil {
        return syncd.RenderTaskError(taskUpdateRepo.Stdout() + taskUpdateRepo.Stderr())
    }

    taskTagListRepo.TaskRun()
    if taskTagListRepo.LastError() != nil {
        return syncd.RenderTaskError(taskTagListRepo.Stdout() + taskTagListRepo.Stderr())
    }

    tagList := gostring.StrFilterSliceEmpty(strings.Split(taskTagListRepo.Stdout(), "\n"))
    tagList = gostring.StringSliceRsort(tagList)
    return syncd.RenderJson(c, goweb.JSON{
        "list": tagList,
    })
}

func ApplyRepoCommitList(c *goweb.Context) error {
    project, err := applyCheckAndGetProjectDetail(c.QueryInt("id"), c.GetInt("user_id"))
    if err != nil {
        return err
    }
    repo := &repoService.Repo{
        ID: project.ID,
        Url: project.RepoUrl,
    }
    if repo, err = repoService.RepoNew(repo); err != nil {
        return syncd.RenderAppError(err.Error())
    }

    updateRepoCmd, err := repo.UpdateRepo(project.RepoBranch)
    if err != nil {
        return syncd.RenderTaskError(err.Error())
    }
    taskUpdateRepo := taskService.TaskCreate(taskService.TASK_REPO_UPDATE, []string{
        updateRepoCmd,
    })

    taskCommitListRepo := taskService.TaskCreate(taskService.TASK_REPO_COMMIT_LIST, []string{
        repo.CommitListRepo(),
    })

    c.CloseCallback(func() {
        taskUpdateRepo.Terminate()
        taskCommitListRepo.Terminate()
    }, 60)

    taskUpdateRepo.TaskRun()
    if taskUpdateRepo.LastError() != nil {
        return syncd.RenderTaskError(taskUpdateRepo.Stdout() + taskUpdateRepo.Stderr())
    }

    taskCommitListRepo.TaskRun()
    if taskCommitListRepo.LastError() != nil {
        return syncd.RenderTaskError(taskCommitListRepo.Stdout() + taskCommitListRepo.Stderr())
    }
    commitList := gostring.StrFilterSliceEmpty(strings.Split(taskCommitListRepo.Stdout(), "\n"))
    return syncd.RenderJson(c, goweb.JSON{
        "list": commitList,
    })
}

func ApplySubmit(c *goweb.Context) error {
    params := ApplyParamValid{
        ProjectId: c.PostFormInt("project_id"),
        Name: c.PostForm("name"),
        Description: c.PostForm("description"),
    }
    if valid := govalidate.NewValidate(&params); !valid.Pass() {
        return syncd.RenderParamError(valid.LastFailed().Msg)
    }
    tag, commit := c.PostForm("tag"), c.PostForm("commit")

    project, err := applyCheckAndGetProjectDetail(c.PostFormInt("project_id"), c.GetInt("user_id"))
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if project.Status != 1 {
        return syncd.RenderParamError("project not enabled")
    }
    if project.RepoMode == 1 && commit == "" {
        return syncd.RenderParamError("commit can not be empty")
    }
    if project.RepoMode == 2 && tag == "" {
        return syncd.RenderParamError("tag can not be empty")
    }
    status := 1
    if project.NeedAudit == 0 {
        status = 3
    }

    apply := &deployService.Apply{
        ProjectId: project.ID,
        SpaceId: project.SpaceId,
        Name: params.Name,
        Description: params.Description,
        Status: status,
        UserId: c.GetInt("user_id"),
        RepoData: deployService.ApplyRepoData{
            RepoUrl: project.RepoUrl,
            RepoMode: project.RepoMode,
            RepoBranch: project.RepoBranch,
            Tag: tag,
            Commit: commit,
        },
    }
    pkId, err := apply.Create()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }

    logService.Record(&logService.OperateLog{
        DataId: pkId,
        OpType: logService.OP_TYPE_APPLY,
        OpName: logService.OP_NAME_APPLY_CREATE,
        UserId: c.GetInt("user_id"),
        UserName: c.GetString("user_name"),
    })

    return syncd.RenderJson(c, nil)
}

func ApplyUpdate(c *goweb.Context) error {
    id, tag, commit, name, description := c.PostFormInt("id"), c.PostForm("tag"), c.PostForm("commit"), c.PostForm("name"), c.PostForm("description")
    if id == 0 {
        return syncd.RenderParamError("id can not be empty")
    }
    if name == "" {
        return syncd.RenderParamError("name can not be empty")
    }
    if description == "" {
        return syncd.RenderParamError("description can not be empty")
    }

    apply, err := deployService.ApplyGetByPk(id)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if apply.UserId != c.GetInt("user_id") {
        return syncd.RenderAppError("no priv update apply")
    }

    project, err := applyCheckAndGetProjectDetail(apply.ProjectId, c.GetInt("user_id"))
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if project.Status != 1 {
        return syncd.RenderParamError("project not enabled")
    }
    if project.RepoMode == 1 && commit == "" {
        return syncd.RenderParamError("commit can not be empty")
    }
    if project.RepoMode == 2 && tag == "" {
        return syncd.RenderParamError("tag can not be empty")
    }
    if apply.Status == 2 {
        apply.Status = 1
    }
    apply.Name = name
    apply.Description = description
    apply.RepoData.Tag = tag
    apply.RepoData.Commit = commit

    if err := apply.Update(); err != nil {
        return syncd.RenderAppError(err.Error())
    }

    logService.Record(&logService.OperateLog{
        DataId: id,
        OpType: logService.OP_TYPE_APPLY,
        OpName: logService.OP_NAME_APPLY_UPDATE,
        UserId: c.GetInt("user_id"),
        UserName: c.GetString("user_name"),
    })

    return syncd.RenderJson(c, nil)
}

func ApplyList(c *goweb.Context) error {
    offset, limit, keyword, userId := c.QueryInt("offset"), c.GetInt("limit"), c.Query("keyword"), c.GetInt("user_id")
    searchTimeDuring, searchProjectId, searchStatus := c.QueryInt("time"), c.QueryInt("project_id"), c.QueryInt("status")

    apply := deployService.Apply{
        ProjectId: searchProjectId,
        Status: searchStatus,
    }
    if searchTimeDuring > 0 {
        apply.Ctime = int(time.Now().Unix()) - searchTimeDuring * 86400
    }
    if havePriv := userService.PrivIn(userService.DEPLOY_VIEW_ALL, c.GetIntSlice("priv")); !havePriv {
        apply.UserId = userId
    }
    userSpaceIds, err := projectService.SpaceGetIdListByUserId(userId)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    list, total, err := apply.List(keyword, userSpaceIds, offset, limit)

    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    var projectIds, spaceIds, userIds []int
    for _, l := range list {
        projectIds = append(projectIds, l.ProjectId)
        spaceIds = append(spaceIds, l.SpaceId)
        userIds = append(userIds, l.UserId)
    }
    projMaps, err := projectService.ProjectGetMapByIds(projectIds)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    spaceMaps, err := projectService.SpaceGetMapByIds(spaceIds)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    userMaps, err := userService.UserGetMapByIds(userIds)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }

    var newList []map[string]interface{}
    for _, l := range list {
        var (
            projectName, spaceName, userName, userEmail string
            userId int
        )
        if proj, exists := projMaps[l.ProjectId]; exists {
            projectName = proj.Name
        }
        if space, exists := spaceMaps[l.SpaceId]; exists {
            spaceName = space.Name
        }
        if user, exists := userMaps[l.UserId]; exists {
            userId = user.ID
            userName = user.Name
            userEmail = user.Email
        }
        newList = append(newList, map[string]interface{}{
            "id": l.ID,
            "name": l.Name,
            "project_name": projectName,
            "space_name": spaceName,
            "status": l.Status,
            "user_id": userId,
            "user_name": userName,
            "user_email": userEmail,
            "ctime": l.Ctime,
        })
    }

    return syncd.RenderJson(c, goweb.JSON{
        "list": newList,
        "total": total,
    })
}

func ApplyDetail(c *goweb.Context) error {
    apply := &deployService.Apply{
        ID: c.QueryInt("id"),
    }
    if err := apply.Detail(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if havePriv := userService.PrivIn(userService.DEPLOY_VIEW_ALL, c.GetIntSlice("priv")); !havePriv {
        if apply.UserId != c.GetInt("user_id") {
            return syncd.RenderCustomerError(syncd.CODE_ERR_NO_PRIV, "no priv")
        }
    }
    var spaceName, projectName, userName, userEmail string
    if space, err := projectService.SpaceGetByPk(apply.SpaceId); err == nil {
        spaceName = space.Name
    }
    if project, err := projectService.ProjectGetByPk(apply.ProjectId); err == nil {
        projectName = project.Name
    }
    if user, err := userService.UserGetByPk(apply.UserId); err == nil {
        userName, userEmail = user.Name, user.Email
    }

    return syncd.RenderJson(c, goweb.JSON{
        "id": apply.ID,
        "name": apply.Name,
        "description": apply.Description,
        "space_name": spaceName,
        "project_name": projectName,
        "project_id": apply.ProjectId,
        "user_name": userName,
        "user_email": userEmail,
        "repo_branch": apply.RepoData.RepoBranch,
        "repo_commit": apply.RepoData.Commit,
        "repo_tag": apply.RepoData.Tag,
        "repo_mode": apply.RepoData.RepoMode,
        "status": apply.Status,
        "ctime": apply.Ctime,
    })
}

func ApplyAudit(c *goweb.Context) error {
    audit, id, status := c.PostFormInt("audit"), c.PostFormInt("id"), 0
    var opContent string

    if audit == 1 {
        status = 3
        opContent = "审核通过"
    } else {
        status = 2
        opContent = "审核拒绝, "
    }
    apply, err := deployService.ApplyGetByPk(id)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    tmpApply := &deployService.Apply{
        ID: apply.ID,
        SpaceId: apply.SpaceId,
        UserId: apply.UserId,
        Status: apply.Status,
    }

    logService.Record(&logService.OperateLog{
        DataId: id,
        OpType: logService.OP_TYPE_APPLY,
        OpName: logService.OP_NAME_APPLY_AUDIT,
        OpContent: gostring.JoinStrings(opContent, c.PostForm("reject_reason")),
        UserId: c.GetInt("user_id"),
        UserName: c.GetString("user_name"),
    })

    return applyUpdateStatus(c, tmpApply, userService.DEPLOY_AUDIT_ALL, status, []int{1})
}

func ApplyUnAudit(c *goweb.Context) error {
    id := c.PostFormInt("id")
    apply, err := deployService.ApplyGetByPk(id)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    tmpApply := &deployService.Apply{
        ID: apply.ID,
        Status: apply.Status,
        UserId: apply.UserId,
        SpaceId: apply.SpaceId,
    }

    logService.Record(&logService.OperateLog{
        DataId: id,
        OpType: logService.OP_TYPE_APPLY,
        OpName: logService.OP_NAME_APPLY_UNAUDIT,
        UserId: c.GetInt("user_id"),
        UserName: c.GetString("user_name"),
    })

    return applyUpdateStatus(c, tmpApply, userService.DEPLOY_AUDIT_ALL, 1, []int{3})
}

func ApplyDiscard(c *goweb.Context) error {
    id := c.PostFormInt("id")
    apply, err := deployService.ApplyGetByPk(id)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    tmpApply := &deployService.Apply{
        ID: apply.ID,
        SpaceId: apply.SpaceId,
        UserId: apply.UserId,
        Status: apply.Status,
    }

    logService.Record(&logService.OperateLog{
        DataId: id,
        OpType: logService.OP_TYPE_APPLY,
        OpName: logService.OP_NAME_APPLY_DISCARD,
        UserId: c.GetInt("user_id"),
        UserName: c.GetString("user_name"),
    })

    return applyUpdateStatus(c, tmpApply, userService.DEPLOY_DROP_ALL, 7, []int{1, 2, 3, 6})
}

func ApplyLog(c *goweb.Context) error {
    id := c.QueryInt("id")
    apply, err := deployService.ApplyGetByPk(id)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if havePriv := userService.PrivIn(userService.DEPLOY_VIEW_ALL, c.GetIntSlice("priv")); !havePriv {
        if apply.UserId != c.GetInt("user_id") {
            return syncd.RenderCustomerError(syncd.CODE_ERR_NO_PRIV, "no priv")
        }
    }
    if err := applyCheckUserInSpace(apply.SpaceId, c.GetInt("user_id")); err != nil {
        return err
    }

    log := &logService.OperateLog{
        DataId: id,
        OpType: logService.OP_TYPE_APPLY,
    }
    list, err := log.List()
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, goweb.JSON{
        "list": list,
    })
}

func applyUpdateStatus(c *goweb.Context, apply *deployService.Apply, privCode int, status int, allowStatus []int) error {
    if havePriv := userService.PrivIn(privCode, c.GetIntSlice("priv")); !havePriv {
        if apply.UserId != c.GetInt("user_id") {
            return syncd.RenderCustomerError(syncd.CODE_ERR_NO_PRIV, "no priv")
        }
    }
    if err := applyCheckUserInSpace(apply.SpaceId, c.GetInt("user_id")); err != nil {
        return err
    }
    if !goslice.InSliceInt(apply.Status, allowStatus)  {
        return syncd.RenderAppError("apply status wrong")
    }
    apply.Status = status
    if err := apply.UpdateStatus(); err != nil {
        return syncd.RenderAppError(err.Error())
    }
    return syncd.RenderJson(c, nil)
}

func applyCheckAndGetProjectDetail(id, userId int) (*projectService.Project, error) {
    if id == 0 {
        return nil, syncd.RenderParamError("id can not empty")
    }
    project, err := projectService.ProjectGetByPk(id)
    if err != nil {
        return nil, syncd.RenderAppError(err.Error())
    }
    user := &projectService.User{
        UserId: userId,
        SpaceId: project.SpaceId,
    }
    if exists, err := user.CheckUserInSpace(); err != nil || !exists {
        return nil, syncd.RenderAppError("user have no privilege to access project")
    }
    return project, nil
}

func applyCheckUserInSpace(spaceId, userId int) error {
    user := &projectService.User{
        UserId: userId,
        SpaceId: spaceId,
    }
    if exists, err := user.CheckUserInSpace(); err != nil || !exists {
        return syncd.RenderAppError("user have no privilege to access project")
    }
    return nil
}

