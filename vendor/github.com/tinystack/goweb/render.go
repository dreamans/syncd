// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goweb

import (
    "encoding/json"
    "net/http"
)

type JSON map[string]interface{}

func renderJson(w http.ResponseWriter, obj interface{}) error {
    jsonBytes, err := json.Marshal(obj)
    if err != nil {
        return err
    }
    w.Write(jsonBytes)
    return nil
}
