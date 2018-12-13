// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package group

type Group struct {
    ID      int         `gorm:"primary_key" json:"id"`
    Name    string      `gorm:"type:varchar(100);not null;default:''" json:"name"`
    Ctime   int         `gorm:"type:int(11);not null;default:0" json:"ctime"`
    Mtime   int         `gorm:"type:int(11);not null;default:0" json:"mtime"`
}

const (
    TableName = "syd_server_group"
)
