// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golog

import (
    "io"
    "os"
    "fmt"
    "time"
    "bytes"
    "strings"
    "sync"
)

const (
    LEVEL_DEBUG = iota
    LEVEL_INFO
    LEVEL_NOTICE
    LEVEL_WARNING
    LEVEL_ERROR
    LEVEL_PANIC
)

var levelMessage = map[int]string{
    LEVEL_DEBUG: "debug",
    LEVEL_INFO: "info",
    LEVEL_NOTICE: "notice",
    LEVEL_WARNING: "warning",
    LEVEL_ERROR: "error",
    LEVEL_PANIC: "panic",
}

type Logger struct {
    output  io.Writer
    mu      sync.Mutex
}

func New(out io.Writer) *Logger {
    logger := &Logger{
        output: out,
    }
    return logger
}

var stdLog = New(os.Stderr)

func (l *Logger) Output(level int, s string) error {
    levelMsg, _ := levelMessage[level]
    var buf bytes.Buffer
    buf.WriteByte('[')
    buf.WriteString(time.Now().Format("2006-01-02 15:04:05"))
    buf.WriteString("] ")
    buf.WriteString(strings.ToUpper(levelMsg))
    buf.WriteString(": ")
    buf.WriteString(s)
    if len(s) == 0 || s[len(s)-1] != '\n' {
        buf.WriteByte('\n')
    }
    l.output.Write(buf.Bytes())
    return nil
}

func (l *Logger) SetOutput(w io.Writer) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.output = w
}

func (l *Logger) GetOutput() io.Writer {
    return l.output
}

func (l *Logger) Debug(s string, v ...interface{}) {
    l.Output(LEVEL_DEBUG, fmt.Sprintf(s, v...))
}

func (l *Logger) Info(s string, v ...interface{}) {
    l.Output(LEVEL_INFO, fmt.Sprintf(s, v...))
}

func (l *Logger) Notice(s string, v ...interface{}) {
    l.Output(LEVEL_NOTICE, fmt.Sprintf(s, v...))
}

func (l *Logger) Warning(s string, v ...interface{}) {
    l.Output(LEVEL_WARNING, fmt.Sprintf(s, v...))
}

func (l *Logger) Error(s string, v ...interface{}) {
    l.Output(LEVEL_ERROR, fmt.Sprintf(s, v...))
}

func (l *Logger) Panic(s string, v ...interface{}) {
    l.Output(LEVEL_PANIC, fmt.Sprintf(s, v...))
    panic(s)
}

func SetOutput(w io.Writer) {
    stdLog.SetOutput(w)
}

func Debug(s string, v ...interface{}) {
    stdLog.Debug(s, v...)
}

func Info(s string, v ...interface{}) {
    stdLog.Info(s, v...)
}

func Notice(s string, v ...interface{}) {
    stdLog.Notice(s, v...)
}

func Warning(s string, v ...interface{}) {
    stdLog.Warning(s, v...)
}

func Error(s string, v ...interface{}) {
    stdLog.Error(s, v...)
}

func Panic(s string, v ...interface{}) {
    stdLog.Panic(s, v...)
}