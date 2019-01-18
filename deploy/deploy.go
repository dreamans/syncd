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
    Srvs            []*Server
    PreCmd          string
    PostCmd         string
    DeployPath      string
    DeployTmpPath   string
    PackFile        string
    PackFileName    string
    status          int
    wg              sync.WaitGroup
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

func (deploy *Deploy) Deploy() {
    deploy.status = STATUS_RUNNING
    for _, srv := range deploy.Srvs {
        srv.Deploy(deploy)
    }
    deploy.status = STATUS_DONE
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

