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
	"os"
	"time"

	"github.com/theckman/yacspin"
)

type loggerBuilder struct {
	debug   bool
	spinCfg yacspin.Config
	err     error
	writer  io.Writer
}

func NewLoggerBuilder() builder {
	return &loggerBuilder{
		spinCfg: yacspin.Config{
			SuffixAutoColon:   true,
			StopCharacter:     "✓",
			StopColors:        []string{"fgGreen"},
			StopFailCharacter: "✗",
			StopFailColors:    []string{"fgRed"},
		},
		writer: os.Stdout,
	}
}

func (l *loggerBuilder) WithWriter(writer io.Writer) builder {
	l.writer = writer
	l.spinCfg.Writer = writer
	return l
}

func (l *loggerBuilder) WithDebug(b bool) builder {
	l.debug = b
	return l
}

func (l *loggerBuilder) SpinnerFrequency(freq time.Duration) builder {

	if freq <= time.Duration(0) {
		l.err = fmt.Errorf("spinner frequency must be greater than 0, %w", l.err)
		return l
	}

	if l.err != nil {
		return l
	}

	l.spinCfg.Frequency = freq
	return l
}

// // SpinnerCharSet sets the character set to use for the spinner animation [0 to 90].
func (l *loggerBuilder) SpinnerCharSet(charset int) builder {
	if charset < 0 || charset > 90 {
		l.err = fmt.Errorf("spinner character set must be between 0 and 90, %w", l.err)
		return l
	}
	if l.err != nil {
		return l
	}
	l.spinCfg.CharSet = yacspin.CharSets[charset]
	return l
}

// // Build returns a new logger instance.
func (l *loggerBuilder) Build() (FLogger, error) {

	if l.err != nil {
		return nil, l.err
	}

	s, err := yacspin.New(l.spinCfg)
	if err != nil {
		return nil, err
	}

	logger := &FLoggerImpl{
		debug:          l.debug,
		currentMessage: "",
		spinner:        s,
		writer:         l.writer,
	}
	return logger, nil
}
