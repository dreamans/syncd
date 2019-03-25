// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "errors"
    "fmt"

    "github.com/dreamans/syncd/model"
    "github.com/dreamans/syncd/util/gostring"
)

type Project struct {
    ID                  int     `json:"id"`
    SpaceId             int     `json:"space_id"`
    Name                string  `json:"name"`
    Description         string  `json:"description"`
    NeedAudit           int     `json:"need_audit"`
    Status              int     `json:"status"`
    RepoUrl             string  `json:"repo_url"`
    RepoBranch          string  `json:"repo_branch"`
    DeployMode          int     `json:"deploy_mode"`
    OnlineCluster       []int   `json:"online_cluster"`
    DeployUser          string  `json:"deploy_user"`
    DeployPath          string  `json:"deploy_path"`
    BuildScript         string  `json:"build_script"`
    BuildHookScript     string  `json:"build_hook_script"`
    DeployHookScript    string  `json:"deploy_hook_script"`
    PreDeployCmd        string  `json:"pre_deploy_cmd"`
    AfterDeployCmd      string  `json:"after_deploy_cmd"`
    AuditNotice         string  `json:"audit_notice"`
    DeployNotice        string  `json:"deploy_notice"`
    Ctime               int     `json:"ctime"`
}

func ProjectListByIds(projectIds []int) ([]Project, error) {
    project := &model.Project{}
    list, ok := project.List(model.QueryParam{
        Fields: "id, name",
        Order: "id DESC",
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Tag: "IN",
                Prepare: projectIds,
            },
        },
    })
    if !ok {
        return nil, errors.New("get project list failed")
    }

    var projList []Project
    for _, l := range list {
        projList = append(projList, Project{
            ID: l.ID,
            Name: l.Name,
        })
    }
    return projList, nil
}

func ProjectAllBySpaceIds(spaceIds []int) ([]Project, error) {
    project := &model.Project{}
    list, ok := project.List(model.QueryParam{
        Fields: "id, name, space_id",
        Order: "id DESC",
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "space_id",
                Tag: "IN",
                Prepare: spaceIds,
            },
        },
    })
    if !ok {
        return nil, errors.New("get project list failed")
    }

    var projList []Project
    for _, l := range list {
        projList = append(projList, Project{
            ID: l.ID,
            Name: l.Name,
            SpaceId: l.SpaceId,
        })
    }
    return projList, nil
}

func (p *Project) UpdateBuildScript() error {
    project := &model.Project{}
    updateData := map[string]interface{}{
        "build_script": p.BuildScript,
    }
    if ok := project.UpdateByFields(updateData, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Prepare: p.ID,
            },
        },
    }); !ok {
        return errors.New("project build_script update failed")
    }
    return nil
}

func (p *Project) UpdateHookScript() error {
    project := &model.Project{}
    updateData := map[string]interface{}{
        "build_hook_script": p.BuildHookScript,
        "deploy_hook_script": p.DeployHookScript,
    }
    if ok := project.UpdateByFields(updateData, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Prepare: p.ID,
            },
        },
    }); !ok {
        return errors.New("project hook script update failed")
    }
    return nil
}

func (p *Project) Detail() error {
    project := &model.Project{}
    if ok := project.Get(p.ID); !ok {
        return errors.New("get project detail failed")
    }
    if project.ID == 0 {
        return errors.New("project detail not exists")
    }

    p.ID = project.ID
    p.SpaceId = project.SpaceId
    p.Name = project.Name
    p.Description = project.Description
    p.NeedAudit = project.NeedAudit
    p.Status = project.Status
    p.RepoUrl = project.RepoUrl
    p.DeployMode = project.DeployMode
    p.RepoBranch = project.RepoBranch
    p.OnlineCluster = gostring.StrSplit2IntSlice(project.OnlineCluster, ",")
    p.DeployUser = project.DeployUser
    p.DeployPath = project.DeployPath
    p.PreDeployCmd = project.PreDeployCmd
    p.AfterDeployCmd = project.AfterDeployCmd
    p.Ctime = project.Ctime
    p.BuildScript = project.BuildScript
    p.BuildHookScript = project.BuildHookScript
    p.DeployHookScript = project.DeployHookScript
    p.AuditNotice = project.AuditNotice
    p.DeployNotice = project.DeployNotice

    return nil
}

func (p *Project) UpdateStatus() error {
    project := &model.Project{}
    updateData := map[string]interface{}{
        "status": p.Status,
    }
    ok := project.UpdateByFields(updateData, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Prepare: p.ID,
            },
        },
    })
    if !ok {
        return errors.New("project status update failed")
    }
    return nil
}

func (p *Project) Total(keyword string, spaceId int) (int, error) {
    project := &model.Project{}
    total, ok := project.Count(model.QueryParam{
        Where: p.parseWhereConds(keyword, spaceId),
    })
    if !ok {
        return 0, errors.New("get project count failed")
    }
    return total, nil
}

func (p *Project) List(keyword string, spaceId, offset, limit int) ([]Project, error) {
    project := &model.Project{}
    list, ok := project.List(model.QueryParam{
        Fields: "id, name, need_audit, status",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: p.parseWhereConds(keyword, spaceId),
    })
    if !ok {
        return nil, errors.New("get project list failed")
    }

    var projList []Project
    for _, l := range list {
        projList = append(projList, Project{
            ID: l.ID,
            Name: l.Name,
            NeedAudit: l.NeedAudit,
            Status: l.Status,
        })
    }
    return projList, nil
}

func (p *Project) CreateOrUpdate() error {
    project := &model.Project{
        ID: p.ID,
        SpaceId: p.SpaceId,
        Name: p.Name,
        Description: p.Description,
        NeedAudit: p.NeedAudit,
        RepoUrl: p.RepoUrl,
        DeployMode: p.DeployMode,
        RepoBranch: p.RepoBranch,
        OnlineCluster: gostring.JoinIntSlice2String(p.OnlineCluster, ","),
        DeployUser: p.DeployUser,
        DeployPath: p.DeployPath,
        PreDeployCmd: p.PreDeployCmd,
        AfterDeployCmd: p.AfterDeployCmd,
        AuditNotice: p.AuditNotice,
        DeployNotice: p.DeployNotice,
    }
    if project.ID > 0 {
        updateData := map[string]interface{}{
            "name": p.Name,
            "description": p.Description,
            "need_audit": p.NeedAudit,
            "repo_url": p.RepoUrl,
            "deploy_mode": p.DeployMode,
            "repo_branch": p.RepoBranch,
            "online_cluster": gostring.JoinIntSlice2String(p.OnlineCluster, ","),
            "deploy_user": p.DeployUser,
            "deploy_path": p.DeployPath,
            "pre_deploy_cmd": p.PreDeployCmd,
            "after_deploy_cmd": p.AfterDeployCmd,
            "audit_notice": p.AuditNotice,
            "deploy_notice": p.DeployNotice,
        }
        if ok := project.UpdateByFields(updateData, model.QueryParam{
            Where: []model.WhereParam{
                model.WhereParam{
                    Field: "id",
                    Prepare: project.ID,
                },
            },
        }); !ok {
            return errors.New("project update failed")
        }
    } else {
        if ok := project.Create(); !ok {
            return errors.New("project create failed")
        }
    }
    return nil
}

func (p *Project) Delete() error {
    project := &model.Project{
        ID: p.ID,
    }
    if ok := project.Delete(); !ok {
        return errors.New("project detail failed")
    }
    return nil
}

func (p *Project) parseWhereConds(keyword string, spaceId int) []model.WhereParam {
    var where []model.WhereParam
    where = append(where, model.WhereParam{
        Field: "space_id",
        Prepare: spaceId,
    })
    if keyword != "" {
        where = append(where, model.WhereParam{
            Field: "name",
            Tag: "LIKE",
            Prepare: fmt.Sprintf("%%%s%%", keyword),
        })
    }
    return where
}
