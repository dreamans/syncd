// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
    "fmt"

    "github.com/dreamans/syncd"
    "github.com/tinystack/goutil/gostring"
)

type Repo struct {
    ID              int
    ApplyId         int
    Url             string
    localPath       string
    fd              Repository
}

type Repository interface {

    SetRepo(r *Repo)

    UpdateRepo(branch string) (string, error)

    ResetRepo() string

    TagListRepo() string

    CommitListRepo() string

    Update2CommitRepo(branch, commit string) string
}

func RepoNew(r *Repo) (*Repo, error) {
    r.fd = &Git{}
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

func (r *Repo) Update2CommitRepo(branch, commit string) string {
    return r.fd.Update2CommitRepo(branch, commit)
}

func (r *Repo) PackRepo(exFiles []string) string {
    var excludeCmds []string
    if len(exFiles) > 0 {
        for _, f := range exFiles {
            excludeCmds = append(excludeCmds, fmt.Sprintf("--exclude='%s'", f))
        }
    }
    tarFile := r.packFilePath()
    cmd := gostring.JoinSepStrings(
        " && ",
        fmt.Sprintf("rm -f %s", tarFile),
        fmt.Sprintf("cd %s", r.localPath),
        fmt.Sprintf("tar %s -zcvf %s *", gostring.JoinSepStrings(" ", excludeCmds...), tarFile),
    )
    return cmd
}

func (r *Repo) DeployRepo(sshPort, sshIp, sshUser, deployPath, preCmd, postCmd string) []string {
    var cmd []string
    cmd = append(cmd, fmt.Sprintf("/usr/bin/env ssh -p %s %s@%s 'mkdir -p %s; mkdir -p %s'", sshPort, sshUser, sshIp, syncd.RemoteTmpDir, deployPath))
    cmd = append(cmd, fmt.Sprintf("/usr/bin/env scp -P %s %s %s@%s:%s/", sshPort, r.packFilePath(), sshUser, sshIp, syncd.RemoteTmpDir))
    if preCmd != "" {
        cmd = append(cmd, fmt.Sprintf("/usr/bin/env ssh -p %s %s@%s '%s'", sshPort, sshUser, sshIp, preCmd))
    }
    cmd = append(cmd, fmt.Sprintf("/usr/bin/env ssh -p %s %s@%s 'cd %s; tar -zxf %s -C %s; rm -f %s'", sshPort, sshUser, sshIp, syncd.RemoteTmpDir, r.packFileName(), deployPath, r.packFileName()))
    if postCmd != "" {
        cmd = append(cmd, fmt.Sprintf("/usr/bin/env ssh -p %s %s@%s '%s'", sshPort, sshUser, sshIp, postCmd))
    }
    return cmd
}

func (r *Repo) packFilePath() string {
    return gostring.JoinStrings(syncd.TmpDir, "/", r.packFileName())
}

func (r *Repo) packFileName() string {
    return gostring.JoinStrings("syncd_", gostring.Int2Str(r.ID), "_", gostring.Int2Str(r.ApplyId), ".tar.gz")
}
