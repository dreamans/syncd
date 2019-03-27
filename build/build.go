// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package build

import (
    "time"
    "fmt"

    "github.com/dreamans/syncd/util/gostring"
    "github.com/dreamans/syncd/util/command"
    "github.com/dreamans/syncd/util/gofile"
    "github.com/dreamans/syncd/util/gopath"
)

type Build struct {
    repo        *Repo
    local       string
    tmp         string
    packFile    string
    scriptFile  string
    task        *command.Task
    result      *Result
}

const (
    STATUS_INIT = 1
    STATUS_ING = 2
    STATUS_DONE = 3
    STATUS_FAILED = 4
)

const (
    COMMAND_TIMEOUT = 86400
)

func NewBuild(repo *Repo, local, tmp, packFile, scripts string) (*Build, error) {
    build := &Build{
        repo: repo,
        local: local,
        tmp: tmp,
        packFile: packFile,
        result: &Result{
            status: STATUS_INIT,
        },
    }
    if err := build.createScriptFile(scripts); err != nil {
        return build, err
    }
    build.initBuildTask()

    return build, nil
}

func (b *Build) createScriptFile(scripts string) error {
    b.scriptFile = gostring.JoinStrings(b.tmp, "/", gostring.StrRandom(24), ".sh")
    s := gostring.JoinStrings(
        "#!/bin/bash\n\n",
        "#--------- build scripts env ---------\n",
        fmt.Sprintf("env_workspace=%s\n", b.local),
        fmt.Sprintf("env_pack_file=%s\n", b.packFile),
        scripts,
    )
    if err := gofile.CreateFile(b.scriptFile, []byte(s), 0744); err != nil {
        return err
    }
    return nil
}

func (b *Build) initBuildTask() {
    cmds := b.repo.Fetch()
    cmds = append(cmds, []string{
        "echo \"Now is\" `date`",
        "echo \"Run user is\" `whoami`",
        fmt.Sprintf("rm -f %s", b.packFile),
        fmt.Sprintf("/bin/bash -c %s", b.scriptFile),
        fmt.Sprintf("rm -f %s", b.scriptFile),
        fmt.Sprintf("rm -fr %s", b.local),
        "echo \"Compile completed\" `date`",
    }...)
    b.task = command.NewTask(cmds, COMMAND_TIMEOUT)
}

func (b *Build) Run() {
    b.result.status = STATUS_ING
    b.result.stime = int(time.Now().Unix())
    b.task.Run()
    if err := b.task.GetError(); err != nil {
        b.result.status = STATUS_FAILED
        b.result.err = err
    } else {
        b.result.status = STATUS_DONE
    }
    b.result.etime = int(time.Now().Unix())
}

func (b *Build) Result() *Result {
    return b.result
}

func (b *Build) Output() []*command.TaskResult{
    return b.task.Result()
}

func (b *Build) PackFile() string {
    return b.packFile
}

func (b *Build) PackRealFile() string {
    if gopath.Exists(b.packFile) {
        return b.packFile
    }
    return ""
}

func (b *Build) Terminate() {
    if b.task != nil {
        b.task.Terminate()
    }
}

type Result struct {
    err     error
    status  int
    stime   int
    etime   int
}

func (r *Result) During() int {
    return r.etime - r.stime
}

func (r *Result) Status() int {
    return r.status
}

func (r *Result) GetError() error {
    return r.err
}