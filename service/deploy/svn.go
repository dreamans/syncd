// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
)

type Svn struct {
    repo    *Repo
}

func (s *Svn) SetRepo(r *Repo) {
    s.repo = r
}

func (s *Svn) UpdateRepo(branch string) error {

    return nil
}

func (s *Svn) ResetRepo() error {
    return nil
}

func (s *Svn) TagListRepo() ([]string, error) {
    return nil, nil
}

func (s *Svn) CommitListRepo() ([]string, error) {
    return nil, nil
}
