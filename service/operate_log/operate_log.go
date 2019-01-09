// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package operate_log

import (
    "errors"

    baseModel "github.com/tinystack/syncd/model"
    operateLogModel "github.com/tinystack/syncd/model/operate_log"
)

type OperateLog struct {
    DataId      int     `json:"data_id"`
    OpType      string  `json:"op_type"`
    OpName      string  `json:"op_name"`
    OpContent   string  `json:"op_content"`
    UserId      int     `json:"user_id"`
    UserName    string  `json:"user_name"`
    Ctime       int     `json:"ctime"`
}

const (
    OP_TYPE_APPLY = "apply"
)

const (
    OP_NAME_APPLY_CREATE = "apply_create"
    OP_NAME_APPLY_UPDATE = "apply_update"
    OP_NAME_APPLY_AUDIT = "apply_audit"
    OP_NAME_APPLY_UNAUDIT = "apply_unaudit"
    OP_NAME_APPLY_DISCARD = "apply_discard"
    OP_NAME_APPLY_START = "apply_start"
    OP_NAME_APPLY_END = "apply_end"
    OP_NAME_APPLY_STOP = "apply_stop"
)

func Record(oplog *OperateLog) {
    oplog.Create()
}

func (op *OperateLog) Create() error {
    oper := &operateLogModel.OperateLog{
        DataId: op.DataId,
        OpType: op.OpType,
        OpName: op.OpName,
        OpContent: op.OpContent,
        UserId: op.UserId,
        UserName: op.UserName,
    }
    if ok := operateLogModel.Create(oper); !ok {
        return errors.New("operate log create failed")
    }
    return nil
}

func (op *OperateLog) List() ([]OperateLog, error) {
    var where []baseModel.WhereParam
    if op.DataId > 0 {
        where = append(where, baseModel.WhereParam{
            Field: "data_id",
            Prepare: op.DataId,
        })
    }
    if op.OpType != "" {
        where = append(where, baseModel.WhereParam{
            Field: "op_type",
            Prepare: op.OpType,
        })
    }
    list, ok := operateLogModel.List(baseModel.QueryParam{
        Fields: "id, data_id, op_type, op_name, op_content, user_id, user_name, ctime",
        Order: "id ASC",
        Where: where,
    })
    if !ok {
        return nil, errors.New("get operate log list failed")
    }
    var operList []OperateLog
    for _, l := range list {
        operList = append(operList, OperateLog{
            DataId: l.DataId,
            OpType: l.OpType,
            OpName: l.OpName,
            OpContent: l.OpContent,
            UserId: l.UserId,
            UserName: l.UserName,
            Ctime: l.Ctime,
        })
    }
    return operList, nil
}

