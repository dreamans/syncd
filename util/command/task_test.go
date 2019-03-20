// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package command

import (
    "testing"
)

func TestTaskRun(t *testing.T) {
    cmds := []string{
        "echo 'syncd'",
        "whoami",
        "date",
    }
    task := TaskNew(cmds, 10)
    task.Run()
    if err := task.GetError(); err != nil {
        t.Errorf("cmd task running error: %s", err.Error())
    }
}