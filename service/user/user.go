// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "errors"
    "fmt"

    "github.com/tinystack/goutil"
    baseModel "github.com/tinystack/syncd/model"
    userModel "github.com/tinystack/syncd/model/user"
)

type User struct {
    ID          int     `json:"id"`
    GroupId     int     `json:"group_id"`
    Name        string  `json:"name"`
    Password    string  `json:"password"`
    Email       string  `json:"Email"`
    TrueName    string  `json:"true_name"`
    Mobile      string  `json:"mobile"`
}

type UserItem struct {
    ID              int     `json:"id"`
    GroupId         int     `json:"group_id"`
    Name            string  `json:"name"`
    LastLoginTime   int     `json:"last_login_time"`
    LastLoginIp     string  `json:"last_login_ip"`
}

func (u *User) CreateOrUpdate() error {
    var ok bool
    user := userModel.User{
        ID: u.ID,
        GroupId: u.GroupId,
        Name: u.Name,
        Email: u.Email,
        TrueName: u.TrueName,
        Mobile: u.Mobile,
    }
    if u.ID > 0 {
        ok = userModel.Update(u.ID, map[string]interface{}{
            "group_id": user.GroupId,
            "name": user.Name,
            "email": user.Email,
            "true_name": user.TrueName,
            "mobile": user.Mobile,
        })
    } else {
        ok = userModel.Create(&user)
    }
    if !ok {
        return errors.New("user data update failed")
    }
    return nil
}

func (u *User) List(keyword string, offset, limit int) ([]UserItem, int, error) {
    var (
        ok bool
        userId, total int
        where []baseModel.WhereParam
        userList []UserItem
    )
    if keyword != "" {
        var w *baseModel.WhereParam
        switch {
        case goutil.IsInteger(keyword):
            userId = goutil.Str2Int(keyword)
            if userId > 0 {
                w = &baseModel.WhereParam{
                    Field: "id",
                    Prepare: userId,
                }
            }
        case goutil.IsEmail(keyword):
            w = &baseModel.WhereParam{
                Field: "email",
                Tag: "LIKE",
                Prepare: fmt.Sprintf("%%%s%%", keyword),
            }
        case goutil.IsMobile(keyword):
            w = &baseModel.WhereParam{
                Field: "mobile",
                Tag: "LIKE",
                Prepare: fmt.Sprintf("%%%s%%", keyword),
            }
        default:
            w = &baseModel.WhereParam{
                Field: "name",
                Tag: "LIKE",
                Prepare: fmt.Sprintf("%%%s%%", keyword),
            }
        }
        if w != nil {
            where = append(where, *w)
        }
    }
    list, ok := userModel.List(baseModel.QueryParam{
        Fields: "id, name",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get user group list data failed")
    }
    total, ok = userModel.Total(baseModel.QueryParam{
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get user group total count failed")
    }

    for _, u := range list {
        userList = append(userList, UserItem{
            ID: u.ID,
            Name: u.Name,
            GroupId: u.GroupId,
            LastLoginIp: u.LastLoginIp,
            LastLoginTime: u.LastLoginTime,
        })
    }
    return userList, total, nil
}
