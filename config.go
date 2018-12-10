// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package syncd

type (
    Config struct {
        Serve   *ServeConfig
        Db      *DbConfig
    }

    ServeConfig struct {
        Addr            string
        Log             string
        ReadTimeout     int
        WriteTimeout    int
        IdleTimeout     int
    }

    DbConfig struct {
        Host    string
        Port    string
        User    string
        Pass    string
        DBname  string
    }
)
