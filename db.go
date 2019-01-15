// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package syncd

import (
    "fmt"
    "time"

    "github.com/jinzhu/gorm"
    _ "github.com/go-sql-driver/mysql"
)

type DB struct {
    DbHandler   *gorm.DB
    cfg         *DbConfig
}

func NewDatabase(cfg *DbConfig) *DB {
    return &DB{
        cfg: cfg,
    }
}

func (db *DB) Open() error {
    c, err := gorm.Open("mysql", db.parseConnConfig())
    if err != nil {
        return err
    }

    //c.LogMode(true)
    c.DB().SetMaxIdleConns(db.cfg.MaxIdleConns)
    c.DB().SetMaxOpenConns(db.cfg.MaxOpenConns)
    c.DB().SetConnMaxLifetime(time.Second * time.Duration(db.cfg.ConnMaxLifeTime))

    db.DbHandler = c
    return nil
}

func (db *DB) Close() {
    db.DbHandler.Close()
}

func (db *DB) parseConnConfig() string {
    var connHost string
    if db.cfg.Unix != "" {
        connHost = fmt.Sprintf("unix(%s)", db.cfg.Unix)
    } else {
        connHost = fmt.Sprintf("tcp(%s:%s)", db.cfg.Host, db.cfg.Port)
    }
    s := fmt.Sprintf("%s:%s@%s/%s?charset=%s&parseTime=True&loc=Local", db.cfg.User, db.cfg.Pass, connHost, db.cfg.DbName, "utf8mb4")

    return s
}

func (db *DB) GetTablePrefix() string {
    return db.cfg.TablePrefix
}
