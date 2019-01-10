// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package govalidate

import (
    "reflect"
    "regexp"
)

func init() {
    validateHandlerMap[VALID_EMAIL] = emailHandler
}

func emailHandler(value reflect.Value, params []string) bool {
    if value.Kind() != reflect.String {
        return false
    }
    email := value.String()
    if email == "" {
        return true
    }
    pattern := `^[0-9A-Za-z][\.\-_0-9A-Za-z]*\@[0-9A-Za-z\-]+(\.[0-9A-Za-z]+)+$`
    reg := regexp.MustCompile(pattern)
    return reg.MatchString(email)
}
