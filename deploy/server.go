// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "fmt"
    "path"

    "github.com/dreamans/syncd/util/command"
)

const (
    COMMAND_TIMEOUT = 3600
)

type Server struct {
    ID              int
    Addr            string
    User            string
    Port            int
    PreCmd          string
    PostCmd         string
    Key             string
    PackFile        string
    DeployTmpPath   string
    DeployPath      string
    task            *command.Task
    result          *ServerResult
}

type ServerResult struct {
    ID          int
    TaskResult  []*command.TaskResult
    Status      int
    Error       error
}

func NewServer(srv *Server) {
    srv.result = &ServerResult{
        ID: srv.ID,
        Status: STATUS_INIT,
    }
    srv.task = command.NewTask(
        srv.deployCmd(),
        COMMAND_TIMEOUT,
    )
}

func (srv *Server) Deploy() {
    srv.result.Status = STATUS_ING
    srv.task.Run()
    if err := srv.task.GetError(); err != nil {
        srv.result.Error = err
        srv.result.Status = STATUS_FAILED
    } else {
        srv.result.Status = STATUS_DONE
    }
}

func (srv *Server) Terminate() {
    if srv.result.Status == STATUS_ING {
        srv.task.Terminate()
    }
}

func (srv *Server) Result() *ServerResult {
    srv.result.TaskResult = srv.task.Result()
    return srv.result
}

func (srv *Server) deployCmd() []string {
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
    var cmds []string
    if srv.PackFile == "" {
        cmds = append(cmds, "echo 'packfile empty' && exit 1")
    }

    cmds = append(cmds, []string{
        fmt.Sprintf(
            "/usr/bin/env ssh -o StrictHostKeyChecking=no %s %s %s@%s 'mkdir -p %s; mkdir -p %s'",
            useCustomKey,
            useSshPort,
            srv.User,
            srv.Addr,
            srv.DeployTmpPath,
            srv.DeployPath,
        ),
        fmt.Sprintf(
            "/usr/bin/env scp -o StrictHostKeyChecking=no -q %s %s %s %s@%s:%s/",
            useCustomKey,
            useScpPort,
            srv.PackFile,
            srv.User,
            srv.Addr,
            srv.DeployTmpPath,
        ),
    }...)
    if srv.PreCmd != "" {
        cmds = append(
            cmds,
            fmt.Sprintf(
                "/usr/bin/env ssh -o StrictHostKeyChecking=no %s %s %s@%s '%s'",
                useCustomKey,
                useSshPort,
                srv.User,
                srv.Addr,
                srv.PreCmd,
            ),
        )
    }
    packFileName := path.Base(srv.PackFile)
    cmds = append(
        cmds,
        fmt.Sprintf(
            "/usr/bin/env ssh -o StrictHostKeyChecking=no %s %s %s@%s 'cd %s; tar -zxf %s -C %s; rm -f %s'",
            useCustomKey,
            useSshPort,
            srv.User,
            srv.Addr,
            srv.DeployTmpPath,
            packFileName,
            srv.DeployPath,
            packFileName,
        ),
    )
    if srv.PostCmd != "" {
        cmds = append(
            cmds,
            fmt.Sprintf("/usr/bin/env ssh -o StrictHostKeyChecking=no %s %s %s@%s '%s'",
                useCustomKey,
                useSshPort,
                srv.User,
                srv.Addr,
                srv.PostCmd,
            ),
        )
    }
    return cmds
}