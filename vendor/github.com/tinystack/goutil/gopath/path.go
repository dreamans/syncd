// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gopath

import (
    "errors"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

func CurrentPath() (string, error) {
    file, err := exec.LookPath(os.Args[0])
    if err != nil {
        return "", err
    }
    path, err := filepath.Abs(file)
    if err != nil {
        return "", err
    }
    i := strings.LastIndex(path, "/")
    if i < 0 {
        i = strings.LastIndex(path, "\\")
    }
    if i < 0 {
        return "", errors.New(`error: Can't find "/" or "\".`)
    }
    return string(path[0 : i]), nil
}

func CreatePath(path string) error {
    exists, err := PathExists(path)
    if err != nil {
        return err
    }
    if !exists {
        if err := os.Mkdir(path, os.ModePerm); err != nil {
            return err
        }
    }
    return nil
}

func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

func RemovePath(path string) error {
    return os.RemoveAll(path)
}
