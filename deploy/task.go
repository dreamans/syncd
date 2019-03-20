// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "sync"
    "errors"
    "fmt"
)

type deployTask struct {
    deploys     map[int][]*Deploy
    mu          sync.Mutex
}

type CallbackFn func(int, int, int, []*ServerResult)

type TaskCallbackFn func(int, int)

var task = &deployTask{
    deploys: map[int][]*Deploy{},
}

func NewTask(id, mode int, deploys []*Deploy, startFn, finishFn CallbackFn, taskFn TaskCallbackFn) error {
	if exists := task.exists(id); exists {
        return errors.New(fmt.Sprintf("deploy task [id: %d] have exists", id))
    }
    task.append(id, deploys)
    go func() {
        taskStatus := STATUS_DONE
        for _, deploy := range deploys {
            if startFn != nil {
                rest, status := deploy.Result()
                startFn(id, deploy.ID, status, rest)
            }
            switch mode {
            case DEPLOY_PARALLEL:
                deploy.Parallel()
            default:
                deploy.Serial()
            }
            resultList, status := deploy.Result()
            if finishFn != nil {
                finishFn(id, deploy.ID, status, resultList)
            }
            if status == STATUS_FAILED {
                taskStatus = STATUS_FAILED
            }
        }
        task.remove(id)
        if taskFn != nil {
            taskFn(id, taskStatus)
        }
    }()
    return nil
}

func StopTask(id int) {
    task.stop(id)
}

func ExistsTask(id int) bool {
    return task.exists(id)
}

func StatusTask(id int) []*DeployResult {
    deploys, exists := task.get(id)
    if !exists {
        return nil
    }
    rests := []*DeployResult{}
    for _, deploy := range deploys {
        rest, s := deploy.Result()
        rests = append(rests, &DeployResult{
            ID: deploy.ID,
            Status: s,
            ServerRest: rest,
        })
    }
    return rests
}

func (t *deployTask) exists(id int) bool {
    t.mu.Lock()
    defer t.mu.Unlock()
    _, exists := t.deploys[id]
    return exists
}

func (t *deployTask) append(id int, deploys []*Deploy) {
    t.mu.Lock()
    defer t.mu.Unlock()
    t.deploys[id] = deploys
}

func (t *deployTask) remove(id int) {
    t.mu.Lock()
    defer t.mu.Unlock()
    delete(t.deploys, id)
}

func (t *deployTask) get(id int) ([]*Deploy, bool) {
    t.mu.Lock()
    defer t.mu.Unlock()
    deploys, exists := t.deploys[id]
    return deploys, exists
}

func (t *deployTask) stop(id int) {
    t.mu.Lock()
    defer t.mu.Unlock()
    deploys, exists := t.deploys[id]
    if exists {
        for _, deploy := range deploys {
            deploy.Terminate()
        }
    }
}