//  Copyright (c) 2018 The cflion Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

// Package log includes logging related manipulations.
//
// log.SetLevel("Info")
//
// log.Info("info message")
//
// logger := log.NeLogger(os.Stdout)
//
// logger.SetLevel("Debug")
//
// logger.Trace("trace message)
// logger.Debug("debug message")
// logger.Info("info message")
// logger.Warn("warning message")
// logger.Error("error message")
//
// logger.Errorf("formatted %s message", "error")
package log

import (
	"fmt"
	"io"
	golog "log"
	"os"
	"strings"
)

// Logging level.
const (
	OFF = iota
	TRACE
	DEBUG
	INFO
	WARN
	ERROR
)

// defaultLogger prints message to the stdout.
var defaultLogger = &Logger{
	level: INFO,
	trace: golog.New(os.Stdout, "[TRACE] ", golog.Ldate|golog.Ltime|golog.Lshortfile),
	debug: golog.New(os.Stdout, "[DEBUG] ", golog.Ldate|golog.Ltime|golog.Lshortfile),
	info:  golog.New(os.Stdout, "[INFO] ", golog.Ldate|golog.Ltime|golog.Lshortfile),
	warn:  golog.New(os.Stdout, "[WARN] ", golog.Ldate|golog.Ltime|golog.Lshortfile),
	error: golog.New(os.Stdout, "[ERROR] ", golog.Ldate|golog.Ltime|golog.Lshortfile),
	depth: 3,
}

// Logger represents a simple defaultLogger with level.
// The inline logger is the standard Go logging "log".
type Logger struct {
	level int
	trace *golog.Logger
	debug *golog.Logger
	info  *golog.Logger
	warn  *golog.Logger
	error *golog.Logger
	// logger *golog.Logger
	depth int
}

// NewLogger creates a logger.
func NewLogger(out io.Writer) *Logger {
	logger := &Logger{
		level: INFO,
		trace: golog.New(out, "[TRACE] ", golog.Ldate|golog.Ltime|golog.Lshortfile),
		debug: golog.New(out, "[DEBUG] ", golog.Ldate|golog.Ltime|golog.Lshortfile),
		info:  golog.New(out, "[INFO] ", golog.Ldate|golog.Ltime|golog.Lshortfile),
		warn:  golog.New(out, "[WARN] ", golog.Ldate|golog.Ltime|golog.Lshortfile),
		error: golog.New(out, "[ERROR] ", golog.Ldate|golog.Ltime|golog.Lshortfile),
		depth: 2,
	}
	return logger
}

// getLevel gets the logging level int value corresponding to the specified level.
func getLevel(level string) int {
	level = strings.ToLower(level)
	switch level {
	case "off":
		return OFF
	case "trace":
		return TRACE
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	default:
		return INFO
	}
}

// SetLevel sets the logging level of the default logger.
func SetLevel(level string) {
	defaultLogger.SetLevel(level)
}

// SetLevel sets the logging level of a logger.
func (l *Logger) SetLevel(level string) {
	l.level = getLevel(level)
}

// SetOutput sets the writer of the default logger.
func SetOutput(w io.Writer) {
	defaultLogger.trace.SetOutput(w)
	defaultLogger.debug.SetOutput(w)
	defaultLogger.info.SetOutput(w)
	defaultLogger.warn.SetOutput(w)
	defaultLogger.error.SetOutput(w)
}

// IsTraceEnabled determines whether the trace level of the default logger is enabled.
func IsTraceEnabled() bool {
    return defaultLogger.IsTraceEnabled()
}

// IsTraceEnabled determines whether the trace level is enabled.
func (l *Logger) IsTraceEnabled() bool {
	return l.level <= TRACE
}

// IsTraceEnabled determines whether the debug level of the default logger is enabled.
func IsDebugEnabled() bool {
    return defaultLogger.IsDebugEnabled()
}

// IsDebugEnabled determines whether the debug level is enabled.
func (l *Logger) IsDebugEnabled() bool {
	return l.level <= DEBUG
}

// IsTraceEnabled determines whether the info level of the default logger is enabled.
func IsInfoEnabled() bool {
    return defaultLogger.IsInfoEnabled()
}

// IsInfoEnabled determines whether the info level is enabled.
func (l *Logger) IsInfoEnabled() bool {
	return l.level <= INFO
}

// IsTraceEnabled determines whether the warn level of the default logger is enabled.
func IsWarnEnabled() bool {
    return defaultLogger.IsWarnEnabled()
}

// IsWarnEnabled determines whether the warn level is enabled.
func (l *Logger) IsWarnEnabled() bool {
	return l.level <= WARN
}

// IsTraceEnabled determines whether the error level of the default logger is enabled.
func IsErrorEnabled() bool {
    return defaultLogger.IsErrorEnabled()
}

// IsErrorEnabled determines whether the error level is enabled.
func (l *Logger) IsErrorEnabled() bool {
	return l.level <= ERROR
}

// Trace prints trace level message of the default logger.
func Trace(v ...interface{}) {
	defaultLogger.Trace(v...)
}

// TRACE prints trace level message.
func (l *Logger) Trace(v ...interface{}) {
	if !l.IsTraceEnabled() {
		return
	}
	l.trace.Output(l.depth, fmt.Sprint(v...))
}

// Tracef prints trace level message of the default logger with format.
func Tracef(format string, v ...interface{}) {
	defaultLogger.Tracef(format, v...)
}

// Tracef prints trace level message with format.
func (l *Logger) Tracef(format string, v ...interface{}) {
	if !l.IsTraceEnabled() {
		return
	}
	l.trace.Output(l.depth, fmt.Sprintf(format, v...))
}

// Debug prints debug level message of the default logger.
func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}

// DEBUG prints debug level message.
func (l *Logger) Debug(v ...interface{}) {
	if !l.IsDebugEnabled() {
		return
	}
	l.debug.Output(l.depth, fmt.Sprint(v...))
}

// Debugf prints debug level message of the default logger with format.
func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

// Debugf prints debug level message with format.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if !l.IsDebugEnabled() {
		return
	}
	l.debug.Output(l.depth, fmt.Sprintf(format, v...))
}

// Info prints info level message of the default logger.
func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

// INFO prints info level message.
func (l *Logger) Info(v ...interface{}) {
	if !l.IsInfoEnabled() {
		return
	}
	l.info.Output(l.depth, fmt.Sprint(v...))
}

// Infof prints info level message of the default logger with format.
func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

// Infof prints info level message with format.
func (l *Logger) Infof(format string, v ...interface{}) {
	if !l.IsInfoEnabled() {
		return
	}
	l.info.Output(l.depth, fmt.Sprintf(format, v...))
}

// Warn prints warn level message of the default logger.
func Warn(v ...interface{}) {
	defaultLogger.Warn(v...)
}

// Warn prints warn level message.
func (l *Logger) Warn(v ...interface{}) {
	if !l.IsWarnEnabled() {
		return
	}
	l.warn.Output(l.depth, fmt.Sprint(v...))
}

// Warnf prints warn level message of the default logger with format.
func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}

// Warnf prints warn level message with format.
func (l *Logger) Warnf(format string, v ...interface{}) {
	if !l.IsWarnEnabled() {
		return
	}
	l.warn.Output(l.depth, fmt.Sprintf(format, v...))
}

// Error prints error level message of the default logger.
func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

// ERROR prints error level message.
func (l *Logger) Error(v ...interface{}) {
	if !l.IsErrorEnabled() {
		return
	}
	l.error.Output(l.depth, fmt.Sprint(v...))
}

// Errorf prints error level message of the default logger withe format.
func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

// Errorf prints error level message with format.
func (l *Logger) Errorf(format string, v ...interface{}) {
	if !l.IsErrorEnabled() {
		return
	}
	l.error.Output(l.depth, fmt.Sprintf(format, v...))
}
