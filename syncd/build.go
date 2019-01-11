// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// +build dev

package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/exec"
)

var (
    goos, goarch, goarm string
    race bool
)

func init() {
    flag.StringVar(&goos, "goos", "", "GOOS for which to build")
    flag.StringVar(&goarch, "goarch", "", "GOARCH for which to build")
    flag.StringVar(&goarm, "goarm", "", "GOARM for which to build")
    flag.BoolVar(&race, "race", false, "Enable race detector")
}

func main() {
    flag.Parse()

    gopath := os.Getenv("GOPATH")
    args := []string{
        "build",
        "-asmflags", fmt.Sprintf("-trimpath=%s", gopath),
        "-gcflags", fmt.Sprintf("-trimpath=%s", gopath),
    }
    if race {
        args = append(args, "-race")
    }

    env := os.Environ()
    env = append(env, "GOOS=" + goos, "GOARCH=" + goarch, "GOARM=" + goarm)
    if !race {
        env = append(env , "CGO_ENABLED=0")
    }

    cmd := exec.Command("go", args...)
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout
    cmd.Env = env
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
}

