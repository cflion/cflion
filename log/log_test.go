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

package log

import (
	"os"
	"testing"
)

var logger = NewLogger(os.Stdout)

func TestSetLevel(t *testing.T) {
	SetLevel("TRACE")
}

func TestSetOutput(t *testing.T) {
	SetOutput(os.Stdout)
}

func TestLogger_IsTraceEnabled(t *testing.T) {
	logger.SetLevel("TRACE")
	if !logger.IsTraceEnabled() {
		t.FailNow()
	}
	logger.SetLevel("ERROR")
	if logger.IsTraceEnabled() {
		t.FailNow()
	}
}

func TestLogger_IsDebugEnabled(t *testing.T) {
	logger.SetLevel("OFF")
	if !logger.IsTraceEnabled() {
		t.FailNow()
	}
	logger.SetLevel("DEBUG")
	if !logger.IsTraceEnabled() {
		t.FailNow()
	}
	logger.SetLevel("INFO")
	if logger.IsTraceEnabled() {
		t.FailNow()
	}
}

func TestLogger_IsInfoEnabled(t *testing.T) {
	logger.SetLevel("INFO")
	if !logger.IsTraceEnabled() {
		t.FailNow()
	}
	logger.SetLevel("OTHER")
	if !logger.IsTraceEnabled() {
		t.FailNow()
	}
	logger.SetLevel("WARN")
	if logger.IsTraceEnabled() {
		t.FailNow()
	}
}

func TestLogger_IsWarnEnabled(t *testing.T) {
	logger.SetLevel("WARN")
	if !logger.IsTraceEnabled() {
		t.FailNow()
	}
	logger.SetLevel("ERROR")
	if logger.IsTraceEnabled() {
		t.FailNow()
	}
}

func TestLogger_IsErrorEnabled(t *testing.T) {
	logger.SetLevel("ERROR")
	if !logger.IsTraceEnabled() {
		t.FailNow()
	}
}

func TestTrace(t *testing.T) {
	Trace("trace")
}

func TestLogger_Trace(t *testing.T) {
	logger.SetLevel("TRACE")
	logger.Trace("trace")
	logger.SetLevel("INFO")
	logger.Trace("trace")
}

func TestTracef(t *testing.T) {
	Tracef("tracef")
}

func TestLogger_Tracef(t *testing.T) {
	logger.SetLevel("TRACE")
	logger.Tracef("tracef")
	logger.SetLevel("INFO")
	logger.Tracef("tracef")
}

func TestDebug(t *testing.T) {
	Debug("debug")
}

func TestLogger_Debug(t *testing.T) {
	logger.SetLevel("DEBUG")
	logger.Debug("debug")
	logger.SetLevel("INFO")
	logger.Debug("debug")
}

func TestDebugf(t *testing.T) {
	Debugf("debugf")
}

func TestLogger_Debugf(t *testing.T) {
	logger.SetLevel("DEBUG")
	logger.Debugf("debugf")
	logger.SetLevel("INFO")
	logger.Debugf("debugf")
}

func TestInfo(t *testing.T) {
	Info("info")
}

func TestLogger_Info(t *testing.T) {
	logger.SetLevel("INFO")
	logger.Info("info")
	logger.SetLevel("ERROR")
	logger.Info("info")
}

func TestInfof(t *testing.T) {
	Infof("infof")
}

func TestLogger_Infof(t *testing.T) {
	logger.SetLevel("INFO")
	logger.Infof("infof")
	logger.SetLevel("ERROR")
	logger.Infof("infof")
}

func TestWarn(t *testing.T) {
	Warn("warn")
}

func TestLogger_Warn(t *testing.T) {
	logger.SetLevel("WARN")
	logger.Warn("warn")
	logger.SetLevel("ERROR")
	logger.Warn("warn")
}

func TestWarnf(t *testing.T) {
	Warnf("warnf")
}

func TestLogger_Warnf(t *testing.T) {
	logger.SetLevel("WARN")
	logger.Warnf("warnf")
	logger.SetLevel("ERROR")
	logger.Warnf("warnf")
}

func TestError(t *testing.T) {
	Error("error")
}

func TestLogger_Error(t *testing.T) {
	logger.SetLevel("ERROR")
	logger.Error("error")
}

func TestErrorf(t *testing.T) {
	Errorf("errorf")
}

func TestLogger_Errorf(t *testing.T) {
	logger.SetLevel("ERROR")
	logger.Errorf("errorf")
}
