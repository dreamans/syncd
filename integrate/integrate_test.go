// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package integrate

import (
    "testing"
)

func TestFetch(t *testing.T) {
    workSpace := "/tmp/laravel_test"
    repo, _ := NewRepo("https://gitee.com/dreamans/syncd.git", "1.1.0", "", workSpace)
    build := NewBuild(workSpace, workSpace + "/" + "laravel.tar.gz", repo)
    result, err := build.Fetch()
    if err != nil {
        t.Errorf("fetch run failed, err is %s", err.Error())
    }
    for _, r := range result {
        t.Logf("cmd[%s], stdout[%s], stderr[%s], success[%v]\n", r.Cmd, r.Stdout, r.Stderr, r.Success)
    }
}
