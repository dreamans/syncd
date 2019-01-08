// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package syncd

import (
    "net/http"
    "errors"
    "fmt"

    "github.com/tinystack/goweb"
)

const (
    CODE_OK = 0
    CODE_ERR_SYSTEM = 1000
    CODE_ERR_APP = 1001
    CODE_ERR_PARAM = 1002
    CODE_ERR_DATA_REPEAT = 1003
    CODE_ERR_LOGIN_FAILED = 1004
    CODE_ERR_NO_LOGIN = 1005
    CODE_ERR_NO_PRIV = 1006
    CODE_ERR_TASK_ERROR = 1007
    CODE_ERR_USER_OR_PASS_WRONG = 1008
)

func RenderParamError(msg string) error {
    return RenderCustomerError(CODE_ERR_PARAM, msg)
}

func RenderAppError(msg string) error {
    return RenderCustomerError(CODE_ERR_APP, msg)
}

func RenderTaskError(msg string) error {
    return RenderCustomerError(CODE_ERR_TASK_ERROR, msg)
}

func RenderCustomerError(code int, msg string) error {
    return errors.New(fmt.Sprintf("%d=>%s", code, msg))
}

func RenderJson(c *goweb.Context, data interface{}) error {
    c.Json(http.StatusOK, goweb.JSON{
        "code": CODE_OK,
        "message": "success",
        "data": data,
    })
    return nil
}

