// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "errors"
    "time"
    "fmt"

    "github.com/dreamans/syncd/model"
)

type Apply struct {
    ID              int     `json:"id"`
    SpaceId         int     `json:"space_id"`
    ProjectId       int     `json:"project_id"`
    Name            string  `json:"name"`
    Description     string  `json:"description"`
    BranchName      string  `json:"branch_name"`
    CommitVersion   string  `json:"commit_version"`
    AuditStatus     int     `json:"audit_status"`
    Status          int     `json:"status"`
    UserId          int     `json:"user_id"`
    Ctime           int     `json:"ctime"`
}

const (
    AUDIT_STATUS_PENDING = 1
    AUDIT_STATUS_OK = 2
    AUDIT_STATUS_REFUSE = 3
)

const (
    STATUS_DEPLOY_NONE = 1
    STATUS_DEPLOY_ING = 2
    STATUS_DEPLOY_SUCCESS = 3
    STATUS_DEPLOY_FAILED = 4
    STATUS_DEPLOY_DROP = 5
)

func (a *Apply) Total(keyword string, spaceIds []int) (int, error) {
    apply := &model.DeployApply{}
    total, ok := apply.Count(model.QueryParam{
        Where: a.parseWhereConds(keyword, spaceIds),
    })
    if !ok {
        return 0, errors.New("get apply count failed")
    }
    return total, nil
}

func (a *Apply) List(keyword string, spaceIds []int, offset, limit int) ([]Apply, error) {
    apply := &model.DeployApply{}
    list, ok := apply.List(model.QueryParam{
        Fields: "id, space_id, project_id, name, user_id, audit_status, status, ctime",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: a.parseWhereConds(keyword, spaceIds),
    })
    if !ok {
        return nil, errors.New("get project list failed")
    }
    var applyList []Apply
    for _, l := range list {
        applyList = append(applyList, Apply{
            ID: l.ID,
            SpaceId: l.SpaceId,
            ProjectId: l.ProjectId,
            Name: l.Name,
            UserId: l.UserId,
            AuditStatus: l.AuditStatus,
            Status: l.Status,
            Ctime: l.Ctime,
        })
    }
    return applyList, nil
}

func (a *Apply) Create() error {
    apply := &model.DeployApply{
        SpaceId: a.SpaceId,
        ProjectId: a.ProjectId,
        Name: a.Name,
        Description: a.Description,
        BranchName: a.BranchName,
        CommitVersion: a.CommitVersion,
        Status: a.Status,
        UserId: a.UserId,
        AuditStatus: a.AuditStatus,
    }
    if ok := apply.Create(); !ok {
        return errors.New("create deploy apply failed")
    }
    return nil
}

func (a *Apply) parseWhereConds(keyword string, spaceIds []int) []model.WhereParam {
    var where []model.WhereParam
    where = append(where, model.WhereParam{
        Field: "space_id",
        Tag: "IN",
        Prepare:  spaceIds,
    })
    if a.Ctime != 0 {
        where = append(where, model.WhereParam{
            Field: "ctime",
            Tag: ">=",
            Prepare:  int(time.Now().Unix()) - a.Ctime * 86400,
        })
    }
    if a.AuditStatus != 0 {
        where = append(where, model.WhereParam{
            Field: "audit_status",
            Prepare:  a.AuditStatus,
        })
    }
    if a.Status != 0 {
        where = append(where, model.WhereParam{
            Field: "status",
            Prepare:  a.Status,
        })
    }
    if a.ProjectId != 0 {
        where = append(where, model.WhereParam{
            Field: "project_id",
            Prepare:  a.ProjectId,
        })
    }
    if keyword != "" {
        where = append(where, model.WhereParam{
            Field: "name",
            Tag: "LIKE",
            Prepare: fmt.Sprintf("%%%s%%", keyword),
        })
    }
    return where
}
