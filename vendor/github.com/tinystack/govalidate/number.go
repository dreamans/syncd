// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package govalidate

import (
    "reflect"
    "strings"
    "strconv"
)

func init() {
    validateHandlerMap[VALID_NUM_MIN] = numMinHandler
    validateHandlerMap[VALID_NUM_MAX] = numMaxHandler
    validateHandlerMap[VALID_NUM_RANGE] = numRangeHandler
}

func numMinHandler(value reflect.Value, params []string) bool {
    return numMinMaxHandler(value, params, 'l')
}

func numMaxHandler(value reflect.Value, params []string) bool {
    return numMinMaxHandler(value, params, 'g')
}

func numRangeHandler(value reflect.Value, params []string) bool {
    return numMinMaxHandler(value, params, 'r')
}

func numMinMaxHandler(value reflect.Value, params []string, op byte) bool {
    var val, pa1, pa2 float64
    if op == 'r' {
        if len(params) != 2 {
            return true
        }
        pa2, _ = strconv.ParseFloat(params[1], 64)
    } else {
        if len(params) != 1 {
            return true
        }
    }
    pa1, _ = strconv.ParseFloat(params[0], 64)

    kind := value.Kind().String()
    switch true {
    case strings.Index(kind, "int") == 0:
        val = float64(value.Int())
    case strings.Index(kind, "uint") == 0:
        val = float64(value.Uint())
    case strings.Index(kind, "float") == 0:
        val = float64(value.Float())
    }

    switch op {
    case 'l':
        if val < pa1 {
            return false
        }
    case 'g':
        if val > pa1 {
            return false
        }
    case 'r':
        if val < pa1 || val > pa2 {
            return false
        }
    }

    return true
}
