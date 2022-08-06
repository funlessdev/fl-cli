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
	"time"

	"github.com/theckman/yacspin"
)

type FLogger interface {
	SpinnerSuffix(suffix string)
	SpinnerMessage(msg string)
	StartSpinner(msg string)
	StopSpinner(success bool)

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

type BaseFLogger struct {
	debug   bool
	spinner *yacspin.Spinner
}

func NewBaseLogger(debug bool) (*BaseFLogger, error) {
	cfg := yacspin.Config{
		Frequency:         150 * time.Millisecond,
		Colors:            []string{"fgYellow"},
		CharSet:           yacspin.CharSets[59],
		SuffixAutoColon:   true,
		StopCharacter:     "✓",
		StopColors:        []string{"fgGreen"},
		StopMessage:       "done",
		StopFailCharacter: "✗",
		StopFailColors:    []string{"fgRed"},
		StopFailMessage:   "failed",
	}

	s, err := yacspin.New(cfg)
	if err != nil {
		return nil, err
	}
	return &BaseFLogger{debug: debug, spinner: s}, nil
}

func (l *BaseFLogger) SpinnerSuffix(suffix string) {
	l.spinner.Suffix(suffix)
}

func (l *BaseFLogger) SpinnerMessage(msg string) {
	l.spinner.Message(msg)
}

func (l *BaseFLogger) StartSpinner(msg string) {
	l.spinner.Message(msg)
	l.spinner.Start()
}

func (l *BaseFLogger) StopSpinner(success bool) {
	if success {
		l.spinner.Stop()
	} else {
		l.spinner.StopFail()
	}
}

func (l *BaseFLogger) Info(args ...interface{}) {
	fmt.Println(args...)
}

func (l *BaseFLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *BaseFLogger) Debug(args ...interface{}) {
	if l.debug {
		fmt.Print("DEBUG: ")
		fmt.Println(args...)
	}
}

func (l *BaseFLogger) Debugf(format string, args ...interface{}) {
	if l.debug {
		fmt.Print("DEBUG: ")
		fmt.Printf(format, args...)
	}
}
