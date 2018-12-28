// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    privModel "github.com/tinystack/syncd/model/user/priv"
)

func PrivList(c *goweb.Context) error {
    syncd.RenderJson(c, goweb.JSON{
        "list": privModel.PrivList,
    })
    return nil
}
