// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package build

import (
    "sync"
    "errors"
    "fmt"

    "github.com/dreamans/syncd/util/command"
)

type buildTask struct {
    builds  map[int]*Build
    mu      sync.Mutex
}

type CallbackFn func(int, string, *Result, []*command.TaskResult)

var task = &buildTask{
    builds: make(map[int]*Build),
}

func NewTask(id int, build *Build, fn CallbackFn) error {
    if exists := task.exists(id); exists {
        return errors.New(fmt.Sprintf("build task [id: %d] have exists", id))
    }
    task.append(id, build)
    go func() {
        build.Run()
        task.remove(id)
        if fn != nil {
            fn(id, build.PackRealFile(), build.Result(), build.Output())
        }
    }()
    return nil
}

func StopTask(id int) {
    task.stop(id)
}

func StatusTask(id int) (*Result, []*command.TaskResult, error) {
    build, exists := task.get(id)
    if !exists {
        return nil, nil, errors.New(fmt.Sprintf("build task [id: %d] not exists", id))
    }
    return build.Result(), build.Output(), nil
}

func (t *buildTask) exists(id int) bool {
    t.mu.Lock()
    defer t.mu.Unlock()
    _, exists := t.builds[id]
    return exists
}

func (t *buildTask) append(id int, build *Build) {
    t.mu.Lock()
    defer t.mu.Unlock()
    t.builds[id] = build
}

func (t *buildTask) remove(id int) {
    t.mu.Lock()
    defer t.mu.Unlock()
    delete(t.builds, id)
}

func (t *buildTask) get(id int) (*Build, bool) {
    t.mu.Lock()
    defer t.mu.Unlock()
    build, exists := t.builds[id]
    return build, exists
}

func (t *buildTask) stop(id int) {
    t.mu.Lock()
    defer t.mu.Unlock()
    build, exists := t.builds[id]
    if exists {
        build.Terminate()
    }
}