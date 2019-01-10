
// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package model

import (
    "fmt"
    "strings"

    "github.com/dreamans/syncd"
    "github.com/jinzhu/gorm"
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

func Create(tableName string, data interface{}) bool {
    db := syncd.Orm.Table(setTableName(tableName)).Create(data)
    if err := db.Error; err != nil {
        syncd.Logger.Warning("mysql query error: %v, sql[%v]", err, db.QueryExpr())
        return false
    }
    return true
}

func GetMulti(tableName string, data interface{}, query QueryParam) bool {
    db := syncd.Orm.Table(setTableName(tableName)).Offset(query.Offset)
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
    db.Find(data)
    if err := db.Error; err != nil {
        syncd.Logger.Warning("mysql query error: %v, sql[%v]", err, db.QueryExpr())
        return false
    }
    return true
}

func Count(tableName string, count *int, query QueryParam) bool {
    db := syncd.Orm.Table(setTableName(tableName))
    db = parseWhereParam(db, query.Where)
    db = db.Count(count)
    if err := db.Error; err != nil {
        syncd.Logger.Warning("mysql query error: %v, sql[%v]", err, db.QueryExpr())
        return false
    }
    return true
}

func GetOne(tableName string, data interface{}, query QueryParam) bool {
    db := syncd.Orm.Table(setTableName(tableName))
    if query.Fields != "" {
        db = db.Select(query.Fields)
    }
    db = parseWhereParam(db, query.Where)
    db = db.First(data)
    if err := db.Error; err != nil && !db.RecordNotFound() {
        syncd.Logger.Warning("mysql query error: %v, sql[%v]", err, db.QueryExpr())
        return false
    }
    return true
}

func GetByPk(tableName string, data interface{}, id interface{}) bool {
    db := syncd.Orm.Table(setTableName(tableName))
    db.First(data, id)
    if err := db.Error; err != nil && !db.RecordNotFound() {
        syncd.Logger.Warning("mysql query error: %v, sql[%v]", err, db.QueryExpr())
        return false
    }
    return true
}

func Update(tableName string, data interface{}, query QueryParam) bool {
    db := syncd.Orm.Table(setTableName(tableName))
    db = parseWhereParam(db, query.Where)
    db = db.Updates(data)
    if err := db.Error; err != nil {
        syncd.Logger.Warning("mysql query error: %v, sql[%v]", err, db.QueryExpr())
        return false
    }
    return true
}

func Delete(tableName string, data interface{}, query QueryParam) bool {
    if len(query.Where) == 0 {
        syncd.Logger.Warning("mysql query error: delete failed, where conditions cannot be empty")
        return false
    }
    db := syncd.Orm.Table(setTableName(tableName))
    db = parseWhereParam(db, query.Where)
    db = db.Delete(data)
    if err := db.Error; err != nil {
        syncd.Logger.Warning("mysql query error: %v, sql[%v]", err, db.QueryExpr())
        return false
    }
    return true
}

func DeleteByPk(tableName string, data interface{}) bool {
    db := syncd.Orm.Table(setTableName(tableName))
    db.Delete(data)
    if err := db.Error; err != nil {
        syncd.Logger.Warning("mysql query error: %v, sql[%v]", err, db.QueryExpr())
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

func setTableName(rawName string) string {
    return strings.Join([]string{
        syncd.DbInstance.GetTablePrefix(),
        rawName,
    }, "")
}
