// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
    "errors"
    "fmt"

    "github.com/tinystack/syncd"
    "github.com/tinystack/goutil/gostring"
)

type Repo struct {
    ID              int
    Repo            string
    Url             string
    User            string
    Pass            string
    localPath       string
    fd              Repository
}

type Repository interface {

    SetRepo(r *Repo)

    UpdateRepo(branch string) (string, error)

    ResetRepo() string

    TagListRepo() string

    CommitListRepo() string
}

func RepoNew(r *Repo) (*Repo, error) {
    switch r.Repo {
    case "git":
        r.fd = &Git{}
    case "svn":
        r.fd = &Svn{}
    default:
        return nil, errors.New(fmt.Sprintf("repository type error, want '%s', but '%s'", "git or svn", r.Repo))
    }
    r.localPath = gostring.JoinStrings(syncd.DataDir, "/", gostring.Int2Str(r.ID))
    r.fd.SetRepo(r)
    return r, nil
}

func (r *Repo) ResetRepo() string {
    return r.fd.ResetRepo()
}

func (r *Repo) TagListRepo() string {
    return r.fd.TagListRepo()
}

func (r *Repo) CommitListRepo() string {
    return r.fd.CommitListRepo()
}

func (r *Repo) UpdateRepo(branch string) (string, error) {
    return r.fd.UpdateRepo(branch)
}

