// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package integrate

import (
    "github.com/dreamans/syncd/util/command"
)

const (
    COMMAND_LONG_TIMEOUT = 86400
    COMMAND_QUICK_TIMEOUT = 3600
    COMMAND_FLASH_TIMEOUT = 300
)

type Build struct {
    WorkSpace   string
    PackFile    string
	Repo        *Repo
	Scripts		string
    task        *command.Task
}

func NewBuild(workSpace, packFile string, scripts string, repo *Repo) *Build {
    build := &Build{
        WorkSpace: workSpace,
        PackFile: packFile,
		Repo: repo,
		Scripts: scripts,
    }
    return build
}

func (b *Build) Build() {
	go func() {
		b.Fetch()
	}()
}

func (b *Build) fetch() ([]*command.TaskResult, error) {
    b.task = command.TaskNew(b.Repo.FetchCmd(), COMMAND_QUICK_TIMEOUT)
    b.task.Run()
    err := b.task.GetError()
    result := b.task.Result()
    b.task = nil
    return result, err
}

func (b *Build) integra() {
	
}

func (b *Build) Terminate() {
    if b.task != nil {
        b.task.Terminate()
    }
}
