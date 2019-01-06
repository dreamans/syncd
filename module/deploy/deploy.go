// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
    "strings"

    "github.com/tinystack/goutil/gostring"
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    projectService "github.com/tinystack/syncd/service/project"
    deployService "github.com/tinystack/syncd/service/deploy"
    userService "github.com/tinystack/syncd/service/user"
    repoService "github.com/tinystack/syncd/service/repo"
    serverService "github.com/tinystack/syncd/service/server"
)

func DeployStart(c *goweb.Context) error {
    id := c.PostFormInt("id")
    apply, err := deployService.ApplyGetByPk(id)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    if err := deployCheckUserPriv(c, apply); err != nil {
        return err
    }
    if apply.Status != 3 {
        return syncd.RenderAppError("apply status wrong")
    }

    project, err := projectService.ProjectGetByPk(apply.ProjectId)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }

    repo := &repoService.Repo{
        ID: apply.ProjectId,
        ApplyId: apply.ID,
        Url: apply.RepoData.RepoUrl,
    }
    if repo, err = repoService.RepoNew(repo); err != nil {
        return syncd.RenderAppError(err.Error())
    }

    var deployTasks []*deployService.DeployTask
    if apply.RepoData.RepoMode == 1 {
        updateRepoCmd, err := repo.UpdateRepo(apply.RepoData.RepoBranch)
        if err != nil {
            return syncd.RenderAppError(err.Error())
        }
        deployTasks = append(deployTasks, &deployService.DeployTask{
            ApplyId: apply.ID,
            Level: deployService.DEPLOY_LEVEL_UPDATE_REPO,
            Cmd: updateRepoCmd,
        })

        tmpCommit := strings.Split(apply.RepoData.Commit, " - ")
        if len(tmpCommit) == 0 {
            return syncd.RenderAppError("git commit version  wrong")
        }
        commit := tmpCommit[0]
        update2CommitCmd := repo.Update2CommitRepo(apply.RepoData.RepoBranch, commit)
        deployTasks = append(deployTasks, &deployService.DeployTask{
            ApplyId: apply.ID,
            Level: deployService.DEPLOY_LEVEL_UPDATE_COMMIT,
            Cmd: update2CommitCmd,
        })
    } else {
        updateRepoCmd, err := repo.UpdateRepo("")
        if err != nil {
            return syncd.RenderAppError(err.Error())
        }
        deployTasks = append(deployTasks, &deployService.DeployTask{
            ApplyId: apply.ID,
            Level: deployService.DEPLOY_LEVEL_UPDATE_REPO,
            Cmd: updateRepoCmd,
        })
        update2CommitCmd := repo.Update2CommitRepo(apply.RepoData.Tag, "")
        deployTasks = append(deployTasks, &deployService.DeployTask{
            ApplyId: apply.ID,
            Level: deployService.DEPLOY_LEVEL_UPDATE_COMMIT,
            Cmd: update2CommitCmd,
        })
    }

    //tar zcvf 
    exFiles := gostring.StrFilterSliceEmpty(strings.Split(project.ExcludeFiles, "\n"))
    packRepoCmd := repo.PackRepo(exFiles)
    deployTasks = append(deployTasks, &deployService.DeployTask{
        ApplyId: apply.ID,
        Level: deployService.DEPLOY_LEVEL_PACK_REPO,
        Cmd: packRepoCmd,
    })

    srvList, err := serverService.ServerGetListByGroupIds(project.DeployServer)
    if err != nil {
        return syncd.RenderAppError(err.Error())
    }
    for _, srv := range srvList {
        deployCmds := repo.DeployRepo(gostring.Int2Str(srv.SshPort), srv.Ip, project.DeployUser, project.DeployPath, project.PreDeployCmd, project.PostDeployCmd)
        deployCmd := gostring.JsonEncode(deployCmds)
        deployTasks = append(deployTasks, &deployService.DeployTask{
            ApplyId: apply.ID,
            Level: deployService.DEPLOY_LEVEL_DEPLOY,
            Cmd: deployCmd,
        })
    }


    for _, dt := range deployTasks {
        if _, err := dt.Create(); err != nil {
            return syncd.RenderAppError(err.Error())
        }
    }

    return nil
}

func deployCheckUserPriv(c *goweb.Context, apply *deployService.Apply) error {
    if havePriv := userService.PrivIn(userService.DEPLOY_DEPLOY_ALL, c.GetIntSlice("priv")); !havePriv {
        if apply.UserId != c.GetInt("user_id") {
            return syncd.RenderCustomerError(syncd.CODE_ERR_NO_PRIV, "no priv")
        }
    }
    if err := applyCheckUserInSpace(apply.SpaceId, c.GetInt("user_id")); err != nil {
        return err
    }
    return nil
}
