// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package render

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success", 
		"data": data,
	})
}
