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
)

type Build struct {
    repo        *Repo
    local       string
    packFile    string
    scriptFile  string
    task        *command.Task
    status      int
    stime       int
    etime       int
}

const (
    STATUS_INIT = 1
    STATUS_ING = 2
    STATUS_DONE = 3
)

func NewBuild(repo *Repo, local, packFile, scripts string) (*Build, error) {
    build := &Build{
        repo: repo,
        local: local,
        packFile: packFile,
        status: STATUS_INIT,
    }
    if err := build.createScriptFile(scripts); err != nil {
        return build, err
    }
    b.initBuildTask()

    return build, nil
}

func (b *Build) createScriptFile(scripts string) error {
    b.scriptFile := gostring.JoinStrings(b.local, "/", gostring.StrRandom(24), ".sh")
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
        fmt.Sprintf("rm -fr %s", b.local),
    }...)
    b.task := command.NewTask(cmds, 86400)
}

func (b *Build) Run() {
    b.status = STATUS_ING
    b.stime = int(time.Now().Unix())
    b.task.Run()
    b.status = STATUS_DONE
    b.stime = int(time.Now().Unix())
}