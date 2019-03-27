// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package command

import (
    "errors"
    "time"
)

type Task struct {
    Commands    []string
    done        bool
    timeout     int
    termChan    chan int
    err         error
    result      []*TaskResult
}

type TaskResult struct {
    Cmd         string      `json:"cmd"`
    Stdout      string      `json:"stdout"`
    Stderr      string      `json:"stderr"`
    Success     bool        `json:"success"`
}

func NewTask(cmds []string, timeout int) *Task {
    return &Task{
        Commands: cmds,
        termChan:  make(chan int),
        timeout: timeout,
    }
}

func (t *Task) Run() {
    for _, cmd := range t.Commands {
        result, err := t.next(cmd)
        t.result = append(t.result, result)
        if err != nil {
            t.err = errors.New("task run failed, " + err.Error())
            break
        }
    }
    t.done = true
}

func (t *Task) GetError() error {
    return t.err
}

func (t *Task) Result() []*TaskResult {
    return t.result
}

func (t *Task) next(cmd string) (*TaskResult, error) {
    result := &TaskResult{
        Cmd: cmd,
    }
    command := &Command{
        Cmd: cmd,
        Timeout: time.Second * time.Duration(t.timeout),
        TerminateChan: t.termChan,
    }
    var err error
    command, err = NewCmd(command)
    if err != nil {
        return result, err
    }
    if err := command.Run(); err != nil {
        result.Stderr = command.Stderr()
        return result, err
    }
    result.Stdout = command.Stdout()
    result.Success = true
    return result, nil
}

func (t *Task) Terminate() {
    if !t.done {
        t.termChan <- 1
    }
}