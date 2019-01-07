// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "errors"

    baseModel "github.com/tinystack/syncd/model"
    deployTaskModel "github.com/tinystack/syncd/model/deploy_task"
)

type DeployTask struct {
    ID          int     `json:"id"`
    ApplyId     int     `json:"apply_id"`
    Level       int     `json:"level"`
    Cmd         string  `json:"cmd"`
    Status      int     `json:"status"`
    Name        string  `json:"name"`
    Output      string  `json:"output"`
}

const (
    DEPLOY_LEVEL_UPDATE_REPO = 1
    DEPLOY_LEVEL_UPDATE_COMMIT = 2
    DEPLOY_LEVEL_PACK_REPO = 3
    DEPLOY_LEVEL_DEPLOY = 4
)

const (
    DEPLOY_STATUS_INIT = 1
    DEPLOY_STATUS_START = 2
    DEPLOY_STATUS_END = 3
    DEPOLY_STATUS_PANIC = 4
    DEPLOY_STATUS_STOP = 5
)

func (t *DeployTask) Create() (int, error) {
    dt := &deployTaskModel.DeployTask{
        ApplyId: t.ApplyId,
        Level: t.Level,
        Cmd: t.Cmd,
        Status: DEPLOY_STATUS_INIT,
        Name: t.Name,
    }
    if ok := deployTaskModel.Create(dt); !ok {
        return 0, errors.New("deploy task save to db failed")
    }
    return dt.ID, nil
}

func (t *DeployTask) UpdateStatus() error {
    updateData := map[string]interface{}{
        "status": t.Status,
    }
    if t.Output != "" {
        updateData["output"] = t.Output
    }
    ok := deployTaskModel.UpdateByPk(t.ID, updateData)
    if !ok {
        return errors.New("update deploy task status failed")
    }
    return nil
}

func (t *DeployTask) GetTaskItem() ([]DeployTask, error) {
    var where []baseModel.WhereParam
    where = append(where, baseModel.WhereParam{
        Field: "apply_id",
        Prepare: t.ApplyId,
    })
    if t.Status > 0 {
        where = append(where, baseModel.WhereParam{
            Field: "status",
            Prepare: t.Status,
        })
    }
    list, ok := deployTaskModel.List(baseModel.QueryParam{
        Where: where,
        Fields: "id, name, apply_id, level, cmd, status, output",
        Order: "level ASC",
    })
    if !ok {
        return nil, errors.New("get task item list from db failed")
    }

    var taskList []DeployTask
    for _, l := range list {
        taskList = append(taskList, DeployTask{
            ID: l.ID,
            ApplyId: l.ApplyId,
            Level: l.Level,
            Cmd: l.Cmd,
            Status: l.Status,
            Name: l.Name,
            Output: l.Output,
        })
    }
    return taskList, nil
}

func (t *DeployTask) Flush() error {
    if t.ApplyId == 0 {
        return errors.New("apply_id can not empty")
    }
    if ok := deployTaskModel.DeleteByApplyId(t.ApplyId); !ok {
        return errors.New("flush task from db failed")
    }
    return nil
}

