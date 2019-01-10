// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package govalidate

import (
    "reflect"
    "strings"
)

const (
    VALID_REQUIRED = "required"
    VALID_NUM_MIN = "int_min"
    VALID_NUM_MAX = "int_max"
    VALID_NUM_RANGE = "num_range"
    VALID_STR_MIN = "str_min"
    VALID_STR_MAX = "str_max"
    VALID_STR_LEN = "str_len"
    VALID_STR_RANGE = "str_range"
    VALID_EMAIL = "email"
    VALID_MOBILE = "mobile"
)

var validateHandlerMap = make(map[string]validateHandler)

type Validate struct {
    Stru      interface{}
    Items     []*Item
    Faileds   []*Failed
}

type Item struct {
    Name     string
    Valids   []*Unit
    Value    reflect.Value
}

type Failed struct {
    Name    string
    Valid   string
    Msg     string
}

type Unit struct {
    Valid       string
    Params      []string
    MsgFormat   string
}

type validateHandler func(value reflect.Value, param []string) bool

func NewValidate(stru interface{}) *Validate {
    valid := &Validate{
        Stru: stru,
    }
    valid.parse()
    valid.valid()
    return valid
}

func (v *Validate) Pass() (pass bool) {
    pass = true
    if len(v.Faileds) > 0 {
        pass = false
    }
    return
}

func (v *Validate) LastFailed() *Failed {
    if v.Pass() {
        return nil
    }
    return v.Faileds[len(v.Faileds) - 1]
}

func newValidItem(fieldName, valid, errmsg string, value reflect.Value) *Item {
    vitem := &Item{
        Name: fieldName,
        Value: value,
    }

    var errMessage = make(map[string]string)
    errMsgs := strings.Split(errmsg, "|")
    for _, msg := range errMsgs {
        sp := strings.Split(msg, "=")
        if len(sp) != 2 {
            continue
        }
        errMessage[strings.TrimSpace(sp[0])] = strings.TrimSpace(sp[1])
    }

    valids := strings.Split(valid, "|")
    for _, valid := range valids {
        sp := strings.Split(valid, "=")
        u := &Unit{
            Valid: strings.TrimSpace(sp[0]),
        }
        if len(sp) > 1 {
            params := strings.Split(sp[1], ",")
            for _, pa := range params {
                pa = strings.TrimSpace(pa)
                if pa != "" {
                    u.Params = append(u.Params, pa)
                }
            }
        }
        if msg, ok := errMessage[u.Valid]; ok {
            u.MsgFormat = msg
        } else {
            u.MsgFormat = ":attr field valid failed"
        }
        vitem.Valids = append(vitem.Valids, u)
    }

    return vitem
}

func (v *Validate) parse() {
    fields := reflect.ValueOf(v.Stru).Elem()
    struName := fields.Type().Name()
    for i := 0; i < fields.NumField(); i++ {
        field := fields.Type().Field(i)
        valid := field.Tag.Get("valid")
        errmsg := field.Tag.Get("errmsg")
        if valid == "" {
            continue
        }
        name := strings.Join([]string{struName, field.Name}, ".")
        v.Items = append(v.Items, newValidItem(name, valid, errmsg, fields.FieldByName(field.Name)))
    }
}

func (v *Validate) valid() error {
    for _, item := range v.Items {
        for _, u := range item.Valids {
            if handler, ok := validateHandlerMap[u.Valid]; ok {
                if ok := handler(item.Value, u.Params); !ok {
                    fail := &Failed{
                        Name: item.Name,
                        Valid: u.Valid,
                        Msg: strings.Replace(u.MsgFormat, ":attr", item.Name, -1),
                    }
                    v.Faileds = append(v.Faileds, fail)
                }
            }
        }
    }

    return nil
}
