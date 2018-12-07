// Copyright 2018 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
    "io"
    "os"
    "fmt"
    "time"
    "bytes"
    "strings"
)

type Logger struct {
    Output  io.Writer
}

func New(out io.Writer) *Logger {
    logger := &Logger{
        Output: out,
    }
    return logger
}

var stdLog = New(os.Stderr)

func (l *Logger) Log(level, s string) error {
    var buf bytes.Buffer
    buf.WriteByte('[')
    buf.WriteString(time.Now().Format("2006-01-02 15:04:05"))
    buf.WriteString("] ")
    buf.WriteString(strings.ToUpper(level))
    buf.WriteString(": ")
    buf.WriteString(s)
    if len(s) == 0 || s[len(s)-1] != '\n' {
        buf.WriteByte('\n')
    }
    l.Output.Write(buf.Bytes())
    return nil
}

func Info(s string, v ...interface{}) {
    stdLog.Log("info", fmt.Sprintf(s, v...))
}
