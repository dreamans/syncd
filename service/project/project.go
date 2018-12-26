// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "errors"
    "fmt"

    "github.com/tinystack/goutil"
    baseModel "github.com/tinystack/syncd/model"
    projectModel "github.com/tinystack/syncd/model/project"
)

type Project struct {
    ID              int     `json:"id"`
    Name            string  `json:"name"`
    Description     string  `json:"description"`
    SpaceId         int     `json:"space_id"`
    Repo            string  `json:"repo"`
    RepoUrl         string  `json:"repo_url"`
    DeployServer    []int   `json:"deploy_server"`
    DeployUser      string  `json:"deploy_user"`
    DeployPath      string  `json:"deploy_path"`
    DeployHistory   int     `json:"deploy_history"`
    PreDeployCmd    string  `json:"pre_deploy_cmd"`
    PostDeployCmd   string  `json:"post_deploy_cmd"`
    NeedAudit       int     `json:"need_audit"`
    Status          int     `json:"status"`
    RepoUser        string  `json:"repo_user"`
    RepoPass        string  `json:"repo_pass"`
    RepoMode        int     `json:"repo_mode"`
    RepoBranch      string  `json:"repo_branch"`
}

type ProjectItem struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    RepoMode    int     `json:"repo_mode"`
    NeedAudit   int     `json:"need_audit"`
    Status      int     `json:"status"`
}

func ProjectGetByPk(id int) (*Project, error) {
    project := &Project{
        ID: id,
    }
    if err := project.Get(); err != nil {
        return nil, err
    }
    return project, nil
}

func ProjectGetMapByIds(ids []int) (map[int]ProjectItem, error) {
    list, err := ProjectGetListByIds(ids)
    if err != nil {
        return nil, err
    }
    maps := map[int]ProjectItem{}
    for _, l := range list {
        maps[l.ID] = l
    }
    return maps, nil
}

func ProjectGetListByIds(ids []int) ([]ProjectItem, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    list, ok := projectModel.List(baseModel.QueryParam{
        Fields: "id, name, repo_mode, need_audit, status",
        Order: "id DESC",
        Where: []baseModel.WhereParam{
            baseModel.WhereParam{
                Field: "id",
                Tag: "IN",
                Prepare: ids,
            },
        },
    })
    if !ok {
        return nil, errors.New("get project list failed")
    }
    var projList []ProjectItem
    for _, l := range list {
        projList = append(projList, ProjectItem{
            ID: l.ID,
            Name: l.Name,
            RepoMode: l.RepoMode,
            NeedAudit: l.NeedAudit,
            Status: l.Status,
        })
    }
    return projList, nil
}

func (p *Project) List(keyword string, offset, limit int) ([]ProjectItem, int, error) {
    var (
        projectId int
        where []baseModel.WhereParam
    )
    if p.SpaceId > 0 {
        where = append(where, baseModel.WhereParam{
            Field: "space_id",
            Prepare: p.SpaceId,
        })
    }
    if p.Status == 1 {
        where = append(where, baseModel.WhereParam{
            Field: "status",
            Prepare: p.Status,
        })
    }
    if keyword != "" {
        if goutil.IsInteger(keyword) {
            projectId = goutil.Str2Int(keyword)
            if projectId > 0 {
                where = append(where, baseModel.WhereParam{
                    Field: "id",
                    Prepare: projectId,
                })
            }
        } else {
            if goutil.IsIp(keyword) {
                where = append(where, baseModel.WhereParam{
                    Field: "ip",
                    Prepare: keyword,
                })
            } else {
                where = append(where, baseModel.WhereParam{
                    Field: "name",
                    Tag: "LIKE",
                    Prepare: fmt.Sprintf("%%%s%%", keyword),
                })
            }
        }
    }

    list, ok := projectModel.List(baseModel.QueryParam{
        Fields: "id, name, repo_mode, need_audit, status",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get project list failed")
    }

    total, ok := projectModel.Total(baseModel.QueryParam{
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get project total count failed")
    }

    var nlist []ProjectItem
    for _, l := range list {
        nlist = append(nlist, ProjectItem{
            ID: l.ID,
            Name: l.Name,
            RepoMode: l.RepoMode,
            NeedAudit: l.NeedAudit,
            Status: l.Status,
        })
    }
    return nlist, total, nil
}

func (p *Project) Get() error {
    if p.ID == 0 {
        return errors.New("id can not be empty")
    }
    detail, ok := projectModel.Get(p.ID)
    if !ok {
        return errors.New("get project detail data failed")
    }
    p.Name = detail.Name
    p.Description = detail.Description
    p.SpaceId = detail.SpaceId
    p.Repo = detail.Repo
    p.RepoUrl = detail.RepoUrl
    p.DeployServer = goutil.StrSplit2IntSlice(detail.DeployServer, ",")
    p.DeployUser = detail.DeployUser
    p.DeployPath = detail.DeployPath
    p.DeployHistory = detail.DeployHistory
    p.PreDeployCmd = detail.PreDeployCmd
    p.PostDeployCmd = detail.PostDeployCmd
    p.NeedAudit = detail.NeedAudit
    p.Status = detail.Status
    p.RepoUser = detail.RepoUser
    p.RepoPass = detail.RepoPass
    p.RepoMode = detail.RepoMode
    p.RepoBranch = detail.RepoBranch

    return nil
}

func (p *Project) CreateOrUpdate() error {
    project := projectModel.Project{
        Name: p.Name,
        Description: p.Description,
        SpaceId: p.SpaceId,
        Repo: p.Repo,
        RepoUrl: p.RepoUrl,
        DeployServer: goutil.JoinIntSlice2String(p.DeployServer, ","),
        DeployUser: p.DeployUser,
        DeployPath: p.DeployPath,
        DeployHistory: p.DeployHistory,
        PreDeployCmd: p.PreDeployCmd,
        PostDeployCmd: p.PostDeployCmd,
        NeedAudit: p.NeedAudit,
        RepoUser: p.RepoUser,
        RepoPass: p.RepoPass,
        RepoMode: p.RepoMode,
        RepoBranch: p.RepoBranch,
    }
    if p.ID > 0 {
        if ok := projectModel.Update(p.ID, project); !ok {
            return errors.New("project data update failed")
        }
    } else {
        if ok := projectModel.Create(&project); !ok {
            return errors.New("project data create failed")
        }
    }
    return nil
}

func (p *Project) Delete() error {
    if p.ID == 0 {
        return errors.New("id can not be empty")
    }
    ok := projectModel.Delete(p.ID)
    if !ok {
        return errors.New("project delete failed")
    }
    return nil
}

func (p *Project) CheckSpaceHaveProject() (bool, error) {
    where := []baseModel.WhereParam{
        baseModel.WhereParam{
            Field: "space_id",
            Prepare: p.SpaceId,
        },
    }
    detail, ok := projectModel.GetOne(baseModel.QueryParam{
        Where: where,
    })
    if !ok {
        return false, errors.New("project get failed")
    }
    return detail.ID > 0, nil
}

func (p *Project) CheckProjectExists() (bool, error) {
    var where []baseModel.WhereParam
    where = append(where, baseModel.WhereParam{
        Field: "name",
        Prepare: p.Name,
    })
    where = append(where, baseModel.WhereParam{
        Field: "space_id",
        Prepare: p.SpaceId,
    })
    if p.ID > 0 {
        where = append(where, baseModel.WhereParam{
            Field: "id",
            Tag: "!=",
            Prepare: p.ID,
        })
    }
    detail, ok := projectModel.GetOne(baseModel.QueryParam{
        Where: where,
    })
    if !ok {
        return false, errors.New("get project one data failed")
    }
    return detail.ID > 0, nil
}

func (p *Project) ChangeStatus() error {
    if p.ID == 0 {
        return errors.New("id can not be empty")
    }
    ok := projectModel.UpdateFields(p.ID, map[string]interface{}{
        "status": p.Status,
    })
    if !ok {
        return errors.New("project status update failed")
    }
    return nil
}
