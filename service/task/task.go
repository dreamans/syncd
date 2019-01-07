// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
    "time"
    "errors"
    "sync"

    "github.com/tinystack/goutil/gostring"
    "github.com/tinystack/gocmd"
)

type Task struct {
    name        string
    cmd         []string
    timeout     int
    termChan    chan int
    err         []error
    stdout      []string
    stderr      []string
    wg          sync.WaitGroup
}

const (
    TASK_REPO_RESET = "repo_reset"
    TASK_REPO_TAG_LIST = "repo_tag_list"
    TASK_REPO_UPDATE = "repo_update"
    TASK_REPO_COMMIT_LIST = "repo_commit_list"
    TASK_SERVER_CHECK = "server_check"
    TASK_REPO_DEPLOY = "repo_deploy"
)

func TaskCreate(name string, cmd []string, timeout int) *Task {
    task := &Task{
        name: name,
        cmd: cmd,
        timeout: timeout,
    }
    return task
}

func (t *Task) TaskAdd() {
    t.wg.Add(1)
}

func (t *Task) TaskDone() {
    t.wg.Done()
}

func (t *Task) TaskWait() {
    t.wg.Wait()
}

func (t *Task) TaskRun() {
    for _, cmd := range t.cmd {
        if err := t.next(cmd); err != nil {
            t.err = append(t.err, errors.New("task run failed, " + err.Error()))
            break
        }
    }
}

func (t *Task) next(cmd string) error {
    t.termChan = make(chan int)
    command := gocmd.Command{
        Cmd: cmd,
        Timeout: time.Second * time.Duration(t.timeout),
        TerminateChan: t.termChan,
    }
    if err := command.Run(); err != nil {
        stderr := string(command.Stderr())
        if stderr != "" {
            t.stderr = append(t.stderr, stderr)
        }
        return err
    }
    stdout := string(command.Stdout())
    if stdout != "" {
        t.stdout = append(t.stdout, stdout)
    }
    return nil
}

func (t *Task) Stdout() string {
    return gostring.JoinSepStrings("\n", t.stdout...)
}

func (t *Task) Stderr() string {
    return gostring.JoinSepStrings("\n", t.stderr...)
}

func (t *Task) LastError() error {
    if len(t.err) == 0 {
        return nil
    }
    return t.err[len(t.err)-1]
}

func (t *Task) Terminate() {
    t.termChan <- 1
}

