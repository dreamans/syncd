// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "testing"
)

func TestDepoly(t *testing.T) {
    deploy := &Deploy{
        Srvs: []*Server{
            NewServer(1, 22, "127.0.0.1", "work", "~/.ssh/id_rsa"),
            NewServer(1, 22, "127.0.0.1", "work", "~/.ssh/id_rsa"),
        },
        PreCmd: "whoami",
        PostCmd: "whoami",
        DeployPath: "/home/work/deploy",
        PackFile: "/Users/work/project/syncd-deploy.tar.gz",
    }
    deploy, _ = NewDepoly(deploy)
    deploy.Deploy()
    wait := deploy.ParalDeploy()
    wait()
}

