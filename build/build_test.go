// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package build

import (
    "testing"
    "time"
    "fmt"

    "github.com/dreamans/syncd/util/command"
)

func TestBuild(t *testing.T) {
    local := "/tmp/laravel"
    tmp := "/tmp"
    packFile := "/tmp/laravel.tar.gz"
    repo := NewRepo("git@gitee.com:mirrors/laravel.git", local)
    scripts := `
cd ${env_workspace}
tar --exclude=.git --exclude=.gitignore -zcvf ${env_pack_file} *
`
    build, err := NewBuild(repo, local, tmp, packFile, scripts)
    if err != nil {
        t.Errorf("create build task failed: %s", err.Error())
    }
    NewTask(1, build, func(id int, result *Result, taskResult []*command.TaskResult) {
        fmt.Println(id)
        fmt.Println(result)
        for _, t := range taskResult {
            fmt.Println(t)
        }
    })

    time.Sleep(2 * time.Second)

    result, taskResult, err := StatusTask(1)

    fmt.Println(result)
    for _, t := range taskResult {
        fmt.Println(t)
    }

    time.Sleep(10 * time.Second)

}