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
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/funlessdev/funless-cli/pkg/docker"
	"github.com/funlessdev/funless-cli/pkg/log"
)

type Admin struct {
	Deploy deploy `cmd:"" help:"deploy sub sub command"`
}

type deploy struct {
}

func (d *deploy) Run(ctx context.Context, logger log.FLogger) error {
	if err := docker.RunPreflightChecks(logger); err != nil {
		return err
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	err = deployWithDocker(ctx, &dockerClient{cli}, logger)

	return err
}

type dockerClient struct {
	*client.Client
}

// Function to connect to docker, pull images and start containers
func deployWithDocker(ctx context.Context, cli *dockerClient, logger log.FLogger) error {
	logger.SpinnerSuffix("Deploying funless locally")
	logger.StartSpinner("pulling images... ")

	cli.pullImage(ctx, logger, FLCore)
	logger.Info("Core image pulled.")

	cli.pullImage(ctx, logger, FLWorker)
	logger.Info("Worker image pulled.")

	logger.StopSpinner(true)
	return nil
}

func (c *dockerClient) pullImage(ctx context.Context, logger log.FLogger, image string) error {
	out, err := c.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		logger.StopSpinner(false)
		return err
	}
	defer out.Close()

	d := json.NewDecoder(out)

	type Event struct {
		Status         string `json:"status"`
		Error          string `json:"error"`
		Progress       string `json:"progress"`
		ProgressDetail struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"progressDetail"`
	}

	var event *Event
	for {
		if err := d.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			logger.StopSpinner(false)
			return err
		}

		if event.Error != "" {
			logger.StopSpinner(false)
			return fmt.Errorf("error pulling image: %s", event.Error)
		}
	}
	return nil
}
