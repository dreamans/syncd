// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"github.com/gin-gonic/gin"
	"github.com/dreamans/syncd/router/user"
)

func RegisterRoute(g *gin.Engine) {
	api := g.Group("/api")
	{
		api.POST("/login", user.Login)
		api.GET("/login/status", user.LoginStatus)
	}
}