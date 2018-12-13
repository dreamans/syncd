// Copyright 2018 tinystack Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

import (
    "github.com/tinystack/syncd"
)

type QueryParam struct {
    Offset     int
    Limit      int
    Order      string
    Fields     string
    Plain      string
    Prepare    []interface{}
}

func Create(tableName string, data interface{}) bool {
    if err := syncd.Orm.Table(tableName).Create(data).Error; err != nil {
        syncd.Logger.Warning("mysql query error: %v", err)
        return false
    }
    return true
}

func GetMulti(tableName string, data interface{}, query QueryParam) bool {
    db := syncd.Orm.Table(tableName).Offset(query.Offset).Limit(query.Limit)
    if query.Fields != "" {
        db = db.Select(query.Fields)
    }
    if query.Order != "" {
        db = db.Order(query.Order)
    }
    db = db.Find(data)
    if err := db.Error; err != nil {
        syncd.Logger.Warning("mysql query error: %v", err)
        return false
    }

    return true
}

func Count(tableName string, count *int, query QueryParam) bool {
    db := syncd.Orm.Table(tableName).Count(&count)
    if err := db.Error; err != nil {
        syncd.Logger.Warning("mysql query error: %v", err)
        return false
    }
    return true
}

func GetOne(tableName string, data interface{}, query QueryParam) bool {
    db := syncd.Orm.Table(tableName)
    if query.Fields != "" {
        db = db.Select(query.Fields)
    }
    if query.Plain != "" {
        db = db.Where(query.Plain, query.Prepare...)
    }
    db = db.First(data)
    if err := db.Error; err != nil && !db.RecordNotFound() {
        syncd.Logger.Warning("mysql query error: %v", err)
        return false
    }
    return true
}

func Update(tableName string, data interface{}, query QueryParam) bool {
    db := syncd.Orm.Table(tableName)
    if query.Plain != "" {
        db = db.Where(query.Plain, query.Prepare...)
    }
    db = db.Updates(data)
    if err := db.Error; err != nil {
        syncd.Logger.Warning("mysql query error: %v", err)
        return false
    }
    return true
}
