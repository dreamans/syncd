// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package integrate

import (
    "fmt"
    "errors"
)

type Repo struct {
    Url     string
    Branch  string
    Commit  string
    Local   string
}

func NewRepo(url, branch, commit, local string) (*Repo, error) {
    if url == "" {
        return nil, errors.New("repo url can not empty")
    }
    if local == "" {
        return nil, errors.New("repo local path can not empty")
    }
    repo := &Repo{
        Url: url,
        Branch: branch,
        Commit: commit,
        Local: local,
    }

    return repo, nil
}

func (repo *Repo) FetchCmd() []string {
    cmds := []string{
        fmt.Sprintf("rm -fr %s", repo.Local),
        fmt.Sprintf("/usr/bin/env git clone -q %s %s", repo.Url, repo.Local),
        fmt.Sprintf("/usr/bin/env git checkout -q %s", repo.Branch),
    }
    if repo.Commit != "" {
        cmds = append(cmds, fmt.Sprintf("/usr/bin/env git reset -q --hard %s", repo.Commit))
    }
    return cmds
}

