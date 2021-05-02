package logger

import (
	"github.com/issue9/logs"
	"xorm.io/core"
)

// 适配xormlogger
type XormWrapLogger struct {
}

func (l XormWrapLogger) Debug(v ...interface{}) {
	logs.Debug(v...)
}
func (l XormWrapLogger) Debugf(format string, v ...interface{}) {
	logs.Debugf(format, v...)
}

func (l XormWrapLogger) Error(v ...interface{}) {
	logs.Error(v...)
}
func (l XormWrapLogger) Errorf(format string, v ...interface{}) {
	logs.Errorf(format, v...)
}
func (l XormWrapLogger) Info(v ...interface{}) {
	logs.Info(v...)
}
func (l XormWrapLogger) Infof(format string, v ...interface{}) {
	logs.Infof(format, v...)
}
func (l XormWrapLogger) Warn(v ...interface{}) {
	logs.Warn(v...)
}
func (l XormWrapLogger) Warnf(format string, v ...interface{}) {
	logs.Warnf(format, v...)
}

func (l XormWrapLogger) Level() core.LogLevel {
	return 1
}
func (l XormWrapLogger) SetLevel(v core.LogLevel) {
	return
}

func (l XormWrapLogger) ShowSQL(show ...bool) {
}

func (l XormWrapLogger) IsShowSQL() bool {
	return true
}
