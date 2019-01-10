// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
    "fmt"

    "github.com/tinystack/goutil/gopath"
    "github.com/tinystack/goutil/gostring"
)

type Git struct {
    repo    *Repo
}

func (g *Git) SetRepo(r *Repo) {
    g.repo = r
}

func (g *Git) UpdateRepo(branch string) (string, error) {
    if branch == "" {
        branch = "master"
    }
    exists, err := gopath.PathExists(g.repo.localPath + "/.git")
    if err != nil {
        return "", err
    }
    var cmd string
    if exists {
        cmd = gostring.JoinSepStrings(
            " && ",
            fmt.Sprintf("cd %s", g.repo.localPath),
            fmt.Sprintf("/usr/bin/env git checkout -q %s", branch),
            "/usr/bin/env git fetch -p -q --all",
            fmt.Sprintf("/usr/bin/env git reset -q --hard origin/%s", branch),
        )
    } else {
        cmd = gostring.JoinSepStrings(
            " && ",
            fmt.Sprintf("mkdir -p %s", g.repo.localPath),
            fmt.Sprintf("cd %s", g.repo.localPath),
            fmt.Sprintf("/usr/bin/env git clone -q %s .", g.repo.Url),
            fmt.Sprintf("/usr/bin/env git checkout -q %s", branch),
        )
    }
    return cmd, nil
}

func (g *Git) ResetRepo() string {
    cmd := gostring.JoinSepStrings(
        " && ",
        fmt.Sprintf("rm -rf %s", g.repo.localPath),
        fmt.Sprintf("mkdir -p %s", g.repo.localPath),
        fmt.Sprintf("cd %s", g.repo.localPath),
        fmt.Sprintf("git clone -q %s .", g.repo.Url),
    )
    return cmd
}

func (g *Git) TagListRepo() string {
    cmd := fmt.Sprintf("cd %s && git tag -l", g.repo.localPath)
    return cmd
}

func (g *Git) CommitListRepo() string {
    cmd := fmt.Sprintf("cd %s && git log -100 --pretty=format:\"%%h - %%ad - %%an %%s \" --date=format:\"%%Y-%%m-%%d %%H:%%M:%%S\"", g.repo.localPath)
    return cmd
}

func (g *Git) Update2CommitRepo(branch, commit string) string {
    cmd := []string{
        fmt.Sprintf("cd %s", g.repo.localPath),
    }
    if branch != "" {
        cmd = append(cmd, fmt.Sprintf("/usr/bin/env git checkout -q %s", branch))
    }
    if commit != "" {
        cmd = append(cmd, fmt.Sprintf("/usr/bin/env git reset -q --hard %s", commit))
    }
    return gostring.JoinSepStrings(" && ", cmd...)
}

