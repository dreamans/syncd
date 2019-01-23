// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gostring

import (
    "strconv"
    "math/rand"
    "time"
    "crypto/md5"
    "encoding/hex"
)

func Str2Int(s string) int {
    i, _ := strconv.Atoi(s)
    return i
}

func Int2Str(i int) string {
    return strconv.Itoa(i)
}

func StrRandom(l int) string {
    str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    bytes := []byte(str)
    result := []byte{}
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < l; i++ {
        result = append(result, bytes[r.Intn(len(bytes))])
    }
    return string(result)
}

func StrMd5(s string) string {
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(s))
    return hex.EncodeToString(md5Ctx.Sum(nil))
}
