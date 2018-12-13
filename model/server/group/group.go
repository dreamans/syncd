// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package group

import (
    "time"

    "github.com/tinystack/syncd/model"
)

func Create(data *Group) bool {
    data.Ctime = int(time.Now().Unix())
    return model.Create(TableName, data)
}

func Update(id int, data Group) bool {
    data.Mtime = int(time.Now().Unix())
    ok := model.Update(TableName, data, model.QueryParam{
        Plain: "id = ?",
        Prepare: []interface{}{id},
    })
    return ok
}
