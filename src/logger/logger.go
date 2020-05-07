package logger

import (
    "fmt"
    "log"
    "time"
)

const (
    TimeFormat = "2006-01-02 15:04:05"
    ERROR      = "ERROR"
    WARN       = "WARN"
    INFO       = "INFO"
    DEBUG      = "DEBUG"
)

func init() {
    log.SetFlags(0) // 不使用标准日志的格式
}

// Error error级别的日志
func Error(msg string) {
    printf(ERROR, fmt.Sprintf("%v\n", msg))
}

// Errorf error级别的日志
func Errorf(format string, v ...interface{}) {
    printf(ERROR, format, v...)
}

// Warn warn级别的日志
func Warn(msg string) {
    printf(WARN, fmt.Sprintf("%v\n", msg))
}

// Warnf warn级别的日志
func Warnf(format string, v ...interface{}) {
    printf(WARN, format, v...)
}

// Info info级别的日志
func Info(msg string) {
    printf(INFO, fmt.Sprintf("%v\n", msg))
}

// Infof info级别的日志
func Infof(format string, v ...interface{}) {
    printf(INFO, format, v...)
}

// Debug debug级别的日志
func Debug(msg string) {
    printf(DEBUG, fmt.Sprintf("%v\n", msg))
}

// Debugf debug级别的日志
func Debugf(format string, v ...interface{}) {
    printf(DEBUG, format, v...)
}

// printf 打印日志
func printf(level, format string, v ...interface{}) {
    log.Printf(fmt.Sprintf("%v [%v] %v", time.Now().Format(TimeFormat), level, format), v...)
}
