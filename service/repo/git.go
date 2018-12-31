// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
    "net/url"
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
    exists, err := gopath.PathExists(g.repo.localPath)
    if err != nil {
        return "", err
    }
    var cmd string
    if exists {
        cmd = gostring.JoinSepStrings(
            " && ",
            fmt.Sprintf("cd %s", g.repo.localPath),
            fmt.Sprintf("git checkout -q %s", branch),
            "git fetch -p -q --all",
            fmt.Sprintf("git reset -q --hard origin/%s", branch),
        )
    } else {
        cmd = gostring.JoinSepStrings(
            " && ",
            fmt.Sprintf("mkdir %s", g.repo.localPath),
            fmt.Sprintf("cd %s", g.repo.localPath),
            fmt.Sprintf("git clone -q %s -b %s .", branch, g.getRemoteUrl()),
        )
    }
    return cmd, nil
}

func (g *Git) ResetRepo() string {
    cmd := gostring.JoinSepStrings(
        " && ",
        fmt.Sprintf("rm -rf %s", g.repo.localPath),
        fmt.Sprintf("mkdir %s", g.repo.localPath),
        fmt.Sprintf("cd %s", g.repo.localPath),
        fmt.Sprintf("git clone -q %s .", g.getRemoteUrl()),
    )
    return cmd
}

func (g *Git) TagListRepo() string {
    cmd := fmt.Sprintf("cd %s && git tag -l", g.repo.localPath)
    return cmd
}

func (g *Git) CommitListRepo() string {
    cmd := fmt.Sprintf("cd %s && git log -100 --pretty=format:\"%%h - %%ad - %%an %%s \" --date=format:\"%%Y-%%m-  %%d %%H:%%M:%%S\"", g.repo.localPath)
    return cmd
}

func (g *Git) getRemoteUrl() string {
    u, err := url.Parse(g.repo.Url)
    if err != nil {
        return g.repo.Url
    }
    var remoteUrl string
    if u.Scheme == "http" || u.Scheme == "https" {
        repoUrl := &url.URL{
            Scheme: u.Scheme,
            User: url.UserPassword(g.repo.User, g.repo.Pass),
            Host: u.Host,
            Path: u.Path,
            RawQuery: u.RawQuery,
        }
        remoteUrl = repoUrl.String()
    } else {
        return g.repo.Url
    }
    return remoteUrl
}

