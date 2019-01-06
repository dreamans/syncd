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
    ID          int
    ApplyId     int
    Level       int
    Cmd         string
    Status      int
}

const (
    DEPLOY_LEVEL_UPDATE_REPO = 1
    DEPLOY_LEVEL_UPDATE_COMMIT = 2
    DEPLOY_LEVEL_PACK_REPO = 3
    DEPLOY_LEVEL_DEPLOY = 4
)

func (t *DeployTask) Create() (int, error) {
    dt := &deployTaskModel.DeployTask{
        ApplyId: t.ApplyId,
        Level: t.Level,
        Cmd: t.Cmd,
        Status: 1,
    }
    if ok := deployTaskModel.Create(dt); !ok {
        return 0, errors.New("deploy task save to db failed")
    }
    return dt.ID, nil
}

func (t *DeployTask) Flush() {
    if t.ApplyId == 0 {
        return errors.New("apply_id can not empty")
    }
    if ok := deployTaskModel.DeleteByApplyId(); !ok {
        return 0, errors.New("deploy task save to db failed")
    }
}

