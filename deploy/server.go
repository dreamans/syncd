// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "fmt"

    "github.com/dreamans/syncd/util/command"
)

const (
    COMMAND_TIMEOUT = 3600
)

type Server struct {
    ID          int
    Addr        string
    User        string
    Port        int
    Key         string
    task        *command.Task
    status      *ServerStatus
}

type ServerStatus struct {
    TaskResult  []*command.TaskResult
    Status      int
    Error       error
}

func NewServer(id, port int, addr, user, key string) *Server {
    return &Server{
        ID: id,
        Addr: addr,
        User: user,
        Port: port,
        Key: key,
        status: &ServerStatus{
            Status: STATUS_INIT,
        },
    }
}

func (srv *Server) Deploy(deploy *Deploy) {
    srv.status.Status = STATUS_RUNNING
    srv.task = command.TaskNew(
        srv.deployCmd(deploy),
        COMMAND_TIMEOUT,
    )
    srv.task.Run()
    if err := srv.task.GetError(); err != nil {
        srv.status.Status = STATUS_FAILED
        srv.status.Error = err
    } else {
        srv.status.Status = STATUS_DONE
    }
    srv.status.TaskResult = srv.task.Result()
}

func (srv *Server) Terminate() {
    if srv.status.Status == STATUS_RUNNING {
        srv.task.Terminate()
    }
}

func (srv *Server) Status() *ServerStatus {
    return srv.status
}

func (srv *Server) deployCmd(deploy *Deploy) []string {
    var (
        useCustomKey, useSshPort, useScpPort string
    )
    if srv.Key != "" {
        useCustomKey = fmt.Sprintf("-i %s", srv.Key)
    }
    if srv.Port != 0 {
        useSshPort = fmt.Sprintf("-p %d", srv.Port)
        useScpPort = fmt.Sprintf(" -P %d", srv.Port)
    }

    cmds := []string{
        fmt.Sprintf(
            "/usr/bin/env ssh -o StrictHostKeyChecking=no %s %s %s@%s 'mkdir -p %s; mkdir -p %s'",
            useCustomKey,
            useSshPort,
            srv.User,
            srv.Addr,
            deploy.DeployTmpPath,
            deploy.DeployPath,
        ),
        fmt.Sprintf(
            "/usr/bin/env scp -o StrictHostKeyChecking=no -q %s %s %s %s@%s:%s/",
            useCustomKey,
            useScpPort,
            deploy.PackFile,
            srv.User,
            srv.Addr,
            deploy.DeployTmpPath,
        ),
    }
    if deploy.PreCmd != "" {
        cmds = append(
            cmds,
            fmt.Sprintf(
                "/usr/bin/env ssh -o StrictHostKeyChecking=no %s %s %s@%s '%s'",
                useCustomKey,
                useSshPort,
                srv.User,
                srv.Addr,
                deploy.PreCmd,
            ),
        )
    }
    cmds = append(
        cmds,
        fmt.Sprintf(
            "/usr/bin/env ssh -o StrictHostKeyChecking=no %s %s %s@%s 'cd %s; tar -zxf %s -C %s; rm -f %s'",
            useCustomKey,
            useSshPort,
            srv.User,
            srv.Addr,
            deploy.DeployTmpPath,
            deploy.PackFileName,
            deploy.DeployPath,
            deploy.PackFileName,
        ),
    )
    if deploy.PostCmd != "" {
        cmds = append(
            cmds,
            fmt.Sprintf("/usr/bin/env ssh -o StrictHostKeyChecking=no %s %s %s@%s '%s'",
            useCustomKey,
            useSshPort,
            srv.User,
            srv.Addr,
            deploy.PostCmd,
        ),
    )
}
return cmds
}

