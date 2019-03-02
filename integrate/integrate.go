// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package integrate

import (
    "fmt"
    "sync"
    "errors"

    "github.com/dreamans/syncd/util/command"
    "github.com/dreamans/syncd/util/gofile"
    "github.com/dreamans/syncd/util/gostring"
    "github.com/dreamans/syncd/util/gopath"
)

const (
    COMMAND_LONG_TIMEOUT = 86400
    COMMAND_QUICK_TIMEOUT = 3600
    COMMAND_FLASH_TIMEOUT = 300
)

type Build struct {
    ID          int
    WorkSpace   string
    PackFile    string
    Repo        *Repo
    Scripts     string
    FnCallback  CallbackFn
    task        *command.Task
}

type CallbackFn func(int, error, string, string)

type buildTask struct {
    buildList   map[int]*Build
    mu          sync.Mutex
}

var buildTaskList = &buildTask{
    buildList: make(map[int]*Build),
}

func NewBuild(id int, workSpace, packFile string, scripts string, repo *Repo, fn CallbackFn) error {
    build := &Build{
        ID: id,
        WorkSpace: workSpace,
        PackFile: packFile,
        Repo: repo,
        Scripts: scripts,
        FnCallback: fn,
    }
    return build.build()
}

func StopBuild(id int) {
    buildTaskList.stop(id)
}

func (bt *buildTask) append(id int, build *Build) error {
    if bt.exists(id) {
        return errors.New("build task have exists")
    }
    bt.mu.Lock()
    defer bt.mu.Unlock()
    bt.buildList[id] = build
    return nil
}

func (bt *buildTask) remove(id int) {
    bt.mu.Lock()
    defer bt.mu.Unlock()
    delete(bt.buildList, id)
}

func (bt *buildTask) stop(id int) {
    bt.mu.Lock()
    defer bt.mu.Unlock()
    build, exists := bt.buildList[id]
    if exists {
        build.Terminate()
    }
}

func (bt *buildTask) exists(id int) bool {
    bt.mu.Lock()
    defer bt.mu.Unlock()
    _, exists := bt.buildList[id]
    return exists
}

func (b *Build) build() error {
    if err := buildTaskList.append(b.ID, b); err != nil {
        return err
    }
    go func() {
        var (
            err error
            rest, rests []*command.TaskResult
            tar string
        )
        rests, err = b.fetch()
        if err == nil {
            rest, err = b.integra()
            rests = append(rests, rest...)
        }

        buildTaskList.remove(b.ID)

        if gopath.Exists(b.PackFile) {
            tar = b.PackFile
        }
        b.FnCallback(b.ID, err, tar, string(gostring.JsonEncode(rests)))
    }()

    return nil
}

func (b *Build) fetch() ([]*command.TaskResult, error) {
    b.task = command.TaskNew(b.Repo.FetchCmd(), COMMAND_QUICK_TIMEOUT)
    b.task.Run()
    err := b.task.GetError()
    result := b.task.Result()
    b.task = nil
    return result, err
}

func (b *Build) integra() ([]*command.TaskResult, error) {
    scriptFile := gostring.JoinStrings(b.WorkSpace, "/", gostring.StrRandom(24), ".sh")
    if err := gofile.CreateFile(scriptFile, b.appendScriptHead(), 0744); err != nil {
        return nil, err
    }
    cmds := []string{
        "echo \"Now is\" `date`",
        "echo \"Run user is\" `whoami`",
        fmt.Sprintf("rm -f %s", b.PackFile),
        fmt.Sprintf("/bin/bash -c %s", scriptFile),
        fmt.Sprintf("rm -fr %s", b.WorkSpace),
    }
    b.task = command.TaskNew(cmds, COMMAND_LONG_TIMEOUT)
    b.task.Run()
    err := b.task.GetError()
    result := b.task.Result()
    b.task = nil
    return result, err
}

func (b *Build) appendScriptHead() []byte {
    s := gostring.JoinStrings(
        "#!/bin/bash\n\n",
        "#--------- build scripts env ---------\n",
        fmt.Sprintf("env_workspace=%s\n", b.WorkSpace),
        fmt.Sprintf("env_pack_file=%s\n", b.PackFile),
        b.Scripts,
    )
    return []byte(s)
}

func (b *Build) Terminate() {
    if b.task != nil {
        b.task.Terminate()
    }
}
