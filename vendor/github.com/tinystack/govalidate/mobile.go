// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package govalidate

import (
    "regexp"
    "reflect"
)

func init() {
    validateHandlerMap[VALID_MOBILE] = mobileHandler
}

func mobileHandler(value reflect.Value, params []string) bool {
    if value.Kind() != reflect.String {
        return false
    }
    mobile := value.String()
    if mobile == "" {
        return true
    }
    pattern := `^1[3456789]\d{9}$`
    reg := regexp.MustCompile(pattern)
    return reg.MatchString(mobile)
}
