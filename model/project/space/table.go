// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package space

type ProjectSpace struct {
    ID          int         `gorm:"primary_key"`
    Name        string      `gorm:"type:varchar(100);unique;not null;default:''"`
    Description string      `gorm:"type:varchar(500);not null;default:''"`
    Ctime       int         `gorm:"type:int(11);not null;default:0"`
}

const (
    TableName = "project_space"
)
