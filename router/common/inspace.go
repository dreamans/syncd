// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package common

import (
	"github.com/gin-gonic/gin"
	"github.com/dreamans/syncd/render"
	"github.com/dreamans/syncd/module/project"
)

func InSpaceCheck(c *gin.Context, spaceId int) bool {
	member := &project.Member{
        UserId: c.GetInt("user_id"),
        SpaceId: spaceId,
    }
    if in := member.MemberInSpace(); !in {
		render.CustomerError(c, render.CODE_ERR_NO_PRIV, "user is not in the project space")
		return false
	}
	return true
}