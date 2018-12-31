// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

type Svn struct {
    repo    *Repo
}

func (s *Svn) SetRepo(r *Repo) {
    s.repo = r
}

func (s *Svn) UpdateRepo(branch string) (string, error) {

    return "", nil
}

func (s *Svn) ResetRepo() string {
    return ""
}

func (s *Svn) TagListRepo() string {
    return ""
}

func (s *Svn) CommitListRepo() string {
    return ""
}

