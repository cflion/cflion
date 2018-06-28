//  Copyright (c) 2018, cflion
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
var defaultLogger = NewLogger(os.Stdout)

// Logger represents a simple defaultLogger with level.
// The inline logger is the standard Go logging "log".
type Logger struct {
	level  int
	logger *golog.Logger
}

// NewLogger creates a logger.
func NewLogger(out io.Writer) *Logger {
	logger := &Logger{level: INFO, logger: golog.New(out, "", golog.Ldate|golog.Ltime|golog.Lshortfile)}
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
	defaultLogger.level = getLevel(level)
}

// Trace prints trace level message of the default logger.
func Trace(v ...interface{}) {
	defaultLogger.Trace(v...)
}

// Tracef prints trace level message of the default logger with format.
func Tracef(format string, v ...interface{}) {
	defaultLogger.Tracef(format, v...)
}

// Debug prints debug level message of the default logger.
func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}

// Debugf prints debug level message of the default logger with format.
func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

// Info prints info level message of the default logger.
func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Infof prints info level message of the default logger with format.
func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

// Warn prints warn level message of the default logger.
func Warn(v ...interface{}) {
	defaultLogger.Warn(v...)
}

// Warnf prints warn level message of the default logger with format.
func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}

// Error prints error level message of the default logger.
func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

// Errorf prints error level message of the default logger withe format.
func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

// SetLevel sets the logging level of a logger.
func (l *Logger) SetLevel(level string) {
	l.level = getLevel(level)
}

// IsTraceEnabled determines whether the trace level is enabled.
func (l *Logger) IsTraceEnabled() bool {
	return l.level <= TRACE
}

// IsDebugEnabled determines whether the debug level is enabled.
func (l *Logger) IsDebugEnabled() bool {
	return l.level <= DEBUG
}

// IsInfoEnabled determines whether the info level is enabled.
func (l *Logger) IsInfoEnabled() bool {
	return l.level <= INFO
}

// IsWarnEnabled determines whether the warn level is enabled.
func (l *Logger) IsWarnEnabled() bool {
	return l.level <= WARN
}

// IsErrorEnabled determines whether the error level is enabled.
func (l *Logger) IsErrorEnabled() bool {
	return l.level <= ERROR
}

// TRACE prints trace level message.
func (l *Logger) Trace(v ...interface{}) {
	if !l.IsTraceEnabled() {
		return
	}
	l.logger.SetPrefix("T ")
	l.logger.Output(2, fmt.Sprint(v...))
}

// Tracef prints trace level message with format.
func (l *Logger) Tracef(format string, v ...interface{}) {
	if !l.IsTraceEnabled() {
		return
	}
	l.logger.SetPrefix("T ")
	l.logger.Output(2, fmt.Sprintf(format, v...))
}

// DEBUG prints debug level message.
func (l *Logger) Debug(v ...interface{}) {
	if !l.IsDebugEnabled() {
		return
	}
	l.logger.SetPrefix("D ")
	l.logger.Output(2, fmt.Sprint(v...))
}

// Debugf prints debug level message with format.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if !l.IsDebugEnabled() {
		return
	}
	l.logger.SetPrefix("D ")
	l.logger.Output(2, fmt.Sprintf(format, v...))
}

// INFO prints info level message.
func (l *Logger) Info(v ...interface{}) {
	if !l.IsInfoEnabled() {
		return
	}
	l.logger.SetPrefix("I ")
	l.logger.Output(2, fmt.Sprint(v...))
}

// Infof prints info level message with format.
func (l *Logger) Infof(format string, v ...interface{}) {
	if !l.IsInfoEnabled() {
		return
	}
	l.logger.SetPrefix("I ")
	l.logger.Output(2, fmt.Sprintf(format, v...))
}

// Warn prints warn level message.
func (l *Logger) Warn(v ...interface{}) {
	if !l.IsWarnEnabled() {
		return
	}
	l.logger.SetPrefix("D ")
	l.logger.Output(2, fmt.Sprint(v...))
}

// Warnf prints warn level message with format.
func (l *Logger) Warnf(format string, v ...interface{}) {
	if !l.IsWarnEnabled() {
		return
	}
	l.logger.SetPrefix("D ")
	l.logger.Output(2, fmt.Sprintf(format, v...))
}

// ERROR prints error level message.
func (l *Logger) Error(v ...interface{}) {
	if !l.IsErrorEnabled() {
		return
	}
	l.logger.SetPrefix("E ")
	l.logger.Output(2, fmt.Sprint(v...))
}

// Errorf prints error level message with format.
func (l *Logger) Errorf(format string, v ...interface{}) {
	if !l.IsErrorEnabled() {
		return
	}
	l.logger.SetPrefix("E ")
	l.logger.Output(2, fmt.Sprintf(format, v...))
}
