// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "errors"
    "path"
    "sync"
)

const (
    STATUS_INIT = 0
    STATUS_RUNNING = 1
    STATUS_DONE = 2
    STATUS_FAILED = 3
    STATUS_TERMINATE = 4
)

type Deploy struct {
    ID              int
    Srvs            []*Server
    PreCmd          string
    PostCmd         string
    DeployPath      string
    DeployTmpPath   string
    PackFile        string
    PackFileName    string
    CallbackFn      deployCallbackFn
    StartCallbackFn	deployCallbackFn
    status          int
    wg              sync.WaitGroup
}

type taskCallbackFn func(int, bool)

type deployCallbackFn func(int, int, *ServerStatus)

var deployTaskList = &deployTask{
    deployList: make(map[int][]*Deploy),
}

func NewDepoly(deploy *Deploy) (*Deploy, error) {
    if deploy.PackFile== "" {
        return nil, errors.New("new deploy failed, deploy.PackFile can not be empty")
    }
    if deploy.DeployPath == "" {
        return nil, errors.New("new deploy failed, deploy.DeployPath can not be empty")
    }
    if deploy.DeployTmpPath == "" {
        deploy.DeployTmpPath = "~/.syncd"
    }
    if deploy.PackFileName == "" {
        deploy.PackFileName = path.Base(deploy.PackFile)
    }

    return deploy, nil
}

func DeployGroups (id int, deployGroups []*Deploy, callbackFn taskCallbackFn) error {
    if err := deployTaskList.append(id, deployGroups); err != nil {
        return err
    }
    go func() {
        haveError := false
        for _, dep := range deployGroups {
            dep.deploy()
            if dep.status == STATUS_FAILED {
                haveError = true
            }
        }
        deployTaskList.remove(id)
        callbackFn(id, haveError)
    }()
    return nil
}

func StopDeploy(id int) {
    deployTaskList.stop(id)
    deployTaskList.remove(id)
}

func (deploy *Deploy) deploy() {
    deploy.status = STATUS_RUNNING
    var srvError error
    for _, srv := range deploy.Srvs {
        if deploy.StartCallbackFn != nil {
            deploy.StartCallbackFn(deploy.ID, srv.ID, nil)
        }
        // check task exists
        var srvStatus *ServerStatus
        if exists := deployTaskList.exists(deploy.ID); exists {
            srv.Deploy(deploy)
            srvStatus = srv.Status()
        } else {
            srvStatus = &ServerStatus{
                Status: STATUS_TERMINATE, 
                Error: errors.New("deploy task run failed, cmd is terminated"),
            }
        }

        if deploy.CallbackFn != nil {
            deploy.CallbackFn(deploy.ID, srv.ID, srvStatus)
        }
        if srvStatus.Error != nil {
            srvError = srvStatus.Error
        }
    }
    if srvError == nil {
        deploy.status = STATUS_DONE
    } else {
        deploy.status = STATUS_FAILED
    }
}

func (deploy *Deploy) ParalDeploy() func() {
    deploy.status = STATUS_RUNNING
    for _, srv := range deploy.Srvs {
        deploy.wg.Add(1)
        go func() {
            srv.Deploy(deploy)
            defer deploy.wg.Done()
        }()
    }

    return func() {
        deploy.wg.Wait()
        deploy.status = STATUS_DONE
    }
}

func (deploy *Deploy) Terminate() {
    if deploy.status == STATUS_RUNNING {
        for _, srv := range deploy.Srvs {
            srv.Terminate()
        }
    }
}

type deployTask struct {
    deployList  map[int][]*Deploy
    mu          sync.Mutex
}

func (dt *deployTask) exists(id int) bool {
    dt.mu.Lock()
    defer dt.mu.Unlock()
    _, exists := dt.deployList[id]
    return exists
}

func (dt *deployTask) append(id int, deploy []*Deploy) error {
    if dt.exists(id) {
        return errors.New("deploy task have exists")
    }
    dt.mu.Lock()
    defer dt.mu.Unlock()
    dt.deployList[id] = deploy
    return nil
}

func (dt *deployTask) remove(id int) {
    dt.mu.Lock()
    defer dt.mu.Unlock()
    delete(dt.deployList, id)
}

func (dt *deployTask) stop(id int) {
    deploy, exists := dt.deployList[id]
    if exists {
        for _, d := range deploy {
            d.Terminate()
        }
    }
}
