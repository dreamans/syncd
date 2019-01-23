// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package render

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

const (
    CODE_SUCCESS = 0
    CODE_PARAM_ERROR = 10001
    CODE_APP_ERROR = 10002
)

func JSON(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, gin.H{
        "code": CODE_SUCCESS,
        "message": "success",
        "data": data,
    })
}

func ParamError(c *gin.Context, message string) {
    c.JSON(http.StatusOK, gin.H{
        "code": CODE_PARAM_ERROR,
        "message": message,
    })
}

func AppError(c *gin.Context, message string) {
    c.JSON(http.StatusOK, gin.H{
        "code": CODE_APP_ERROR,
        "message": message,
    })
}

func Success(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "code": CODE_SUCCESS,
        "message": "success",
    })
}
