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
        Mail    *MailConfig
    }

    SyncdConfig struct {
        Dir     string
        Cipher  string
    }

    LogConfig struct {
        Path    string
    }

    MailConfig struct {
        Enable  int
        Smtp    string
        Port    int
        User    string
        Pass    string
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
