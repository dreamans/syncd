// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golog

import (
    "os"
    "sync"
)

type FileHandler struct {
    filename    string
    file        *os.File
    mu          sync.Mutex
}

func NewFileHandler(path string) *FileHandler {
    return &FileHandler{
        filename: path,
    }
}

func (f *FileHandler) Write(p []byte) (n int, err error) {
    f.mu.Lock()
    defer f.mu.Unlock()

    if err = f.openFileHandler(); err != nil {
        return 0, err
    }

    return f.file.Write(p)
}

func (f *FileHandler) Close() error {
    f.mu.Lock()
    defer f.mu.Unlock()
    return f.close()
}

func (f *FileHandler) close() error {
    if f.file == nil {
        return nil
    }
    err := f.file.Close()
    f.file = nil
    return err
}

func (f *FileHandler) openFileHandler() error {
    file, err := os.OpenFile(f.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
    if err != nil {
        return err
    }
    f.file = file
    return nil
}
