// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package group

type ServerGroup struct {
    ID      int         `gorm:"primary_key"`
    Name    string      `gorm:"type:varchar(100);not null;default:''"`
    Utime   int         `gorm:"type:int(11);not null;default:0"`
}

const (
    TableName = "server_group"
)
