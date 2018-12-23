// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "errors"

    baseModel "github.com/tinystack/syncd/model"
    projectUserModel "github.com/tinystack/syncd/model/project/user"
)

type User struct {
    ID          int     `json:"id"`
    UserId      int     `json:"user_id"`
    SpaceId     int     `json:"space_id"`
    Ctime       int     `json:"ctime"`
}

func (u *User) Add() error {
    user := projectUserModel.ProjectUser{
        SpaceId: u.SpaceId,
        UserId: u.UserId,
    }
    if ok := projectUserModel.Create(&user); !ok {
        return errors.New("project user add failed")
    }
    return nil
}

func (u *User) CheckUserExists() (bool, error) {
    where := []baseModel.WhereParam{
        baseModel.WhereParam{
            Field: "space_id",
            Prepare: u.SpaceId,
        },
        baseModel.WhereParam{
            Field: "user_id",
            Prepare: u.UserId,
        },
    }
    detail, ok := projectUserModel.GetOne(baseModel.QueryParam{
        Where: where,
    })
    if !ok {
        return false, errors.New("get project user data failed")
    }
    return detail.ID > 0, nil
}

func (u *User) List(offset, limit int) ([]User, int, error) {
    list, ok := projectUserModel.List(baseModel.QueryParam{
        Offset: offset,
        Limit: limit,
        Order: "id DESC",
    })
    if !ok {
        return nil, 0, errors.New("get project user list failed")
    }

    total, ok := projectUserModel.Total(baseModel.QueryParam{})
    if !ok {
        return nil, 0, errors.New("get project user total count failed")
    }

    var newList []User
    for _, l := range list {
        newList = append(newList, User{
            ID: l.ID,
            UserId: l.UserId,
            SpaceId: l.SpaceId,
            Ctime: l.Ctime,
        })
    }
    return newList, total, nil
}

func (u *User) Delete() error {
    if u.ID == 0 {
        return errors.New("id can not be empty")
    }
    ok := projectUserModel.Delete(u.ID)
    if !ok {
        return errors.New("project user remove failed")
    }
    return nil
}
