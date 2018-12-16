// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package priv

type PrivItem struct {
    Label   string  `json:"label"`
    Value   int     `json:"value"`
}

type PrivGroup struct {
    Label   string      `json:"label"`
    Items   []PrivItem  `json:"items"`
}

const (
    PROJECT_VIEW = 101
    PROJECT_NEW = 102
    PROJECT_EDIT = 103
    PROJECT_DEL = 104

    SERVER_GROUP_VIEW = 201
    SERVER_GROUP_NEW = 202
    SERVER_GROUP_EDIT = 203
    SERVER_GROUP_DEL = 204
    SERVER_VIEW = 205
    SERVER_NEW = 206
    SERVER_EDIT = 207
    SERVER_DEL = 208
)

var projectPriv = PrivGroup{
    Label: "项目",
    Items: []PrivItem{
        PrivItem{
            Label: "查看项目",
            Value: PROJECT_VIEW,
        },
        PrivItem{
            Label: "新建项目",
            Value: PROJECT_NEW,
        },
        PrivItem{
            Label: "编辑项目",
            Value: PROJECT_EDIT,
        },
        PrivItem{
            Label: "删除项目",
            Value: PROJECT_DEL,
        },
    },
}

var serverPriv = PrivGroup{
    Label: "服务器",
    Items: []PrivItem{
        PrivItem{
            Label: "查看分组",
            Value: SERVER_GROUP_VIEW,
        },
        PrivItem{
            Label: "新建分组",
            Value: SERVER_GROUP_NEW,
        },
        PrivItem{
            Label: "编辑分组",
            Value: SERVER_GROUP_EDIT,
        },
        PrivItem{
            Label: "删除分组",
            Value: SERVER_GROUP_DEL,
        },
        PrivItem{
            Label: "查看服务器",
            Value: SERVER_VIEW,
        },
        PrivItem{
            Label: "新建服务器",
            Value: SERVER_NEW,
        },
        PrivItem{
            Label: "编辑服务器",
            Value: SERVER_EDIT,
        },
        PrivItem{
            Label: "删除服务器",
            Value: SERVER_DEL,
        },
    },
}

var PrivList = []PrivGroup{
    projectPriv,
    serverPriv,
}
