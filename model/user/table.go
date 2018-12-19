// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

type User struct {
    ID              int     `gorm:"primary_key"`
    GroupId         int     `gorm:"type:int(11);not null;default:0"`
    Name            string  `gorm:"type:varchar(100);not null;default:''"`
    Password        string  `gorm:"type:char(32);not null;default:''"`
    Email           string  `gorm:"type:varchar(100);not null;default:''"`
    TrueName        string  `gorm:"type:varchar(20);not null;default:''"`
    Mobile          string  `gorm:"type:varchar(20);not null;default:''"`
    Salt            string  `gorm:"type:varchar(10);not null;default:''"`
    LockStatus      int     `gorm:"type:int(11);not null;default:1"`
    LastLoginIp     string  `gorm:"type:varchar(45);not null;default:''"`
    LastLoginTime   int     `gorm:"type:int(11);not null;default:0"`
    Ctime           int     `gorm:"type:int(11);not null;default:0"`
}

const (
    TableName = "user"
)
