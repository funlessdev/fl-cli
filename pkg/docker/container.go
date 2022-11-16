// Copyright 2022 Giuseppe De Palma, Matteo Trentin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type FLContainerHandler struct{}

func NewFLContainerHandler() ContainerHandler {
	return &FLContainerHandler{}
}

// Creates and starts a container and then waits for it to exit
func (dr *FLContainerHandler) RunAndWait(ctx context.Context, dockerClient *client.Client, conf ContainerConfigs) error {
	resp, err := dockerClient.ContainerCreate(ctx, conf.Container, conf.Host, conf.Networking, nil, conf.ContName)
	if err != nil {
		return err
	}

	if err := dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	okC, errC := dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	for {
		select {
		case <-okC:
			dr.Remove(ctx, dockerClient, resp.ID)
			return nil
		case err = <-errC:
			return err
		}
	}
}

// Creates and starts a container and returns without waiting
func (dr *FLContainerHandler) RunAsync(ctx context.Context, dockerClient *client.Client, conf ContainerConfigs) error {
	resp, err := dockerClient.ContainerCreate(ctx, conf.Container, conf.Host, conf.Networking, nil, conf.ContName)
	if err != nil {
		return err
	}

	if err := dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func (dr *FLContainerHandler) Remove(ctx context.Context, dockerClient *client.Client, containerID string) error {
	return dockerClient.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
}

// Checks if a container exists and returns (true, ID, nil) if it does
func (dr *FLContainerHandler) Exists(ctx context.Context, dockerClient *client.Client, containerName string) (bool, string, error) {
	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName}),
	})
	if err != nil {
		return false, "", err
	}

	if len(containers) == 0 {
		return false, "", nil
	}

	return true, containers[0].ID, nil
}
