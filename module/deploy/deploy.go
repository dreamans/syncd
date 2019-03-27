// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "errors"

    "github.com/dreamans/syncd/model"
)

type Deploy struct{
    ID          int     `json:"id"`
    ApplyId     int     `json:"apply_id"`
    GroupId     int     `json:"group_id"`
    Status      int     `json:"status"`
    Content     string  `json:"content"`
    Ctime       int     `json:"ctime"`
}

const (
    DEPLOY_STATUS_NONE = 0
    DEPLOY_STATUS_START = 1
    DEPLOY_STATUS_SUCCESS = 2
    DEPLOY_STATUS_FAILED = 3
)

func (d *Deploy) TaskList() ([]Deploy, error) {
    dt := &model.DeployTask{}
    list, ok := dt.List(model.QueryParam{
        Fields: "id, apply_id, group_id, status, content, ctime",
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "apply_id",
                Prepare: d.ApplyId,
            },
        },
    })
    if !ok {
        return nil, errors.New("get deploy task list failed")
    }
    var (
        deployList []Deploy
    )
    for _, l := range list {
        deployList = append(deployList, Deploy{
            ID: l.ID,
            ApplyId: l.ApplyId,
            GroupId: l.GroupId,
            Status: l.Status,
            Content: l.Content,
            Ctime: l.Ctime,
        })
    }
    return deployList, nil
}

func (d *Deploy) Create() error {
    deploy := model.DeployTask{
        ApplyId: d.ApplyId,
        GroupId: d.GroupId,
        Status: d.Status,
    }
    if ok := deploy.Create(); !ok {
        return errors.New("create deploy task failed")
    }
    return nil
}

func (d *Deploy) UpdateStatus() error {
    dt := &model.DeployTask{}
    updateData := map[string]interface{}{
        "status": d.Status,
    }
    if ok := dt.UpdateByFields(updateData, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "apply_id",
                Prepare: d.ApplyId,
            },
            model.WhereParam{
                Field: "group_id",
                Prepare: d.GroupId,
            },
        },
    }); !ok {
        return errors.New("update deploy task result failed")
    }
    return nil
}

func (d *Deploy) UpdateResult() error {
    dt := &model.DeployTask{}
    updateData := map[string]interface{}{
        "status": d.Status,
        "content": d.Content,
    }
    if ok := dt.UpdateByFields(updateData, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "apply_id",
                Prepare: d.ApplyId,
            },
            model.WhereParam{
                Field: "group_id",
                Prepare: d.GroupId,
            },
        },
    }); !ok {
        return errors.New("update deploy task result failed")
    }
    return nil
}

func (d *Deploy) DeleteByApplyId() error {
    dep := &model.DeployTask{
        ApplyId: d.ApplyId,
    }
    if ok := dep.Delete(model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "apply_id",
                Prepare: d.ApplyId,
            },
        },
    }); !ok {
        return errors.New("remove deploy task failed")
    }
    return nil
}