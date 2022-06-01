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

	"github.com/funlessdev/funless-cli/pkg/docker"
	"github.com/funlessdev/funless-cli/pkg/log"
)

type Admin struct {
	Deploy deploy `cmd:"" help:"deploy sub sub command"`
}

type deploy struct {
}

func (d *deploy) Run(logger log.FLogger) error {
	err := docker.RunPreflightChecks(logger)
	if err != nil {
		return err
	}

	logger.SpinnerSuffix("Deploying funless locally")

	logger.StartSpinner("pulling component images")

	time.Sleep(2 * time.Second)

	logger.StopSpinner(true)

	logger.StartSpinner("uploading data")

	time.Sleep(2 * time.Second)

	logger.StopSpinner(true)

	return err
}
