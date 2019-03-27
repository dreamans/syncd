// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package common

import (
	"fmt"

	"github.com/dreamans/syncd"
	"github.com/dreamans/syncd/util/gofile"
	"github.com/dreamans/syncd/util/command"
	"github.com/dreamans/syncd/util/gostring"
	"github.com/dreamans/syncd/module/project"
	"github.com/dreamans/syncd/module/deploy"
)

func HookBuild(applyId int) {
	apply := &deploy.Apply{
		ID: applyId,
	}
	if err := apply.Detail(); err != nil {
		return
	}
	proj := &project.Project{
		ID: apply.ProjectId,
	}
	if err := proj.Detail(); err != nil {
		return
	}
	if proj.BuildHookScript == "" {
		return
	}

	build := &deploy.Build{
		ApplyId: applyId,
	}
	if err := build.Detail(); err != nil {
		return
	}

	var buildStatus int
	if build.Status == deploy.BUILD_STATUS_SUCCESS {
		buildStatus = 1
	}

	s := gostring.JoinStrings(
        "#!/bin/bash\n\n",
        "#--------- build hook scripts env ---------\n",
        fmt.Sprintf("env_apply_id=%d\n", apply.ID),
		fmt.Sprintf("env_apply_name='%s'\n", apply.Name),
		fmt.Sprintf("env_pack_file='%s'\n", build.Tar),
		fmt.Sprintf("env_build_output='%s'\n", build.Output),
		fmt.Sprintf("env_build_errmsg='%s'\n", build.Errmsg),
		fmt.Sprintf("env_build_status=%d\n", buildStatus),
        proj.BuildHookScript,
	)
	hookCommandTaskRun(s, applyId)
}

func HookDeploy(applyId int) {
	apply := &deploy.Apply{
		ID: applyId,
	}
	if err := apply.Detail(); err != nil {
		return
	}
	proj := &project.Project{
		ID: apply.ProjectId,
	}
	if err := proj.Detail(); err != nil {
		return
	}
	if proj.DeployHookScript == "" {
		return
	}

	var deployStatus int
	if apply.Status == deploy.APPLY_STATUS_SUCCESS {
		deployStatus = 1
	}

	s := gostring.JoinStrings(
        "#!/bin/bash\n\n",
        "#--------- deploy hook scripts env ---------\n",
        fmt.Sprintf("env_apply_id=%d\n", apply.ID),
		fmt.Sprintf("env_apply_name='%s'\n", apply.Name),
		fmt.Sprintf("env_deploy_status=%d\n", deployStatus),
        proj.DeployHookScript,
	)
	hookCommandTaskRun(s, applyId)
}

func hookCommandTaskRun(s string, applyId int) {
	scriptFile := gostring.JoinStrings(syncd.App.LocalTmpSpace, "/", gostring.StrRandom(24), ".sh")
	if err := gofile.CreateFile(scriptFile, []byte(s), 0744); err != nil {
		syncd.App.Logger.Error("hook script file create failed, err[%s], apply_id[%d]", err.Error(), applyId)
        return
	}
	cmds := []string{
		fmt.Sprintf("/bin/bash -c %s", scriptFile),
		fmt.Sprintf("rm -f %s", scriptFile),
	}
	task := command.NewTask(cmds, 86400)
	task.Run()
	if err := task.GetError(); err != nil {
		syncd.App.Logger.Error("hook script run failed, err[%s], output[%s], apply_id[%d]", err.Error(), gostring.JsonEncode(task.Result()), applyId)
    }
}