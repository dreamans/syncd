// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "errors"
    "fmt"

    "github.com/dreamans/syncd/util/gostring"
    "github.com/dreamans/syncd/util/gois"
    "github.com/dreamans/syncd/model"
)

type User struct {
    ID              int         `json:"id"`
    RoleId          int         `json:"role_id"`
    RoleName	    string	`json:"role_name"`
    Username	    string      `json:"username"`
    Password        string      `json:"password"`
    Email	    string      `json:"email"`
    Truename        string      `json:"truename"`
    Mobile          string      `json:"mobile"`
    Status          int         `json:"status"`
    LastLoginTime   int         `json:"last_login_time"`
    LastLoginIp     string      `json:"last_login_ip"`
    Ctime           int         `json:"ctime"`
}

func (u *User) Total(keyword string) (int, error) {
    user := &model.User{}
    total, ok := user.Count(model.QueryParam{
        Where: u.parseWhereConds(keyword),
    })
    if !ok {
        return 0, errors.New("get user list count failed")
    }
    return total, nil
}

func (u *User) List(keyword string, offset, limit int) ([]User, error) {
    user := &model.User{}
    list, ok := user.List(model.QueryParam{
        Fields: "id, role_id, username, email, status, last_login_time, last_login_ip",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: u.parseWhereConds(keyword),
    })
    if !ok {
        return nil, errors.New("get user list failed")
    }
    var roleIdList []int
    for _, l := range list {
        roleIdList = append(roleIdList, l.RoleId)
    }
    roleMap , err := RoleGetMapByIds(roleIdList)
    if err != nil {
        return nil, errors.New("get user map list failed")
    }

    var userList []User
    for _, l := range list {
        user := User{
            ID: l.ID,
            RoleId: l.RoleId,
            Username: l.Username,
            Email: l.Email,
            Status: l.Status,
            LastLoginTime: l.LastLoginTime,
            LastLoginIp: l.LastLoginIp,
        }
        if r, exists := roleMap[user.RoleId]; exists {
            user.RoleName = r.Name
        }
        userList = append(userList, user)
    }
    return userList, nil
}

func (u *User) UserExists() (bool, error) {
    var where []model.WhereParam
    if u.Username != "" {
        where = append(where, model.WhereParam{
            Field: "username",
            Prepare: u.Username,
        })
    }
    if u.ID != 0 {
        where = append(where, model.WhereParam{
            Field: "id",
            Tag: "!=",
            Prepare: u.ID,
        })
    }
    if u.Email != "" {
        where = append(where, model.WhereParam{
            Field: "email",
            Prepare: u.Email,
        })
    }
    user := &model.User{}
    count, ok := user.Count(model.QueryParam{
        Where: where,
    })
    if !ok {
        return false, errors.New("check user failed")
    }

    return count > 0, nil
}

func (u *User) CreateOrUpdate() error {
    var salt, password string
    if u.Password != "" {
        salt = gostring.StrRandom(10)
        password = gostring.StrMd5(gostring.JoinStrings(u.Password, salt))
    }
    user := &model.User{
        ID: u.ID,
        RoleId: u.RoleId,
        Username: u.Username,
        Email: u.Email,
        Truename: u.Truename,
        Mobile: u.Mobile,
        Status: u.Status,		
    }
    if u.ID > 0 {
        updateData := map[string]interface{}{
            "role_id": u.RoleId,
            "username": u.Username,
            "email": u.Email,
            "truename": u.Truename,
            "mobile": u.Mobile,
            "status": u.Status,
        }
        if password != "" {
            updateData["password"] = password
            updateData["salt"] = salt
        }
        ok := user.UpdateByFields(updateData, model.QueryParam{
            Where: []model.WhereParam{
                model.WhereParam{
                    Field: "id",
                    Prepare: u.ID,
                },
            },
        })
        if !ok {
            return errors.New("user create failed")
        }
    } else {
        user.Password = password
        user.Salt = salt
        ok := user.Create()
        if !ok {
            return errors.New("user update failed")
        }
    }
    return nil
}

func (u *User) parseWhereConds(keyword string) []model.WhereParam {
    var where []model.WhereParam
    if keyword != "" {
        if gois.IsInteger(keyword) {
            where = append(where, model.WhereParam{
                Field: "id",
                Prepare: keyword,
            })
        } else {
            if gois.IsEmail(keyword) {
                where = append(where, model.WhereParam{
                    Field: "email",
                    Prepare: keyword,
                })
            } else {
                where = append(where, model.WhereParam{
                    Field: "username",
                    Tag: "LIKE",
                    Prepare: fmt.Sprintf("%%%s%%", keyword),
                })
            }
        }
    }
    return where
}
