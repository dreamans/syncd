// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/render"
)

func Login(c *gin.Context) {
    render.JSON(c, gin.H{"id": 1,})
}

func LoginStatus(c *gin.Context) {
    render.JSON(c, gin.H{"id": 1,})
}