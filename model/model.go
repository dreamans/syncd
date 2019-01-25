// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

import (
    "fmt"
    "strings"

    "github.com/jinzhu/gorm"
    "github.com/dreamans/syncd"
)

type WhereParam struct {
    Field   string
    Tag     string
    Prepare interface{}
}

type QueryParam struct {
    Fields     string
    Offset     int
    Limit      int
    Order      string
    Where      []WhereParam
}

func Create(model interface{}) bool {
    db := syncd.App.DB.DbHandler.Create(model)
    if err := db.Error; err != nil {
        syncd.App.Logger.Warning("mysql execute error: %s, sql [%v]", err.Error(), db.QueryExpr())
        return false
    }
    return true
}

func GetMulti(model interface{}, query QueryParam) bool {
    db := syncd.App.DB.DbHandler.Offset(query.Offset)
    if query.Limit > 0 {
        db = db.Limit(query.Limit)
    }
    if query.Fields != "" {
        db = db.Select(query.Fields)
    }
    if query.Order != "" {
        db = db.Order(query.Order)
    }
    db = parseWhereParam(db, query.Where)
    db.Find(model)
    if err := db.Error; err != nil {
        syncd.App.Logger.Warning("mysql query error: %s, sql[%v]", err.Error(), db.QueryExpr())
        return false
    }
    return true
}

func Count(model interface{}, count *int, query QueryParam) bool {
    db := syncd.App.DB.DbHandler.Model(model)
    db = parseWhereParam(db, query.Where)
    db = db.Count(count)
    if err := db.Error; err != nil {
        syncd.App.Logger.Warning("mysql query error: %s, sql[%v]", err.Error(), db.QueryExpr())
        return false
    }
    return true
}

func Delete(model interface{}, query QueryParam) bool {
    if len(query.Where) == 0 {
        syncd.App.Logger.Warning("mysql query error: delete failed, where conditions cannot be empty")
        return false
    }
    db := syncd.App.DB.DbHandler.Model(model)
    db = parseWhereParam(db, query.Where)
    db = db.Delete(model)
    if err := db.Error; err != nil {
        syncd.App.Logger.Warning("mysql query error: %s, sql[%v]", err.Error(), db.QueryExpr())
        return false
    }
    return true
}

func DeleteByPk(model interface{}) bool {
    db := syncd.App.DB.DbHandler.Model(model)
    db.Delete(model)
    if err := db.Error; err != nil {
        syncd.App.Logger.Warning("mysql query error: %s, sql[%v]", err.Error(), db.QueryExpr())
        return false
    }
    return true
}

func GetOne(model interface{}, query QueryParam) bool {
    db := syncd.App.DB.DbHandler.Model(model)
    if query.Fields != "" {
        db = db.Select(query.Fields)
    }
    db = parseWhereParam(db, query.Where)
    db = db.First(model)
    if err := db.Error; err != nil && !db.RecordNotFound() {
        syncd.App.Logger.Warning("mysql query error: %s, sql[%v]", err.Error(), db.QueryExpr())
        return false
    }
    return true
}

func GetByPk(model interface{}, id interface{}) bool {
    db := syncd.App.DB.DbHandler.Model(model)
    db.First(model, id)
    if err := db.Error; err != nil && !db.RecordNotFound() {
        syncd.App.Logger.Warning("mysql query error: %s sql[%v]", err.Error(), db.QueryExpr())
        return false
    }
    return true
}

func UpdateByPk(model interface{}) bool {
    db := syncd.App.DB.DbHandler.Model(model)
    db = db.Updates(model)
    if err := db.Error; err != nil {
        syncd.App.Logger.Warning("mysql query error: %s, sql[%v]", err.Error(), db.QueryExpr())
        return false
    }
    return true
}

func Update(model interface{}, data interface{}, query QueryParam) bool {
    db := syncd.App.DB.DbHandler.Model(model)
    db = parseWhereParam(db, query.Where)
    db = db.Updates(data)
    if err := db.Error; err != nil {
        syncd.App.Logger.Warning("mysql query error: %s, sql[%v]", err.Error(), db.QueryExpr())
        return false
    }
    return true
}

func parseWhereParam(db *gorm.DB, where []WhereParam) *gorm.DB {
    if len(where) == 0 {
        return db
    }
    var (
        plain []string
        prepare []interface{}
    )
    for _, w := range where {
        tag := w.Tag
        if tag == "" {
            tag = "="
        }
        var plainFmt string
        switch tag {
        case "IN":
            plainFmt = fmt.Sprintf("%s IN (?)", w.Field)
        default:
            plainFmt = fmt.Sprintf("%s %s ?", w.Field, tag)
        }
        plain = append(plain, plainFmt)
        prepare = append(prepare, w.Prepare)
    }
    return db.Where(strings.Join(plain, " AND "), prepare...)
}
