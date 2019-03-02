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
    ServerId    int     `json:"server_id"`
    Status      int     `json:"status"`
    Output      string  `json:"output"`
    Errmsg      string  `json:"errmsg"`
    StartTime   int     `json:"start_time"`
    FinishTime  int     `json:"finish_time"`
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
        Fields: "id, apply_id, group_id, server_id, status, output, errmsg, ctime",
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
            ServerId: l.ServerId,
            Status: l.Status,
            Output: l.Output,
            Errmsg: l.Errmsg,
            StartTime: l.StartTime,
            FinishTime: l.FinishTime,
            Ctime: l.Ctime,
        })
    }

    return deployList, nil
}

func (d *Deploy) Create() error {
    deploy := model.DeployTask{
        ApplyId: d.ApplyId,
        GroupId: d.GroupId,
        ServerId: d.ServerId,
        Status: d.Status,
    }
    if ok := deploy.Create(); !ok {
        return errors.New("create deploy task failed")
    }
    return nil
}

func (d *Deploy) UpdateStart() error {
    dt := &model.DeployTask{}
    updateData := map[string]interface{}{
        "status": d.Status,
        "start_time": d.StartTime,
    }
    if ok := dt.UpdateByFields(updateData, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "apply_id",
                Prepare: d.ApplyId,
            },
            model.WhereParam{
                Field: "server_id",
                Prepare: d.ServerId,
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
        "output": d.Output,
        "errmsg": d.Errmsg,
        "finish_time": d.FinishTime,
    }
    if ok := dt.UpdateByFields(updateData, model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "apply_id",
                Prepare: d.ApplyId,
            },
            model.WhereParam{
                Field: "server_id",
                Prepare: d.ServerId,
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
