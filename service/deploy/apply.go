// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "errors"
    "fmt"

    "github.com/tinystack/goutil/gois"
    "github.com/tinystack/goutil/gostring"
    baseModel "github.com/tinystack/syncd/model"
    deployApplyModel "github.com/tinystack/syncd/model/deploy_apply"
)

type Apply struct {
    ID              int             `json:"id"`
    ProjectId       int             `json:"project_id"`
    SpaceId         int             `json:"space_id"`
    Name            string          `json:"name"`
    Description     string          `json:"description"`
    RepoData        ApplyRepoData   `json:"repo_data"`
    Status          int             `json:"status"`
    UserId          int             `json:"user_id"`
    Ctime           int             `json:"ctime"`
}

type ApplyRepoData struct {
    Repo        string      `json:"repo"`
    RepoUrl     string      `json:"repo_url"`
    RepoUser    string      `json:"repo_user"`
    RepoPass    string      `json:"repo_pass"`
    RepoMode    int         `json:"repo_mode"`
    RepoBranch  string      `json:"repo_branch"`
    Tag         string      `json:"repo_tag"`
    Commit      string      `json:"repo_commit"`
}

func ApplyGetByPk(id int) (*Apply, error) {
    apply := &Apply{
        ID: id,
    }
    if err := apply.Detail(); err != nil {
        return nil, err
    }
    return apply, nil
}

func (a *Apply) List(keyword string, spaceIds []int, offset, limit int) ([]Apply, int, error) {
    var where []baseModel.WhereParam
    if keyword != "" {
        if gois.IsInteger(keyword) {
            applyId := gostring.Str2Int(keyword)
            if applyId > 0 {
                where = append(where, baseModel.WhereParam{
                    Field: "id",
                    Prepare: applyId,
                })
            }
        } else {
            where = append(where, baseModel.WhereParam{
                Field: "name",
                Tag: "LIKE",
                Prepare: fmt.Sprintf("%%%s%%", keyword),
            })
        }
    }
    if a.UserId > 0 {
        where = append(where, baseModel.WhereParam{
            Field: "user_id",
            Prepare: a.UserId,
        })
    }
    if a.ProjectId > 0 {
        where = append(where, baseModel.WhereParam{
            Field: "project_id",
            Prepare: a.ProjectId,
        })
    }
    if a.Ctime > 0 {
        where = append(where, baseModel.WhereParam{
            Field: "ctime",
            Tag: ">",
            Prepare: a.Ctime,
        })
    }
    if a.Status > 0 {
        where = append(where, baseModel.WhereParam{
            Field: "status",
            Prepare: a.Status,
        })
    }
    where = append(where, baseModel.WhereParam{
        Field: "space_id",
        Tag: "IN",
        Prepare: spaceIds,
    })
    list, ok := deployApplyModel.List(baseModel.QueryParam{
        Fields: "id, project_id, space_id, name, status, user_id, ctime",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get apply list failed")
    }
    total, ok := deployApplyModel.Total(baseModel.QueryParam{
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get apply total count failed")
    }
    var nlist []Apply
    for _, l := range list {
        nlist = append(nlist, Apply{
            ID: l.ID,
            ProjectId: l.ProjectId,
            SpaceId: l.SpaceId,
            Name: l.Name,
            Status: l.Status,
            UserId: l.UserId,
            Ctime: l.Ctime,
        })
    }
    return nlist, total, nil
}

func (a *Apply) Create() error {
    repoData, err := gostring.JsonEncode(a.RepoData)
    if err != nil {
        return err
    }
    apply := deployApplyModel.DeployApply{
        ProjectId: a.ProjectId,
        Name: a.Name,
        Description: a.Description,
        SpaceId: a.SpaceId,
        RepoData: repoData,
        Status: a.Status,
        UserId: a.UserId,
    }
    if ok := deployApplyModel.Create(&apply); !ok {
        return errors.New("apply submit failed")
    }
    return nil
}

func (a *Apply) UpdateStatus() error {
    ok := deployApplyModel.Update(a.ID, map[string]interface{}{
        "status": a.Status,
    })
    if !ok {
        return errors.New("update apply status failed")
    }
    return nil
}

func (a *Apply) Detail() error {
    if a.ID == 0 {
        return errors.New("id can not empty")
    }
    detail, ok := deployApplyModel.Get(a.ID)
    if !ok {
        return errors.New("apply detail get failed")
    }
    if detail.ID == 0 {
        return errors.New("apply not exists")
    }
    a.ID = detail.ID
    a.ProjectId = detail.ProjectId
    a.SpaceId = detail.SpaceId
    a.Name = detail.Name
    a.Description = detail.Description
    a.Status = detail.Status
    a.UserId = detail.UserId
    a.Ctime = detail.Ctime
    gostring.JsonDecode(detail.RepoData, &a.RepoData)

    return nil
}

