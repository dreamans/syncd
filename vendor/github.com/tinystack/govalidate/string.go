// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package govalidate

import (
    "reflect"
    "strconv"
)

func init() {
    validateHandlerMap[VALID_STR_MIN] = strMinHandler
    validateHandlerMap[VALID_STR_MAX] = strMaxHandler
    validateHandlerMap[VALID_STR_LEN] = strEqLenHandler
    validateHandlerMap[VALID_STR_RANGE] = strRangeHandler
}

func strMinHandler(value reflect.Value, params []string) bool {
    return strMinMaxHandler(value, params, 'l')
}

func strMaxHandler(value reflect.Value, params []string) bool {
    return strMinMaxHandler(value, params, 'g')
}

func strEqLenHandler(value reflect.Value, params []string) bool {
    return strMinMaxHandler(value, params, 'e')
}

func strRangeHandler(value reflect.Value, params []string) bool {
    return strMinMaxHandler(value, params, 'r')
}

func strMinMaxHandler(value reflect.Value, params []string, op byte) bool {
    var smin, smax, slen int
    if len(params) == 1 {
        slen, _ = strconv.Atoi(params[0])
    }
    if len(params) == 2 {
        smin, _ = strconv.Atoi(params[0])
        smax, _ = strconv.Atoi(params[1])
    }
    if value.Kind() == reflect.String {
        s := value.String()
        l := len(s)
        if l == 0 {
            return true
        }
        switch op {
        case 'l':
            if l < slen {
                return false
            }
        case 'e':
            if l != slen {
                return false
            }
        case 'g':
            if l > slen {
                return false
            }
        case 'r':
            if l < smin || l > smax {
                return false
            }
        }
    }
    return true
}
