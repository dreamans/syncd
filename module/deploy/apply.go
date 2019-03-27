// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "errors"
    "time"
    "fmt"

    "github.com/dreamans/syncd/model"
    "github.com/dreamans/syncd/util/gois"
)

type Apply struct {
    ID                  int     `json:"id"`
    SpaceId             int     `json:"space_id"`
    ProjectId           int     `json:"project_id"`
    Name                string  `json:"name"`
    Description         string  `json:"description"`
    BranchName          string  `json:"branch_name"`
    CommitVersion       string  `json:"commit_version"`
    AuditStatus         int     `json:"audit_status"`
    AuditRefusalReasion string  `json:"audit_refusal_reasion"`
    RollbackId          int     `json:"rollback_id"`
    RollbackApplyId     int     `json:"rollback_apply_id"`
    IsRollbackApply     int     `json:"is_rollback_apply"`
    Status              int     `json:"status"`
    UserId              int     `json:"user_id"`
    Username            string  `json:"username"`
    Email               string  `json:"email"`
    RollbackStatus      int     `json:"rollback_status"`
    Ctime               int     `json:"ctime"`
}

const (
    AUDIT_STATUS_PENDING = 1
    AUDIT_STATUS_OK = 2
    AUDIT_STATUS_REFUSE = 3
)

const (
    APPLY_STATUS_NONE = 1
    APPLY_STATUS_ING = 2
    APPLY_STATUS_SUCCESS = 3
    APPLY_STATUS_FAILED = 4
    APPLY_STATUS_DROP = 5
    APPLY_STATUS_ROLLBACK = 6
)

func (a *Apply) RollbackList() ([]Apply, error) {
    apply := &model.DeployApply{}
    list, ok := apply.List(model.QueryParam{
        Fields: "id, name",
        Limit: 10,
        Order: "ctime DESC",
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "project_id",
                Prepare: a.ProjectId,
            },
            model.WhereParam{
                Field: "audit_status",
                Prepare: AUDIT_STATUS_OK,
            },
            model.WhereParam{
                Field: "status",
                Prepare: APPLY_STATUS_SUCCESS,
            },
            model.WhereParam{
                Field: "is_rollback_apply",
                Prepare: 0,
            },
        },
    })
    if !ok {
        return nil, errors.New("get apply list failed")
    }
    var applyList []Apply
    for _, l := range list {
        applyList = append(applyList, Apply{
            ID: l.ID,
            Name: l.Name,
        })
    }
    return applyList, nil
}

func (a *Apply) DropStatus() error {
    apply := &model.DeployApply{}
    updateData := map[string]interface{}{
        "status": APPLY_STATUS_DROP,
    }
    if ok := apply.UpdateByFields(updateData, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Prepare: a.ID,
            },
        },
    }); !ok {
        return errors.New("update deploy apply status failed")
    }

    return nil
}

func (a *Apply) Update() error {
    apply := &model.DeployApply{}
    updateData := map[string]interface{}{
        "branch_name": a.BranchName,
        "audit_status": a.AuditStatus,
        "commit_version": a.CommitVersion,
        "description": a.Description,
    }
    if ok := apply.UpdateByFields(updateData, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Prepare: a.ID,
            },
        },
    }); !ok {
        return errors.New("update deploy apply failed")
    }

    return nil
}

func (a *Apply) UpdateStatus() error {
    apply := &model.DeployApply{}
    updateData := map[string]interface{}{
        "status": a.Status,
    }
    if ok := apply.UpdateByFields(updateData, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Prepare: a.ID,
            },
        },
    }); !ok {
        return errors.New("update deploy apply status failed")
    }

    return nil
}

func (a *Apply) UpdateRollback() error {
    apply := &model.DeployApply{}
    updateData := map[string]interface{}{
        "rollback_apply_id": a.RollbackApplyId,
        "status": APPLY_STATUS_ROLLBACK,
    }
    if ok := apply.UpdateByFields(updateData, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Prepare: a.ID,
            },
        },
    }); !ok {
        return errors.New("update deploy apply rollback_apply_id failed")
    }

    return nil
}

func (a *Apply) UpdateAuditStatus() error {
    apply := &model.DeployApply{}
    updateData := map[string]interface{}{
        "audit_status": a.AuditStatus,
        "audit_refusal_reasion": a.AuditRefusalReasion,
    }
    if ok := apply.UpdateByFields(updateData, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Prepare: a.ID,
            },
        },
    }); !ok {
        return errors.New("update apply audit_status failed")
    }

    return nil
}

func (a *Apply) Detail() error {
    apply := &model.DeployApply{}
    if ok := apply.Get(a.ID); !ok {
        return errors.New("get deploy apply detail failed")
    }
    if apply.ID == 0 {
        return errors.New("deploy apply detail not exists")
    }
    a.SpaceId = apply.SpaceId
    a.ProjectId = apply.ProjectId
    a.Name = apply.Name
    a.Description = apply.Description
    a.BranchName = apply.BranchName
    a.CommitVersion = apply.CommitVersion
    a.AuditStatus = apply.AuditStatus
    a.Status = apply.Status
    a.UserId = apply.UserId
    a.RollbackId = apply.RollbackId
    a.RollbackApplyId = apply.RollbackApplyId
    a.IsRollbackApply = apply.IsRollbackApply
    a.RollbackStatus = APPLY_STATUS_NONE
    a.Ctime = apply.Ctime
    return nil
}

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

func (a *Apply) CheckHaveDeploying() (bool, error) {
    apply := &model.DeployApply{}
    count, ok := apply.Count(model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Tag: "!=",
                Prepare: a.ID,
            },
            model.WhereParam{
                Field: "project_id",
                Prepare: a.ProjectId,
            },
            model.WhereParam{
                Field: "status",
                Prepare: APPLY_STATUS_ING,
            },
            model.WhereParam{
                Field: "ctime",
                Tag: ">=",
                Prepare: int(time.Now().Unix()) - 86400,
            },
        },
    })
    if !ok {
        return false, errors.New("get apply count failed")
    }

    return count == 0, nil
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
        return nil, errors.New("get apply list failed")
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
        RollbackId: a.RollbackId,
        IsRollbackApply: a.IsRollbackApply,
        RollbackApplyId: a.RollbackApplyId,
    }
    if ok := apply.Create(); !ok {
        return errors.New("create deploy apply failed")
    }
    a.ID = apply.ID
    return nil
}

func (a *Apply) parseWhereConds(keyword string, spaceIds []int) []model.WhereParam {
    var where []model.WhereParam
    where = append(where, model.WhereParam{
        Field: "space_id",
        Tag: "IN",
        Prepare:  spaceIds,
    })
    where = append(where, model.WhereParam{
        Field: "is_rollback_apply",
        Prepare:  a.IsRollbackApply,
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
        if gois.IsInteger(keyword) {
            where = append(where, model.WhereParam{
                Field: "id",
                Prepare: keyword,
            })
        } else {
            where = append(where, model.WhereParam{
                Field: "name",
                Tag: "LIKE",
                Prepare: fmt.Sprintf("%%%s%%", keyword),
            })
        }
    }
    return where
}
