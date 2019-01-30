// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    /*
    
    "fmt"

    
    "github.com/dreamans/syncd/util/gostring"
    */
    "errors"

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
    Status          int     `json:"status"`
}

func (a *Apply) Create() error {
    apply := &model.DeployApply{
        SpaceId: a.SpaceId,
        ProjectId: a.SpaceId,
        Name: a.Name,
        Description: a.Description,
        BranchName: a.BranchName,
        CommitVersion: a.CommitVersion,
        Status: a.SpaceId,
    }
    if ok := apply.Create(); !ok {
        return errors.New("create deploy apply failed")
    }
    return nil
}
