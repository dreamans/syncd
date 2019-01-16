// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "errors"
    "time"

    baseModel "github.com/dreamans/syncd/model"
    userTokenModel "github.com/dreamans/syncd/model/user_token"
)

type Token struct {
    ID              int     `json:"id"`
    UserId          int     `json:"user_id"`
    Token           string  `json:"token"`
    ExpireTime      int     `json:"expire_time"`
    Ctime           int     `json:"ctime"`
}

func (t *Token) CreateOrUpdate() error {
    detail, ok := userTokenModel.GetOne(baseModel.QueryParam{
        Where: []baseModel.WhereParam{
            baseModel.WhereParam{
                Field: "user_id",
                Prepare: t.UserId,
            },
        },
    })
    if !ok {
        return errors.New("get user token detail failed")
    }
    if detail.ID == 0 {
        ok := userTokenModel.Create(&userTokenModel.UserToken{
            UserId: t.UserId,
            Token: t.Token,
            ExpireTime: t.ExpireTime,
        })
        if !ok {
            return errors.New("user token create failed")
        }
    } else {
        ok := userTokenModel.Update(detail.ID, map[string]interface{}{
            "token": t.Token,
            "expire_time": t.ExpireTime,
        })
        if !ok {
            return errors.New("user token update failed")
        }
    }
    return nil
}

func (t *Token) ValidateToken() bool {
    if t.UserId == 0 || t.Token == "" {
        return false
    }
    detail, ok := userTokenModel.GetOne(baseModel.QueryParam{
        Where: []baseModel.WhereParam{
            baseModel.WhereParam{
                Field: "user_id",
                Prepare: t.UserId,
            },
            baseModel.WhereParam{
                Field: "token",
                Prepare: t.Token,
            },
        },
    })
    if !ok {
        return false
    }
    if detail.ID == 0 {
        return false
    }
    if detail.ExpireTime < int(time.Now().Unix()) {
        return false
    }
    t.ID = detail.ID
    return true
}

func (t *Token) DeleteByUserId() error {
    if t.UserId == 0 {
        return errors.New("user_id can not be empty")
    }
    if ok := userTokenModel.DeleteByUserId(t.UserId); !ok {
        return errors.New("token delete failed")
    }

    return nil
}

func (t *Token) UpdateExpirationTime() error {
    ok := userTokenModel.Update(t.ID, map[string]interface{}{
        "expire_time": int(time.Now().Unix()) + 3600 * 30,
    })
    if !ok {
        return errors.New("extended token expiration time failed")
    }
    return nil
}
