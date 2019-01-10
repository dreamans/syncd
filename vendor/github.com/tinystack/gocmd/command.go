// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gocmd

import (
    "os/exec"
    "bytes"
    "bufio"
    "io"
    "time"
    "errors"
    "syscall"
    "fmt"
)

type Command struct {
    Cmd             string
    Timeout         time.Duration
    TerminateChan   chan int
    OutputFnHandle  func(line []byte)
    command         *exec.Cmd
    stdout          []byte
    stderr          bytes.Buffer
}

func (c *Command) Run() error {
    var err error

    cmd := c.cmdInit()
    cmd.Stderr = &c.stderr

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return err
    }
    if err := cmd.Start(); err != nil {
        return err
    }
    reader := bufio.NewReader(stdout)

    outputBytesChan := make(chan []byte)
    go func(ch chan []byte){
        var outputBytes []byte
        for {
            line, _, err := reader.ReadLine()
            if err != nil || io.EOF == err {
                break
            }
            if c.OutputFnHandle != nil {
                c.OutputFnHandle(line)
            }
            outputBytes = append(outputBytes, line...)
            outputBytes = append(outputBytes, '\n')
        }
        ch <-outputBytes
        defer close(outputBytesChan)
    }(outputBytesChan)

    errChan := make(chan error)
    go func(){
        errChan <- cmd.Wait()
        defer close(errChan)
    }()

    err = nil
    select {
        case err = <-errChan:
        case <-time.After(c.Timeout):
            err = c.terminate()
            if err == nil {
                err = errors.New(fmt.Sprintf("command run timeout and forced to terminate, cmd [%s], timeout [%v]", c.Cmd, c.Timeout))
            }
        case <-c.TerminateChan:
            err = c.terminate()
            if err == nil {
                err = errors.New(fmt.Sprintf("command is terminated, cmd [%s]", c.Cmd))
            }
    }
    c.stdout = <-outputBytesChan
    return err
}

func (c *Command) Stdout() []byte {
    return c.stdout
}

func (c *Command) Stderr() []byte {
    return c.stderr.Bytes()
}

func (c *Command) cmdInit () *exec.Cmd {
    if c.Timeout == 0 * time.Second {
        c.Timeout = 3600 * time.Second
    }
    if c.TerminateChan == nil {
        c.TerminateChan = make(chan int)
    }

    c.command = exec.Command("/bin/bash", "-c", c.Cmd)
    c.command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

    return c.command
}

func (c *Command) terminate() error {
    return syscall.Kill(-c.command.Process.Pid, syscall.SIGKILL)
}
