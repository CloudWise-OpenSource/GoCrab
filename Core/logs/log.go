// Copyright 2016 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

// Usage:
//
// import "github.com/CloudWise-OpenSource/GoCrab/Core/logs"
//
//	log := NewLogger(10000)
//	log.SetLogger("console", "")
//
//	> the first params stand for how many channel
//
// Use it like this:
//
//	log.Trace("trace")
//	log.Info("info")
//	log.Warn("warning")
//	log.Debug("debug")
//	log.Critical("critical")
//
package logs

import (
	"fmt"
	"path"
	"runtime"
	"sync"
)

// RFC5424 log message levels.
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

// Legacy loglevel constants to ensure backwards compatibility.
//
// Deprecated: will be removed in 1.5.0.
const (
	LevelInfo  = LevelInformational
	LevelTrace = LevelDebug
	LevelWarn  = LevelWarning
)

type loggerType func() LoggerInterface

// LoggerInterface defines the behavior of a log provider.
type LoggerInterface interface {
	Init(config string) error
	WriteMsg(msg string, level int) error
	Destroy()
	Flush()
}

var adapters = make(map[string]loggerType)

// Register makes a log provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, log loggerType) {
	if log == nil {
		panic("logs: Register provide is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("logs: Register called twice for provider " + name)
	}
	adapters[name] = log
}

// Logger is default logger in GoCrab application.
// it can contain several providers and log message into all providers.
type Logger struct {
	lock                sync.Mutex
	level               int
	enableFuncCallDepth bool
	loggerFuncCallDepth int
	msg                 chan *logMsg
	outputs             map[string]LoggerInterface
}

type logMsg struct {
	level int
	msg   string
}

// NewLogger returns a new Logger.
// channellen means the number of messages in chan.
// if the buffering chan is full, logger adapters write to file or other way.
func NewLogger(channellen int64) *Logger {
	bl := new(Logger)
	bl.level = LevelDebug
	bl.loggerFuncCallDepth = 2
	bl.msg = make(chan *logMsg, channellen)
	bl.outputs = make(map[string]LoggerInterface)
	//bl.SetLogger("console", "") // default output to console
	go bl.startLogger()
	return bl
}

// SetLogger provides a given logger adapter into Logger with config string.
// config need to be correct JSON as string: {"interval":360}.
func (bl *Logger) SetLogger(adaptername string, config string) error {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	if log, ok := adapters[adaptername]; ok {
		lg := log()
		err := lg.Init(config)
		bl.outputs[adaptername] = lg
		if err != nil {
			fmt.Println("logs.Logger.SetLogger: " + err.Error())
			return err
		}
	} else {
		return fmt.Errorf("logs: unknown adaptername %q (forgotten Register?)", adaptername)
	}
	return nil
}

// remove a logger adapter in Logger.
func (bl *Logger) DelLogger(adaptername string) error {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	if lg, ok := bl.outputs[adaptername]; ok {
		lg.Destroy()
		delete(bl.outputs, adaptername)
		return nil
	} else {
		return fmt.Errorf("logs: unknown adaptername %q (forgotten Register?)", adaptername)
	}
}

func (bl *Logger) writerMsg(loglevel int, msg string) error {
	if loglevel > bl.level {
		return nil
	}
	lm := new(logMsg)
	lm.level = loglevel
	if bl.enableFuncCallDepth {
		_, file, line, ok := runtime.Caller(bl.loggerFuncCallDepth)
		if _, filename := path.Split(file); filename == "log.go" && (line == 97 || line == 83) {
			_, file, line, ok = runtime.Caller(bl.loggerFuncCallDepth + 1)
		}
		if ok {
			_, filename := path.Split(file)
			lm.msg = fmt.Sprintf("[%s:%d] %s", filename, line, msg)
		} else {
			lm.msg = msg
		}
	} else {
		lm.msg = msg
	}
	bl.msg <- lm
	return nil
}

// Set log message level.
//
// If message level (such as LevelDebug) is higher than logger level (such as LevelWarning),
// log providers will not even be sent the message.
func (bl *Logger) SetLevel(l int) {
	bl.level = l
}

// set log funcCallDepth
func (bl *Logger) SetLogFuncCallDepth(d int) {
	bl.loggerFuncCallDepth = d
}

// enable log funcCallDepth
func (bl *Logger) EnableFuncCallDepth(b bool) {
	bl.enableFuncCallDepth = b
}

// start logger chan reading.
// when chan is not empty, write logs.
func (bl *Logger) startLogger() {
	for {
		select {
		case bm := <-bl.msg:
			for _, l := range bl.outputs {
				err := l.WriteMsg(bm.msg, bm.level)
				if err != nil {
					fmt.Println("ERROR, unable to WriteMsg:", err)
				}
			}
		}
	}
}

// Log EMERGENCY level message.
func (bl *Logger) Emergency(format string, v ...interface{}) {
	msg := fmt.Sprintf("[M] "+format, v...)
	bl.writerMsg(LevelEmergency, msg)
}

// Log ALERT level message.
func (bl *Logger) Alert(format string, v ...interface{}) {
	msg := fmt.Sprintf("[A] "+format, v...)
	bl.writerMsg(LevelAlert, msg)
}

// Log CRITICAL level message.
func (bl *Logger) Critical(format string, v ...interface{}) {
	msg := fmt.Sprintf("[C] "+format, v...)
	bl.writerMsg(LevelCritical, msg)
}

// Log ERROR level message.
func (bl *Logger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf("[E] "+format, v...)
	bl.writerMsg(LevelError, msg)
}

// Log WARNING level message.
func (bl *Logger) Warning(format string, v ...interface{}) {
	msg := fmt.Sprintf("[W] "+format, v...)
	bl.writerMsg(LevelWarning, msg)
}

// Log NOTICE level message.
func (bl *Logger) Notice(format string, v ...interface{}) {
	msg := fmt.Sprintf("[N] "+format, v...)
	bl.writerMsg(LevelNotice, msg)
}

// Log INFORMATIONAL level message.
func (bl *Logger) Informational(format string, v ...interface{}) {
	msg := fmt.Sprintf("[I] "+format, v...)
	bl.writerMsg(LevelInformational, msg)
}

// Log DEBUG level message.
func (bl *Logger) Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf("[D] "+format, v...)
	bl.writerMsg(LevelDebug, msg)
}

// Log WARN level message.
//
// Deprecated: compatibility alias for Warning(), Will be removed in 1.5.0.
func (bl *Logger) Warn(format string, v ...interface{}) {
	bl.Warning(format, v...)
}

// Log INFO level message.
//
// Deprecated: compatibility alias for Informational(), Will be removed in 1.5.0.
func (bl *Logger) Info(format string, v ...interface{}) {
	bl.Informational(format, v...)
}

// Log TRACE level message.
//
// Deprecated: compatibility alias for Debug(), Will be removed in 1.5.0.
func (bl *Logger) Trace(format string, v ...interface{}) {
	bl.Debug(format, v...)
}

// flush all chan data.
func (bl *Logger) Flush() {
	for _, l := range bl.outputs {
		l.Flush()
	}
}

// close logger, flush all chan data and destroy all adapters in Logger.
func (bl *Logger) Close() {
	for {
		if len(bl.msg) > 0 {
			bm := <-bl.msg
			for _, l := range bl.outputs {
				err := l.WriteMsg(bm.msg, bm.level)
				if err != nil {
					fmt.Println("ERROR, unable to WriteMsg (while closing logger):", err)
				}
			}
			continue
		}
		break
	}
	for _, l := range bl.outputs {
		l.Flush()
		l.Destroy()
	}
}
