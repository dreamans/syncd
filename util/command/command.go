// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package command

import (
    "bytes"
    "os/exec"
    "time"
    "syscall"
    "errors"
    "fmt"
    "strings"
)

const (
    DEFAULT_RUM_TIMEOUT = 3600
)

type Command struct {
    Cmd             string
    Timeout         time.Duration
    TerminateChan   chan int
    Setpgid         bool
    command         *exec.Cmd
    stdout          bytes.Buffer
    stderr          bytes.Buffer
}

func NewCmd(c *Command) (*Command, error) {
    if c.Timeout == 0 * time.Second {
        c.Timeout = DEFAULT_RUM_TIMEOUT * time.Second
    }
    if c.TerminateChan == nil {
        c.TerminateChan = make(chan int)
    }
    cmd := exec.Command("/bin/bash", "-c", c.Cmd)
    if c.Setpgid {
        cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
    }
    cmd.Stderr = &c.stderr
    cmd.Stdout = &c.stdout
    c.command = cmd

    return c, nil
}

func (c *Command) Run() error {
    if err := c.command.Start(); err != nil {
        return err
    }

    errChan := make(chan error)
    go func(){
        errChan <- c.command.Wait()
        defer close(errChan)
    }()

    var err error
    select {
    case err = <-errChan:
    case <-time.After(c.Timeout):
        err = c.terminate()
        if err == nil {
            err = errors.New(fmt.Sprintf("cmd run timeout, cmd [%s], time[%v]", c.Cmd, c.Timeout))
        }
    case <-c.TerminateChan:
        err = c.terminate()
        if err == nil {
            err = errors.New(fmt.Sprintf("cmd is terminated, cmd [%s]", c.Cmd))
        }
    }
    return err
}

func (c *Command) Stderr() string {
    return strings.TrimSpace(string(c.stderr.Bytes()))
}

func (c *Command) Stdout() string {
    return strings.TrimSpace(string(c.stdout.Bytes()))
}

func (c *Command) terminate() error {
    if c.Setpgid {
        return syscall.Kill(-c.command.Process.Pid, syscall.SIGKILL)
    } else {
        return syscall.Kill(c.command.Process.Pid, syscall.SIGKILL)
    }
}