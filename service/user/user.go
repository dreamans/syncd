// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "errors"
    "fmt"

    "github.com/tinystack/goutil/gois"
    "github.com/tinystack/goutil/gostring"
    "github.com/dreamans/syncd/model"
    userModel "github.com/dreamans/syncd/model/user"
)

type User struct {
    ID          int     `json:"id"`
    GroupId     int     `json:"group_id"`
    Name        string  `json:"name"`
    Password    string  `json:"password"`
    Email       string  `json:"email"`
    TrueName    string  `json:"true_name"`
    Mobile      string  `json:"mobile"`
    LockStatus  int     `json:"lock_status"`
    Salt        string  `json:"salt"`
}

type UserItem struct {
    ID              int     `json:"id"`
    GroupId         int     `json:"group_id"`
    GroupName       string  `json:"group_name"`
    Name            string  `json:"name"`
    Email           string  `json:"email"`
    LockStatus      int     `json:"lock_status"`
    LastLoginTime   int     `json:"last_login_time"`
    LastLoginIp     string  `json:"last_login_ip"`
}

func UserGetByPk(id int) (*User, error) {
    user := &User{
        ID: id,
    }
    if err := user.Detail(); err != nil {
        return nil, err
    }
    return user, nil
}

func UserGetMapByIds(ids []int) (map[int]UserItem, error) {
    list, err := UserGetListByIds(ids)
    if err != nil {
        return nil, err
    }
    maps := map[int]UserItem{}
    for _, l := range list {
        maps[l.ID] = l
    }
    return maps, nil
}

func UserGetListByIds(ids []int) ([]UserItem, error){
    list, ok := userModel.List(model.QueryParam{
        Fields: "id, name, group_id, email, lock_status, last_login_ip, last_login_time",
        Order: "id DESC",
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "id",
                Tag: "IN",
                Prepare: ids,
            },
        },
    })
    if !ok {
        return nil, errors.New("get user list failed")
    }
    var userList []UserItem
    for _, l := range list {
        userList = append(userList, UserItem{
            ID: l.ID,
            GroupId: l.GroupId,
            Name: l.Name,
            Email: l.Email,
            LockStatus: l.LockStatus,
            LastLoginTime: l.LastLoginTime,
            LastLoginIp: l.LastLoginIp,
        })
    }
    return userList, nil
}

func (u *User) UpdatePassword() error {
    updateData := map[string]interface{}{
        "password": u.Password,
        "salt": u.Salt,
    }
    if ok := userModel.Update(u.ID, updateData); !ok {
        return errors.New("user password update failed")
    }
    return nil
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
        LockStatus: u.LockStatus,
    }
    var salt, password string
    if u.Password != "" {
        salt = gostring.StrRandom(10)
        password = gostring.StrMd5(gostring.JoinStrings(u.Password, salt))
    }
    if u.ID > 0 {
        updateData := map[string]interface{}{
            "group_id": user.GroupId,
            "name": user.Name,
            "email": user.Email,
            "true_name": user.TrueName,
            "mobile": user.Mobile,
            "lock_status": user.LockStatus,
        }
        if password != "" {
            updateData["password"] = password
            updateData["salt"] = salt
        }
        ok = userModel.Update(u.ID, updateData)
    } else {
        user.Password = password
        user.Salt = salt
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
        where []model.WhereParam
        userList []UserItem
    )
    if keyword != "" {
        var w *model.WhereParam
        switch {
        case gois.IsInteger(keyword):
            userId = gostring.Str2Int(keyword)
            if userId > 0 {
                w = &model.WhereParam{
                    Field: "id",
                    Prepare: userId,
                }
            }
        case gois.IsEmail(keyword):
            w = &model.WhereParam{
                Field: "email",
                Tag: "LIKE",
                Prepare: fmt.Sprintf("%%%s%%", keyword),
            }
        case gois.IsMobile(keyword):
            w = &model.WhereParam{
                Field: "mobile",
                Tag: "LIKE",
                Prepare: fmt.Sprintf("%%%s%%", keyword),
            }
        default:
            w = &model.WhereParam{
                Field: "name",
                Tag: "LIKE",
                Prepare: fmt.Sprintf("%%%s%%", keyword),
            }
        }
        if w != nil {
            where = append(where, *w)
        }
    }
    list, ok := userModel.List(model.QueryParam{
        Fields: "id, name, group_id, email, lock_status, last_login_ip, last_login_time",
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get user list data failed")
    }
    total, ok = userModel.Total(model.QueryParam{
        Where: where,
    })
    if !ok {
        return nil, 0, errors.New("get user total count failed")
    }
    for _, u := range list {
        userList = append(userList, UserItem{
            ID: u.ID,
            Name: u.Name,
            Email: u.Email,
            GroupId: u.GroupId,
            LockStatus: u.LockStatus,
            LastLoginIp: u.LastLoginIp,
            LastLoginTime: u.LastLoginTime,
        })
    }
    return userList, total, nil
}

func (u *User) Detail() error {
    if u.ID == 0 {
        return errors.New("id can not be empty")
    }
    detail, ok := userModel.Get(u.ID)
    if !ok {
        return errors.New("get user detail data failed")
    }
    if detail.ID == 0 {
        return errors.New("user not exists")
    }
    u.transmitUserDetail(detail)
    return nil
}

func (u *User) GetByName() error {
    if u.Name == "" {
        return errors.New("name can not be empty")
    }
    detail, ok := userModel.GetOne(model.QueryParam{
        Where: []model.WhereParam{
            model.WhereParam{
                Field: "name",
                Prepare: u.Name,
            },
        },
    })
    if !ok {
        return errors.New("get user detail data failed")
    }
    u.transmitUserDetail(detail)
    return nil
}

func (u *User) GetByNameOrEmail() error {
    if u.Name == "" && u.Email == "" {
        return errors.New("name or email can not be empty")
    }
    var where []model.WhereParam
    if u.Name != "" {
        where = append(where, model.WhereParam{
            Field: "name",
            Prepare: u.Name,
        })
    }
    if u.Email != "" {
        where = append(where, model.WhereParam{
            Field: "email",
            Prepare: u.Email,
        })
    }
    detail, ok := userModel.GetOne(model.QueryParam{
        Where: where,
    })
    if !ok {
        return errors.New("get user detail data failed")
    }
    u.transmitUserDetail(detail)
    return nil
}

func (u *User) CheckUserExists() (bool, error) {
    var where []model.WhereParam
    if u.Name != "" {
        where = append(where, model.WhereParam{
            Field: "name",
            Prepare: u.Name,
        })
    }
    if u.Email != "" {
        where = append(where, model.WhereParam{
            Field: "email",
            Prepare: u.Email,
        })
    }
    if u.ID > 0 {
        where = append(where, model.WhereParam{
            Field: "id",
            Tag: "!=",
            Prepare: u.ID,
        })
    }
    detail, ok := userModel.GetOne(model.QueryParam{
        Where: where,
    })
    if !ok {
        return false, errors.New("get user one data failed")
    }
    return detail.ID > 0, nil
}

func (u *User) Delete() error {
    if u.ID == 0 {
        return errors.New("id can not be empty")
    }
    if err := u.Detail(); err != nil {
        return err
    }
    ok := userModel.Delete(u.ID)
    if !ok {
        return errors.New("user delete failed")
    }
    return nil
}

func (u *User) Search() ([]UserItem, error){
    var where []model.WhereParam
    if u.Name != "" {
        where = append(where, model.WhereParam{
            Field: "name",
            Tag: "LIKE",
            Prepare: fmt.Sprintf("%%%s%%", u.Name),
        })
    }
    if u.Email != "" {
        where = append(where, model.WhereParam{
            Field: "email",
            Prepare: u.Email,
        })
    }
    list, ok := userModel.List(model.QueryParam{
        Fields: "id, name, group_id, email, lock_status",
        Order: "id DESC",
        Where: where,
    })
    if !ok {
        return nil, errors.New("get user list failed")
    }

    var userList []UserItem
    for _, l := range list {
        userList = append(userList, UserItem{
            ID: l.ID,
            GroupId: l.GroupId,
            Name: l.Name,
            Email: l.Email,
            LockStatus: l.LockStatus,
        })
    }
    return userList, nil
}

func (u *User) transmitUserDetail(detail userModel.User) {
    u.ID = detail.ID
    u.GroupId = detail.GroupId
    u.Name = detail.Name
    u.Password = detail.Password
    u.Email = detail.Email
    u.TrueName = detail.TrueName
    u.Mobile = detail.Mobile
    u.LockStatus = detail.LockStatus
    u.Salt = detail.Salt
}

func (u *User) DoNotUpdatePassword() {
    u.Password = ""
}