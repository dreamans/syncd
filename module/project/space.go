// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "github.com/tinystack/goweb"
    "github.com/tinystack/syncd"
    "github.com/tinystack/syncd/route"
    projectService "github.com/tinystack/syncd/service/project"
)

func init() {
    route.Register(route.API_PROJECT_SPACE_UPDATE, updateProjectSpace)
    route.Register(route.API_PROJECT_SPACE_LIST, listProjectSpace)
    route.Register(route.API_PROJECT_SPACE_DETAIL, detailProjectSpace)
    route.Register(route.API_PROJECT_SPACE_DELETE, deleteProjectSpace)
    route.Register(route.API_PROJECT_SPACE_EXISTS, existsProjectSpace)
}

func updateProjectSpace(c *goweb.Context) error {
    id, name := c.PostFormInt("id"), c.PostForm("name")
    if name == "" {
        syncd.RenderParamError(c, "name can not be empty")
        return nil
    }

    spaceExists := &projectService.Space{
        ID: id,
        Name: name,
    }
    exists, err := spaceExists.CheckSpaceExists()
    if err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    if exists {
        syncd.RenderAppError(c, "space data update failed, space name have exists")
        return nil
    }

    projectSpace := &projectService.Space{
        ID: id,
        Name: name,
        Description: c.PostForm("description"),
    }
    if err := projectSpace.CreateOrUpdate(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}

func listProjectSpace(c *goweb.Context) error {
    offset, limit := c.QueryInt("offset"), c.QueryInt("limit")
    keyword := c.Query("keyword")

    projectSpace := &projectService.Space{}
    list, total, err := projectSpace.List(keyword, offset, limit)
    if err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }

    //check if project exists in the space
    var newList []map[string]interface{}
    for _, l := range list {
        projService := &projectService.Project{
            SpaceId: l.ID,
        }
        exists, err := projService.CheckSpaceHaveProject()
        if err != nil {
            syncd.RenderAppError(c, err.Error())
            return nil
        }
        newList = append(newList, map[string]interface{}{
            "id": l.ID,
            "name": l.Name,
            "description": l.Description,
            "have_project": exists,
            "ctime": l.Ctime,
        })
    }
    syncd.RenderJson(c, goweb.JSON{
        "list": newList,
        "total": total,
    })
    return nil
}

func detailProjectSpace(c *goweb.Context) error {
    projectSpace := &projectService.Space{
        ID: c.QueryInt("id"),
    }
    if err := projectSpace.Get(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, projectSpace)
    return nil
}

func deleteProjectSpace(c *goweb.Context) error {
    var (
        id int
        exists bool
        err error
    )
    id = c.PostFormInt("id")
    proj := &projectService.Project{
        SpaceId: id,
    }
    exists, err = proj.CheckSpaceHaveProject()
    if err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    if exists {
        syncd.RenderAppError(c, "space delete failed, project in space is not empty")
        return nil
    }

    projectSpace := &projectService.Space{
        ID: id,
    }
    if err = projectSpace.Delete(); err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, nil)
    return nil
}

func existsProjectSpace(c *goweb.Context) error {
    keyword := c.Query("keyword")
    id := c.QueryInt("id")
    if keyword == "" {
        syncd.RenderParamError(c, "params error")
        return nil
    }
    projectSpace := &projectService.Space{
        ID: id,
        Name: keyword,
    }
    exists, err := projectSpace.CheckSpaceExists()
    if err != nil {
        syncd.RenderAppError(c, err.Error())
        return nil
    }
    syncd.RenderJson(c, goweb.JSON{
        "exists": exists,
    })
    return nil
}
