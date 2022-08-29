// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package log

import (
	"fmt"
	"io"

	"github.com/theckman/yacspin"
)

type FLoggerImpl struct {
	debug          bool
	currentMessage string
	spinner        *yacspin.Spinner
	writer         io.Writer
}

func (l *FLoggerImpl) SpinnerSuffix(suffix string) {
	l.spinner.Suffix(suffix)
}

func (l *FLoggerImpl) SpinnerMessage(msg string) {
	l.currentMessage = msg
	l.spinner.Message(msg)
}

func (l *FLoggerImpl) StartSpinner(msg string) {
	l.currentMessage = msg
	l.spinner.Message(msg)
	_ = l.spinner.Start()
}

func (l *FLoggerImpl) StopSpinner(err error) error {
	if err == nil {
		l.spinner.StopMessage(l.currentMessage)
		return l.spinner.Stop()
	} else {
		l.spinner.StopFailMessage(l.currentMessage)
		_ = l.spinner.StopFail()
		return err
	}
}

func (l *FLoggerImpl) Info(args ...interface{}) {
	fmt.Fprintln(l.writer, args...)
}

func (l *FLoggerImpl) Infof(format string, args ...interface{}) {
	fmt.Fprintf(l.writer, format, args...)
}

func (l *FLoggerImpl) Debug(args ...interface{}) {
	if l.debug {
		fmt.Fprint(l.writer, "DEBUG: ")
		fmt.Fprintln(l.writer, args...)
	}
}

func (l *FLoggerImpl) Debugf(format string, args ...interface{}) {
	if l.debug {
		fmt.Fprint(l.writer, "DEBUG: ")
		fmt.Fprintf(l.writer, format, args...)
	}
}
