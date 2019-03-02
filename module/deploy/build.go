// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "errors"
    "time"

    "github.com/dreamans/syncd/model"
)

type Build struct{
    ID          int     `json:"id"`
    ApplyId     int     `json:"apply_id"`
    StartTime   int     `json:"start_time"`
    FinishTime  int     `json:"finish_time"`
    Status      int     `json:"status"`
    Tar         string  `json:"tar"`
    Output      string  `json:"Output"`
    Errmsg      string  `json:"errmsg"`
    Ctime       int     `json:"ctime"`
}

const (
    BUILD_STATUS_NONE = 0
    BUILD_STATUS_START = 1
    BUILD_STATUS_SUCCESS = 2
    BUILD_STATUS_FAILED = 3
)

func (b *Build) Create() error {
    build := &model.DeployBuild{
        ApplyId: b.ApplyId,
        Status: b.Status,
        StartTime: int(time.Now().Unix()),
    }
    if ok := build.Create(); !ok {
        return errors.New("create deploy build failed")
    }
    return nil
}

func (b *Build) CreateFull() error {
    build := &model.DeployBuild{
        ApplyId: b.ApplyId,
        StartTime: b.StartTime,
        FinishTime: b.FinishTime,
        Status: b.Status,
        Tar: b.Tar,
        Output: b.Output,
        Errmsg: b.Errmsg,
    }
    if ok := build.Create(); !ok {
        return errors.New("create deploy build failed")
    }
    return nil
}

func (b *Build) Detail() error {
    build := &model.DeployBuild{}
    if ok := build.GetByApplyId(b.ApplyId); !ok {
        return errors.New("get deploy build detail failed")
    }
    if build.ID == 0 {
        build.Status = BUILD_STATUS_NONE
        return nil
    }
    b.ID = build.ID
    b.Status = build.Status
    b.StartTime = build.StartTime
    b.FinishTime = build.FinishTime
    b.Tar = build.Tar
    b.Output = build.Output
    b.Errmsg = build.Errmsg
    b.Ctime = build.Ctime

    return nil
}

func (b *Build) Exists() (bool, error) {
    if err := b.Detail(); err != nil {
        return false, err
    }
    if b.ID == 0 {
        return false, nil
    }
    return true, nil 
}

func (b *Build) Finish() error {
    build := &model.DeployBuild{}
    updateData := map[string]interface{}{
        "status": b.Status,
        "tar": b.Tar,
        "output": b.Output,
        "finish_time": int(time.Now().Unix()),
        "errmsg": b.Errmsg,
    }
    if ok := build.UpdateByFields(updateData, model.QueryParam{
        Where:[]model.WhereParam{
            model.WhereParam{
                Field: "apply_id",
                Prepare: b.ApplyId,
            },
        },
    }); !ok {
        return errors.New("update deploy build failed")
    }
    return nil
}

func (b *Build) Delete() error {
    build := &model.DeployBuild{
        ID: b.ID,
    }
    if ok := build.Delete(); !ok {
        return errors.New("remove deploy build failed")
    }
    return nil
}
