// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package build

import (
    "strings"
    "os/user"

    "github.com/tinystack/syncd/pkg/util"
)

type Bash struct {
    bashTpl         []string
    gitCmd          string
    buildId         string
    rootPath        string
    buildPath       string
    repoPath        string
    customBuildFile string
}

const (
    BashRootPath = "~/.syncd"
    BashPkgFile = "package.tgz"
    BashRepoBasePath = "repo"
    BashCustomBuildFileName = ".syncd.build.sh"
)

func NewBash(buildId, gitCmd string) *Bash {
    bash := &Bash{
        gitCmd: gitCmd,
        buildId: buildId,
        rootPath: "~/.syncd",
    }

    if user, err := user.Current(); err == nil {
        bash.rootPath = util.JoinStrings(user.HomeDir, "/.syncd")
    }

    bash.buildPath = util.JoinStrings(bash.rootPath, "/", bash.buildId)
    bash.repoPath = util.JoinStrings(bash.buildPath, "/", BashRepoBasePath)
    bash.customBuildFile = util.JoinStrings(bash.repoPath, "/", BashCustomBuildFileName)
    return bash
}

func (bash *Bash) BashCmd() string {
    bash.appendCmd(
        "#!/bin/bash",
        bash.outputBash(">>> building"),
        bash.outputBash(util.JoinStrings("start checking ", bash.buildPath)),
    )
    bash.appendCmd(
        bash.ifBash(
            util.JoinStrings("-d ", "\"", bash.buildPath, "\""),
            []string{
                util.JoinStrings("rm -fr ", bash.buildPath),
                bash.outputBash(util.JoinStrings("delete old build path ", bash.buildPath, " success")),
            },
        )...
    )
    bash.appendCmd(
        util.JoinStrings("mkdir -p ", bash.buildPath),
    )
    bash.appendCmd(
        bash.ifErrorBash("build path create failed")...
    )
    bash.appendCmd(
        bash.outputBash(util.JoinStrings("new ", bash.buildPath, " create success")),
    )
    bash.appendCmd(
        bash.outputBash(util.JoinStrings("start clone remote repo to ", bash.repoPath)),
    )
    bash.appendCmd(
        util.JoinStrings(bash.gitCmd, " ",bash.repoPath),
    )
    bash.appendCmd(bash.ifErrorBash("git clone repo failed, script exit")...)
    bash.appendCmd(
        bash.outputBash(util.JoinStrings("start run repo ", BashCustomBuildFileName)),
    )
    bash.appendCmd(
        bash.ifBash(
            util.JoinStrings("-f ", "\"", bash.customBuildFile, "\""),
            append([]string{
                    util.JoinStrings("/bin/bash ", bash.customBuildFile),
                },
                bash.ifErrorBash(util.JoinStrings(BashCustomBuildFileName, " occur error and exit"))...
            ),
        )...
    )
    bash.appendCmd(
        bash.ifBash(
            util.JoinStrings("! -f ", "\"", bash.customBuildFile, "\""),
            []string{
                bash.outputBash(util.JoinStrings(BashCustomBuildFileName, " get lost")),
            },
        )...
    )
    bash.appendCmd(
        bash.outputBash("start pack repo files"),
        util.JoinStrings("cd ", bash.buildPath),
    )
    bash.appendCmd(
        bash.ifBash(
            util.JoinStrings("-f ", BashPkgFile),
            []string{
                util.JoinStrings("rm -f ", BashPkgFile),
            },
        )...
    )
    bash.appendCmd(
        util.JoinStrings("cd ", bash.repoPath),
        util.JoinStrings("tar zcf ../", BashPkgFile, " *"),
    )
    bash.appendCmd(
        bash.ifErrorBash("repo files pack failed and exit")...
    )
    bash.appendCmd(
        bash.outputBash(util.JoinStrings("packed repo files to ", bash.repoPath)),
        bash.outputBash("<<< building finish"),
    )

    return bash.combineCmd()
}

func (bash *Bash) appendCmd(cmd ...string) {
    bash.bashTpl = append(bash.bashTpl, cmd...)
}

func (bash *Bash) combineCmd() string {
    return strings.Join(bash.bashTpl, "\n")
}

func (bash *Bash) outputBash(msg string) string {
    return util.JoinStrings(`echo "["$(date "+%F %T")"]" "[syncd]" `, "\"", msg, "\"")
}

func (bash *Bash) ifBash(condition string, cmd []string) []string {
    ifCmd := []string{
        util.JoinStrings("if [ ", condition, " ];then"),
    }
    ifCmd = append(ifCmd, cmd...)
    ifCmd = append(ifCmd, "fi")
    return ifCmd
}

func (bash *Bash) ifErrorBash(errMsg string) []string {
    return []string{
        "if [ $? -ne 0 ]; then",
        bash.outputBash(util.JoinStrings("[failed] ", errMsg)),
        "exit 127",
        "fi",
    }
}

