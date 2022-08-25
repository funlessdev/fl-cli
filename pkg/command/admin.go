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
package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/funlessdev/funless-cli/pkg/admin"
	"github.com/funlessdev/funless-cli/pkg/log"
)

type (
	Admin struct {
		Deploy deploy `cmd:"" help:"deploy 1 core and 1 worker locally with docker containers"`
		Reset  reset  `cmd:"" help:"removes the deployment of local containers"`
	}

	deploy struct{}
	reset  struct{}
)

func (d *deploy) Run(ctx context.Context, logger log.FLogger) error {
	logger.Info("Deploying funless locally...\n")

	// check client manages to connect to docker

	// check if fl_net network already exists
	// if not, create it
	// if yes, use it

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return err
	}

	id, err := admin.ObtainFLNet(ctx, cli, logger)

	if err != nil {
		if client.IsErrConnectionFailed(err) {
			return errors.New("could not connect to docker, please make sure docker is running and accessible")
		}
		return err
	}

	logger.Infof("FLNet network id: %s\n", id)
	// err = deployWithDocker(ctx, &dockerClient{cli}, logger)

	return err
}

type dockerClient struct {
	*client.Client
}

// Function to connect to docker, pull images and start containers
func deployWithDocker(ctx context.Context, cli *dockerClient, logger log.FLogger) error {
	if err := pullFLImages(ctx, cli, logger); err != nil {
		return logger.StopSpinner(err)
	}

	if err := startFLContainers(ctx, cli, logger); err != nil {
		return logger.StopSpinner(err)
	}

	logger.Info("\nDeployment complete!")
	logger.Info("You can now start using Funless! 🎉")
	return nil
}

func pullFLImages(ctx context.Context, cli *dockerClient, logger log.FLogger) error {
	logger.StartSpinner(fmt.Sprintf("pulling Core image (%s) 📦", FLCore))
	if err := logger.StopSpinner(cli.pullImage(ctx, FLCore)); err != nil {
		return err
	}

	logger.StartSpinner(fmt.Sprintf("pulling Worker image (%s) 📦📦", FLWorker))
	if err := logger.StopSpinner(cli.pullImage(ctx, FLWorker)); err != nil {
		return err
	}

	return nil
}

func (c *dockerClient) pullImage(ctx context.Context, image string) error {
	out, err := c.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
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
			return err
		}

		if event.Error != "" {
			return fmt.Errorf("pulling image: %s", event.Error)
		}
	}
	return nil
}

func startFLContainers(ctx context.Context, cli *dockerClient, logger log.FLogger) error {
	logger.StartSpinner("starting Core container 🎛️")

	if err := logger.StopSpinner(startCoreContainer(ctx, cli)); err != nil {
		return err
	}

	logger.StartSpinner("starting Worker container 👷")
	if err := logger.StopSpinner(startWorkerContainer(ctx, cli)); err != nil {
		return err
	}

	return nil

}

func startCoreContainer(ctx context.Context, client *dockerClient) error {
	containerConfig := &container.Config{
		Image: FLCore,
		ExposedPorts: nat.PortSet{
			"4001/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"4001/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "4001",
				},
			},
		},
	}

	return client.startContainer(ctx, containerConfig, hostConfig)

}

func startWorkerContainer(ctx context.Context, client *dockerClient) error {
	containerConfig := &container.Config{
		Image: FLWorker,
	}
	return client.startContainer(ctx, containerConfig, nil)
}

func (c *dockerClient) startContainer(ctx context.Context, containerConfig *container.Config, hostConfig *container.HostConfig) error {
	resp, err := c.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")

	if err != nil {
		return err
	}

	if err := c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}
