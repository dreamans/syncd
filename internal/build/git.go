// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package build

import (
    "net/url"
    "strings"
)

type GitCommand struct {
    remoteRepo  string
    user        string
    pass        string
}

func NewGitCommand(remoteRepo, user, pass string) *GitCommand {
    git := &GitCommand{
        remoteRepo: remoteRepo,
        user: user,
        pass: pass,
    }
    return git
}

func (git *GitCommand) CloneCmd(branch string) string {
    var cmd []string

    cmd = append(cmd, "git clone")
    cmd = append(cmd, "-v")
    cmd = append(cmd, "-b")
    cmd = append(cmd, branch)
    cmd = append(cmd, git.getRemoteUrl())

    return strings.Join(cmd, " ")
}

func (git *GitCommand) getRemoteUrl() string {
    u, err := url.Parse(git.remoteRepo)
    if err != nil {
        return git.remoteRepo
    }
    var remoteUrl string
    if u.Scheme == "http" || u.Scheme == "https" {
        repoUrl := &url.URL{
            Scheme: u.Scheme,
            User: url.UserPassword(git.user, git.pass),
            Host: u.Host,
            Path: u.Path,
            RawQuery: u.RawQuery,
        }
        remoteUrl = repoUrl.String()
    } else {
        return git.remoteRepo
    }
    return remoteUrl
}

