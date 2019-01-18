// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package command

import (
    "testing"
    "time"
    "strings"
)

func TestCmdRun(t *testing.T) {
    c := &Command{
        Cmd: "echo 'syncd'",
        Timeout: time.Second * 60,
    }
    var err error 
    if c, err = NewCmd(c); err != nil {
        t.Error(err)
    }
    if err = c.Run(); err != nil {
        t.Error(err)
    }
    output := c.Stdout()
    if output != "syncd" {
        t.Errorf("cmd `echo syncd` output '%s' not eq 'syncd'", output)
    }
}

func TestCmdTimeout(t *testing.T) {
    c := &Command{
        Cmd: "sleep 2",
        Timeout: time.Second * 1,
    }
    var err error
    if c, err = NewCmd(c); err != nil {
        t.Error(err)
    }
    err = c.Run()
    if err == nil {
        t.Error("cmd should run timeout and output errmsg, but err is nil")
    }
    if strings.Index(err.Error(), "cmd run timeout") == -1 {
        t.Errorf("cmd run timeout output '%s' prefix not eq 'cmd run timeout'", err.Error())
    }
}

func TestCmdTerminate(t *testing.T) {
    tChan := make(chan int)
    c := &Command{
        Cmd: "sleep 5",
        TerminateChan: tChan,
    }
    var err error
    if c, err = NewCmd(c); err != nil {
        t.Error(err)
    }
    go func() {
        err = c.Run()
        if err == nil {
            t.Error("cmd should be terminated and output errmsg, but err is nil")
        }
        if strings.Index(err.Error(), "cmd is terminated") == -1 {
            t.Errorf("cmd is terminated output '%s' prefix not eq 'cmd is terminated'", err.Error())
        }
    }()

    tChan<- 1
}

