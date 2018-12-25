// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package syncd

type (
    Config struct {
        Serve   *ServeConfig
        Db      *DbConfig
        Log     *LogConfig
        Syncd   *SyncdConfig
    }

    SyncdConfig struct {
        Dir string
    }

    LogConfig struct {
        Path    string
    }

    ServeConfig struct {
        Addr            string
        ReadTimeout     int
        WriteTimeout    int
        IdleTimeout     int
    }

    DbConfig struct {
        Unix            string
        Host            string
        Port            string
        Charset         string
        User            string
        Pass            string
        DbName          string
        TablePrefix     string
        MaxIdleConns    int
        MaxOpenConns    int
        ConnMaxLifeTime int
    }
)
