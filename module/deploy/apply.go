// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "errors"
    "fmt"

    "github.com/dreamans/syncd/model"
    "github.com/dreamans/syncd/util/gostring"
)

type Apply struct {
    ID          int     `json:"id"`
    SpaceId     int     `json:"space_id"`
    ProjectId   int     `json:"project_id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    BranchName  string  `json:"branch_name"`
    Status      int     `json:"status"`
}

func (a *Apply) Create() {

}
