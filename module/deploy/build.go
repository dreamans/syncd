// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "errors"
    //"time"
    //"fmt"

    "github.com/dreamans/syncd/model"
)

type Build struct{
    ID          int     `json:"id"`
    ApplyId     int     `json:"apply_id"`
    StartTime   int     `json:"start_time"`
    FinishTime  int     `json:"finish_time"`
    Status      int     `json:"status"`
    Tar         string  `json:"tar"`
    Cmd         string  `json:"cmd"`
    Output      string  `json:"Output"`
    Ctime       int     `json:"ctime"`
}

const (
    BUILD_STATUS_START = 1
    BUILD_STATUS_SUCCESS = 2
    BUILD_STATUS_FAILED = 3
)

func (b *Build) Create() error {
    build := &model.DeployBuild{
        ApplyId: b.ApplyId,
        Status: b.Status,
        Cmd: b.Cmd,
    }
    if ok := build.Create(); !ok {
        return errors.New("create deploy build failed")
    }
    return nil
}

func (b *Build) Exists() (bool, error) {
    build := &model.DeployBuild{}
    if ok := build.GetByApplyId(b.ApplyId); !ok {
        return false, errors.New("get deploy build detail failed")
    }
    if build.ID == 0 {
        return false, nil
    }
    b.ID = build.ID
    b.Status = build.Status	
    return true, nil
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
