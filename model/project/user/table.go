// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

type ProjectUser struct {
    ID          int      `gorm:"primary_key"`
    SpaceId     int      `gorm:"type:int(11);not null;default:0"`
    UserId      int      `gorm:"type:int(11);not null;default:0"`
    Ctime       int      `gorm:"type:int(11);not null;default:0"`
}

const (
    TableName = "project_user"
)
