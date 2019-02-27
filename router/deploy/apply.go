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
    "github.com/dreamans/syncd/module/user"
    "github.com/dreamans/syncd/module/server"
)

type ApplyFormBind struct {
    ProjectId       int     `form:"project_id" binding:"required"`
    SpaceId         int     `form:"space_id" binding:"required"`
    Name            string  `form:"name" binding:"required"`
    BranchName      string  `form:"branch_name"`
    CommitVersion   string  `form:"commit_version"`
    Description     string  `form:"description"`
}

type ApplyQueryBind struct {
    Time        int     `form:"time"`
    AuditStatus int     `form:"audit_status"`
    Status      int     `form:"status"`
    ProjectId   int     `form:"project_id"`
    Keyword     string  `form:"keyword"`
    Offset      int     `form:"offset"`
    Limit       int     `form:"limit" binding:"required,gte=1,lte=999"`
}

type ApplyAuditFormBind struct {
    ID          int     `form:"id" binding:"required"`
    AuditStatus int     `form:"audit_status" binding:"required"`
    AuditRefusalReasion string  `form:"audit_refusal_reasion"`
}

type ApplyUpdateFormBind struct {
    ID              int     `form:"id" binding:"required"`
    BranchName      string  `form:"branch_name" binding:"required"`
    CommitVersion   string  `form:"commit_version"`
    Description     string  `form:"description" binding:"required"`
}

func ApplyDrop(c *gin.Context) {
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
    if deploy.STATUS_DEPLOY_ING == apply.Status {
        render.AppError(c, "deploy apply status incorrect")
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
    if err := apply.DropStatus(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    render.JSON(c, nil)
}

func ApplyUpdate(c *gin.Context) {
    var form ApplyUpdateFormBind
    if err := c.ShouldBind(&form); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    apply := &deploy.Apply{
        ID: form.ID,
    }
    if err := apply.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    if deploy.STATUS_DEPLOY_NONE != apply.Status || deploy.AUDIT_STATUS_OK == apply.AuditStatus {
        render.AppError(c, "deploy apply status incorrect")
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
    branchName := apply.BranchName
    if form.BranchName != "" {
        branchName = form.BranchName
    } 
    apply = &deploy.Apply{
        ID: form.ID,
        BranchName: branchName,
        AuditStatus: deploy.AUDIT_STATUS_PENDING,
        CommitVersion: form.CommitVersion,
        Description: form.Description,
    }

    if err := apply.Update(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    render.JSON(c, nil)
}

func ApplyAudit(c *gin.Context) {
    var form ApplyAuditFormBind
    if err := c.ShouldBind(&form); err != nil {
        render.ParamError(c, err.Error())
        return
    }

    apply := &deploy.Apply{
        ID: form.ID,
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

    apply = &deploy.Apply{
        ID: form.ID,
        AuditStatus: form.AuditStatus,
        AuditRefusalReasion: form.AuditRefusalReasion,
    }
    if err := apply.UpdateAuditStatus(); err != nil {
        render.AppError(c, err.Error())
        return
    }

    render.JSON(c, nil)
}

func ApplyDetail(c *gin.Context) {
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

    u := &user.User{
        ID: apply.UserId,
    }
    if err := u.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    apply.Username = u.Username
    apply.Email = u.Email

    render.JSON(c, apply)
}

func ApplyList(c *gin.Context) {
    var query ApplyQueryBind
    if err := c.ShouldBind(&query); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    m := &project.Member{
        UserId: c.GetInt("user_id"),
    }
    spaceIds, err := m.SpaceIdsByUserId()
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    apply := &deploy.Apply{
        Ctime: query.Time,
        AuditStatus: query.AuditStatus,
        Status: query.Status,
        ProjectId: query.ProjectId,
    }
    list, err := apply.List(query.Keyword, spaceIds, query.Offset, query.Limit)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    var (
        projectIds, userIds []int
    )

    spaceMap := map[int]project.Space{}
    projectMap := map[int]project.Project{}
    userMap := map[int]user.User{}

    for _, l := range list {
        projectIds = append(projectIds, l.ProjectId)
        userIds = append(userIds, l.UserId)
    }
    spaceList, err := project.SpaceListByIds(spaceIds)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    for _, l := range spaceList {
        spaceMap[l.ID] = l
    }

    projectList, err := project.ProjectListByIds(projectIds)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    for _, l := range projectList {
        projectMap[l.ID] = l
    }

    userList, err := user.UserGetListByIds(userIds)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    for _, l := range userList {
        userMap[l.ID] = l
    }

    restList := []map[string]interface{}{}
    for _, l := range list {
        var spaceName, projectName, userName, email string
        if space, exists := spaceMap[l.SpaceId]; exists {
            spaceName = space.Name
        }
        if proj, exists := projectMap[l.ProjectId]; exists {
            projectName = proj.Name
        }
        if u, exists := userMap[l.UserId]; exists {
            userName = u.Username
            email = u.Email
        }
        restList = append(restList, map[string]interface{}{
            "id": l.ID,
            "name": l.Name,
            "space_id": l.SpaceId,
            "space_name": spaceName,
            "project_id": l.ProjectId,
            "project_name": projectName,
            "ctime": l.Ctime,
            "username": userName,
            "email": email,
            "audit_status": l.AuditStatus,
            "status": l.Status,
        })
    }

    total, err := apply.Total(query.Keyword, spaceIds)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    render.JSON(c, gin.H{
        "list": restList,
        "total": total,
    })
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

    m := &project.Member{
        UserId: c.GetInt("user_id"),
        SpaceId: form.SpaceId,
    }
    if in := m.MemberInSpace(); !in {
        render.CustomerError(c, render.CODE_ERR_NO_PRIV, "user is not in the project space")
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
        UserId: c.GetInt("user_id"),
        AuditStatus: deploy.AUDIT_STATUS_PENDING,
        Status: deploy.STATUS_DEPLOY_NONE,
    }

    if proj.NeedAudit == 0 {
        apply.AuditStatus = deploy.AUDIT_STATUS_OK
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

    m := &project.Member{
        UserId: c.GetInt("user_id"),
        SpaceId: proj.SpaceId,
    }
    if in := m.MemberInSpace(); !in {
        render.CustomerError(c, render.CODE_ERR_NO_PRIV, "user is not in the project space")
        return
    }

    clusterIds := proj.OnlineCluster
    if proj.PreReleaseCluster > 0 {
        clusterIds = append(clusterIds, proj.PreReleaseCluster)
    }
    clusterList, err := server.GroupGetMapByIds(clusterIds)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }
    serverList, err := server.ServerGetListByGroupIds(clusterIds)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    restProj := map[string]interface{}{
        "id": proj.ID,
        "name": proj.Name,
        "deploy_mode": proj.DeployMode,
        "repo_branch": proj.RepoBranch,
        "cluster_list": clusterList,
        "pre_cluster_id": proj.PreReleaseCluster,
        "online_cluster_ids": proj.OnlineCluster,
        "server_list": serverList,
    }

    render.JSON(c, restProj)
}
