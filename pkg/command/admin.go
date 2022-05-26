// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package command

import (
	"time"

	"github.com/funlessdev/funless-cli/pkg/spinner"
	"github.com/theckman/yacspin"
)

type Admin struct {
	Deploy deploy `cmd:"" help:"deploy sub sub command"`
}

type deploy struct {
}

func (d *deploy) Run() error {
	cfg := yacspin.Config{
		Frequency:         150 * time.Millisecond,
		Colors:            []string{"fgYellow"},
		CharSet:           yacspin.CharSets[59],
		Suffix:            " deploying funless locally",
		SuffixAutoColon:   true,
		Message:           "pulling component images",
		StopCharacter:     "✓",
		StopColors:        []string{"fgGreen"},
		StopMessage:       "done",
		StopFailCharacter: "✗",
		StopFailColors:    []string{"fgRed"},
		StopFailMessage:   "failed",
	}
	spinner, err := spinner.CreateSpinner(cfg)
	if err != nil {
		return err
	}

	err = spinner.Start()
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)
	err = spinner.Stop()
	// doing some work

	spinner.Message("uploading data")
	err = spinner.Start()
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)

	err = spinner.Stop()

	// upload...
	time.Sleep(2 * time.Second)

	return err
}
