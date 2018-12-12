// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

import (
    "github.com/tinystack/syncd"
)

func Create(data interface{}) bool {
    if db := syncd.Orm.Create(data); db.Error != nil {
        syncd.Logger.Warning("mysql query error: %v", db.Error)
        return false
    }
    return true
}
