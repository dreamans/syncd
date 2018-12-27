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
    deployApplyModel "github.com/tinystack/syncd/model/deploy/apply"
)

type Apply struct {
    ID              int             `json:"id"`
    ProjectId       int             `json:"project_id"`
    SpaceId         int             `json:"space_id"`
    Name            string          `json:"name"`
    Description     string          `json:"description"`
    RepoData        ApplyRepoData   `json:"repo_data"`
    Status          int             `json:"status"`
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

func (a *Apply) List(keyword string, offset, limit int) ([]Apply, int, error) {
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
    list, ok := deployApplyModel.List(baseModel.QueryParam{
        Fields: "id, project_id, space_id, name, status, ctime",
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
    }
    if ok := deployApplyModel.Create(&apply); !ok {
        return errors.New("apply submit failed")
    }
    return nil
}
