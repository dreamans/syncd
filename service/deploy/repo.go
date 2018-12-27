// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "fmt"
    "time"
    "errors"

    "github.com/tinystack/syncd"
    "github.com/tinystack/goutil/gostring"
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
    UpdateRepo(branch string) error
    ResetRepo() error
    TagListRepo() ([]string, error)
    CommitListRepo() ([]string, error)
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
    r.LocalPath = gostring.JoinStrings(syncd.DataDir, "/", gostring.Int2Str(r.ID))
    r.fd.SetRepo(r)
    return r, nil
}

func (r *Repo) ResetRepo() error {
    return r.fd.ResetRepo()
}

func (r *Repo) TagListRepo() ([]string, error) {
    return r.fd.TagListRepo()
}

func (r *Repo) CommitListRepo() ([]string, error){
    return r.fd.CommitListRepo()
}

func (r *Repo) UpdateRepo(branch string) error {
    return r.fd.UpdateRepo(branch)
}

func (r *Repo) NewCommand(cmd string) *gocmd.Command {
    return &gocmd.Command{
        Cmd: cmd,
        Timeout: time.Second * r.CmdRunTimeout,
    }
}
