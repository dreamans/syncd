// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
    "time"
    "errors"
    "fmt"
    "sync"

    "github.com/tinystack/goutil/gostring"
    "github.com/tinystack/gocmd"
    taskModel "github.com/tinystack/syncd/model/task"
    taskLogModel "github.com/tinystack/syncd/model/task_log"
)

type Task struct {
    name        string
    cmd         []string
    key         string
    ctime       int
    ftime       int
    status      int
    termChan    chan int
    err         []error
    stdout      []string
    stderr      []string
    wg          sync.WaitGroup
}

func TaskCreate(name string, cmd []string) (*Task, error) {
    task := &Task{
        name: name,
        cmd: cmd,
        key: gostring.StrRandom(32),
        ctime: int(time.Now().Unix()),
        status: 1,
    }
    if ok := taskModel.Create(&taskModel.Task{
        Name: task.name,
        Key: task.key,
        Cmd: gostring.JoinSepStrings("\n", task.cmd...),
        Status: 1,
    }); !ok {
        return nil, errors.New("create task data failed")
    }

    task.recordTaskLog(fmt.Sprintf("create task [%s]", task.name))

    return task, nil
}

func (t *Task) Key() string {
    return t.key
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

func (t *Task) TaskRun() error {
    var ok bool
    ok = taskModel.UpdateByKey(t.key, map[string]interface{}{
        "status": 2,
    })
    if !ok {
        return errors.New("update task status failed")
    }
    t.recordTaskLog(fmt.Sprintf("start run task [%s]", t.name))
    for _, cmd := range t.cmd {
        if err := t.next(cmd); err != nil {
            t.err = append(t.err, errors.New("task run failed, " + err.Error()))
            break
        }
    }

    status := 3
    if len(t.err) > 0 {
        status = 4
    }
    ok = taskModel.UpdateByKey(t.key, map[string]interface{}{
        "status": status,
        "ftime": int(time.Now().Unix()),
    })
    if !ok {
        return errors.New("update task status failed")
    }

    result := "success"
    if status == 4 {
        result = "failed"
    }
    t.recordTaskLog(fmt.Sprintf("task [%s] run %s", t.name, result))

    return nil
}

func (t *Task) next(cmd string) error {
    t.termChan = make(chan int)
    command := gocmd.Command{
        Cmd: cmd,
        Timeout: time.Second * 3600,
        TerminateChan: t.termChan,
    }

    t.recordTaskLog(fmt.Sprintf("run cmd [%s]", cmd))

    if err := command.Run(); err != nil {
        stderr := string(command.Stderr())
        if stderr != "" {
            t.stderr = append(t.stderr, stderr)
            t.recordTaskLog(stderr)
        }
        t.recordTaskLog(err.Error())
        return err
    }

    stdout := string(command.Stdout())
    if stdout != "" {
        t.stdout = append(t.stdout, stdout)
        t.recordTaskLog(stdout)
    }

    return nil
}

func (t *Task) recordTaskLog(msg string) {
    taskLogModel.Create(&taskLogModel.TaskLog{
        Key: t.key,
        Content: msg,
    })
}

func (t *Task) Stdout() []string {
    return t.stdout
}

func (t *Task) Stderr() []string {
    return t.stderr
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

