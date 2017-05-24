// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

package GoCrab

import (
	"strings"

	"github.com/CloudWise-OpenSource/GoCrab/Core/logs"
)

// Log levels to control the logging output.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLevel(l int) {
	Logger.SetLevel(l)
}

func SetLogFuncCall(b bool) {
	Logger.EnableFuncCallDepth(b)
	Logger.SetLogFuncCallDepth(3)
}

// logger references the used application logger.
var Logger *logs.Logger

// SetLogger sets a new logger.
func SetLogger(adaptername string, config string) error {
	err := Logger.SetLogger(adaptername, config)
	if err != nil {
		return err
	}
	return nil
}

func Emergency(v ...interface{}) {
	Logger.Emergency(generateFmtStr(len(v)), v...)
}

func Alert(v ...interface{}) {
	Logger.Alert(generateFmtStr(len(v)), v...)
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	Logger.Critical(generateFmtStr(len(v)), v...)
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	Logger.Error(generateFmtStr(len(v)), v...)
}

// Warning logs a message at warning level.
func Warning(v ...interface{}) {
	if RunMode == RUNMODE_DEV {
		Logger.Warning(generateFmtStr(len(v)), v...)
	}
}

// Deprecated: compatibility alias for Warning(), Will be removed in 1.5.0.
func Warn(v ...interface{}) {
	Warning(v...)
}

func Notice(v ...interface{}) {
	if RunMode == RUNMODE_DEV {
		Logger.Notice(generateFmtStr(len(v)), v...)
	}
}

// Info logs a message at info level.
func Informational(v ...interface{}) {
	Logger.Informational(generateFmtStr(len(v)), v...)
}

// Deprecated: compatibility alias for Warning(), Will be removed in 1.5.0.
func Info(v ...interface{}) {
	Informational(v...)
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	if RunMode == RUNMODE_DEV {
		Logger.Debug(generateFmtStr(len(v)), v...)
	}
}

// Trace logs a message at trace level.
// Deprecated: compatibility alias for Warning(), Will be removed in 1.5.0.
func Trace(v ...interface{}) {
	Logger.Trace(generateFmtStr(len(v)), v...)
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}
