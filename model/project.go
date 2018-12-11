// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

import (
    //"github.com/jinzhu/gorm"
    "github.com/tinystack/syncd"
)

type Project struct {
    ID          int     `gorm:"primary_key"`
    Name        string  `gorm:"type:varchar(100);not null;default:''"`
    Description string  `gorm:"type:varchar(100);not null;default:''"`
    Space       string  `gorm:"type:varchar(100);not null;default:''"`
    BuildScript string
}

func (p Project) Create() {
    syncd.Orm.Create(&p)
}
