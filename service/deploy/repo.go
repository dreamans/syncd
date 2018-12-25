// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "fmt"
    "time"
    "errors"
    "strings"

    "github.com/tinystack/syncd"
    "github.com/tinystack/goutil"
    "github.com/tinystack/gocmd"
)

type Repo struct {
    ID              int
    LocalPath       string
    Repo            string
    Url             string
    User            string
    Pass            string
    CmdRunTimeout   time.Duration
    fd              Repository
}

type Repository interface {
    SetRepo(r *Repo)
    UpdateRepoCmd(branch string) (string, error)
    ResetRepoCmd() string
    TagListCmd() string
}

func NewRepo(r *Repo) (*Repo, error) {
    switch r.Repo {
    case "git":
        r.fd = &Git{}
    case "svn":
        r.fd = &Svn{}
    default:
        return nil, errors.New(fmt.Sprintf("repository type error, want '%s', but '%s'", "git or svn", r.Repo))
    }
    r.LocalPath = goutil.JoinStrings(syncd.DataDir, "/", goutil.Int2Str(r.ID))
    r.fd.SetRepo(r)
    return r, nil
}

func (r *Repo) ResetRepo() error {
    resetCmd := r.fd.ResetRepoCmd()
    return r.newCommand(resetCmd).Run()
}

func (r *Repo) TagListRepo() ([]string, error) {
    tagListCmd := r.fd.TagListCmd()
    println(tagListCmd)
    cmd := r.newCommand(tagListCmd)
    if err := cmd.Run(); err != nil {
        return nil, errors.New(err.Error() + ", " + string(cmd.Stderr()))
    }
    tagList := strings.Split(string(cmd.Stdout()), "\n")
    tagList = goutil.StringSliceRsort(tagList)

    return tagList, nil
}

func (r *Repo) UpdateRepo(branch string) error {
    updateCmd, err := r.fd.UpdateRepoCmd(branch)
    if err != nil {
        return nil
    }
    return r.newCommand(updateCmd).Run()
}

func (r *Repo) newCommand(cmd string) *gocmd.Command {
    return &gocmd.Command{
        Cmd: cmd,
        Timeout: time.Second * r.CmdRunTimeout,
    }
}
