// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package govalidate

import (
    "reflect"
)

func init() {
    validateHandlerMap[VALID_REQUIRED] = requiredHandler
}

func requiredHandler(value reflect.Value, param []string) bool {
    switch value.Kind() {
    case reflect.String:
        if value.String() == "" {
            return false
        }
    case reflect.Slice:
        if value.Slice(0, 1).Len() == 0 {
            return false
        }
    }

    return true
}
