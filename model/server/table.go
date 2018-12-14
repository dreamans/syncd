// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

type Server struct {
    ID      int         `gorm:"primary_key" json:"id"`
    GroupId int         `gorm:"type:int(11);not null;default:0" json:"group_id"`
    Name    string      `gorm:"type:varchar(100);not null;default:''" json:"name"`
    Ip      string      `gorm:"type:varchar(15);not null;default:''" json:"ip"`
    SshPort int         `gorm:"type:int(11);not null;default:22" json:"ssh_port"`
    Utime   int         `gorm:"type:int(11);not null;default:0" json:"utime"`
}

const (
    TableName = "syd_server"
)
