// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
    "errors"
    "fmt"
    "time"
    "strings"

    "github.com/dreamans/syncd"
    "github.com/dreamans/syncd/util/gostring"
    "github.com/dreamans/syncd/util/goaes"
)

type Login struct {
    UserId      int
    RoleId      int
    Username    string
    Password    string
    Email       string
    Truename    string
    Mobile      string
    Token       string
}

func (login *Login) Logout() error {
    token := &Token{
        UserId: login.UserId,
    }
    return token.DeleteByUserId()
}

func (login *Login) Login() error {
    u := &User{}
    if login.Username != "" {
        u.Username = login.Username
    }
    if login.Email != "" {
        u.Email = login.Email
    }
    if err := u.Detail(); err != nil {
        return errors.New("username or password incorrect")
    }
    loginPassword := gostring.StrMd5(gostring.JoinStrings(login.Password, u.Salt))
    if u.Password != loginPassword {
        return errors.New("password incorrect")
    }

    if u.Status != 1 {
        return errors.New("user is locked")
    }

    login.UserId = u.ID
    if err := login.createToken(); err != nil {
        return errors.New("token create failed")
    }

    return nil
}

func (login *Login) ValidateToken() error {
    authTokenBytes, err := gostring.Base64UrlDecode(login.Token)
    if err != nil {
        return errors.New("token check failed, can not decode raw token")
    }
    tokenValBytes, err := goaes.Decrypt(syncd.App.CipherKey, authTokenBytes)
    if err != nil {
        return errors.New("token check failed, can not decrypt raw token")
    }
    tokenArr := strings.Split(string(tokenValBytes), "\t")
    if len(tokenArr) != 2 {
        return errors.New("token check failed, len wrong")
    }
    token := &Token{
        UserId: gostring.Str2Int(tokenArr[0]),
        Token: tokenArr[1],
    }
    if ok := token.ValidateToken(); !ok {
        return errors.New("token check failed, maybe your account is logged in on another device or token expired")
    }

    //get user detail
    user := &User{
        ID: token.UserId,
    }
    if err := user.Detail(); err != nil {
        return errors.New("token check failed, user detail get failed")
    }

    if user.Status != 1 {
        return errors.New("user is locked")
    }

    login.UserId = user.ID
    login.Username = user.Username
    login.Email = user.Email
    login.Truename = user.Truename
    login.Mobile = user.Mobile
    login.RoleId = user.RoleId

    return nil
}

func (login *Login) createToken() error {
    loginKey := gostring.StrRandom(40)
    loginRaw := fmt.Sprintf("%d\t%s", login.UserId, loginKey)
    var (
        err error
        tokenBytes []byte
    )
    tokenBytes, err = goaes.Encrypt(syncd.App.CipherKey, []byte(loginRaw))
    if err != nil {
        return err
    }
    login.Token = gostring.Base64UrlEncode(tokenBytes)

    token := &Token{
        UserId: login.UserId,
        Token: loginKey,
        Expire: int(time.Now().Unix()) + 86400 * 30,
    }
    if err := token.CreateOrUpdate(); err != nil {
        return err
    }
    return nil
}

