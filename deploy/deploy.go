// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "sync"
)

const (
    STATUS_INIT = 0
    STATUS_ING = 1
    STATUS_DONE = 2
    STATUS_FAILED = 3
)

const (
    DEPLOY_PARALLEL = 1
    DEPLOY_SERIAL = 2
)

type Deploy struct {
    ID              int
    User            string
    PreCmd          string
    PostCmd         string
    DeployPath      string
    DeployTmpPath   string
    PackFile        string
    srvs            []*Server
    status          int
    wg              sync.WaitGroup
}

type DeployResult struct {
    ID          int
    Status      int
    ServerRest  []*ServerResult
}

func (d *Deploy) AddServer(id int, addr string, port int) {
    srv := &Server{
        ID: id,
        Addr: addr,
        User: d.User,
        Port: port,
        PreCmd: d.PreCmd,
        PostCmd: d.PostCmd,
        PackFile: d.PackFile,
        DeployTmpPath: d.DeployTmpPath,
        DeployPath: d.DeployPath,
    }
    NewServer(srv)
    d.srvs = append(d.srvs, srv)
}

func (d *Deploy) Parallel() {
    if d.status == STATUS_FAILED {
        return
    }
    d.status = STATUS_ING
    status := STATUS_DONE
    for _, srv := range d.srvs {
        d.wg.Add(1)
        go func() {
            if d.status == STATUS_ING {
                srv.Deploy()
                if srv.Result().Status == STATUS_FAILED {
                    status = STATUS_FAILED
                }
            }
            defer d.wg.Done()
        }()
    }
    d.wg.Wait()
    d.status = status
}

func (d *Deploy) Serial() {
    if d.status == STATUS_FAILED {
        return
    }
    d.status = STATUS_ING
    status := STATUS_DONE
    for _, srv := range d.srvs {
        if d.status == STATUS_ING {
            srv.Deploy()
            if srv.Result().Status == STATUS_FAILED {
                status = STATUS_FAILED
            }
        }
    }
    d.status = status
}

func (d *Deploy) Result() ([]*ServerResult, int) {
    var rest []*ServerResult
    for _, srv := range d.srvs {
        rest = append(rest, srv.Result())
    }
    return rest, d.status
}

func (d *Deploy) Terminate() {
    d.status = STATUS_FAILED
    for _, srv := range d.srvs {
        srv.Terminate()
    }
}