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
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type FLDockerClient struct {
	innerClient *client.Client
}

func NewDockerClient(client *client.Client) DockerClient {
	return &FLDockerClient{
		innerClient: client,
	}
}

func (c *FLDockerClient) ImageExists(ctx context.Context, image string) (bool, error) {
	_, _, err := c.innerClient.ImageInspectWithRaw(ctx, image)
	notFound := client.IsErrNotFound(err)

	// notFound being false means the error is something else
	// we still return false as we can't be sure the image actually exists
	if err != nil && !notFound {
		return false, err
	}

	if notFound {
		return false, nil
	}

	return true, nil
}

func (c *FLDockerClient) Pull(ctx context.Context, image string) error {
	out, err := c.innerClient.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()

	d := json.NewDecoder(out)

	var event *dockerEvent
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

// struct for decoding docker events, used in PullImage to check if an error occurred during pulling
type dockerEvent struct {
	Status         string `json:"status"`
	Error          string `json:"error"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

// Creates and starts a container and then waits for it to exit
func (c *FLDockerClient) RunAndWait(ctx context.Context, conf ContainerConfigs) error {
	resp, err := c.innerClient.ContainerCreate(ctx, conf.Container, conf.Host, conf.Networking, nil, conf.ContName)
	if err != nil {
		return err
	}

	if err := c.innerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	okC, errC := c.innerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	for {
		select {
		case <-okC:
			c.RemoveCtr(ctx, resp.ID)
			return nil
		case err = <-errC:
			return err
		}
	}
}

// Creates and starts a container and returns without waiting
func (c *FLDockerClient) RunAsync(ctx context.Context, conf ContainerConfigs) error {
	resp, err := c.innerClient.ContainerCreate(ctx, conf.Container, conf.Host, conf.Networking, nil, conf.ContName)
	if err != nil {
		return err
	}

	if err := c.innerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func (c *FLDockerClient) RemoveCtr(ctx context.Context, containerID string) error {
	return c.innerClient.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
}

// Checks if a container exists and returns (true, ID, nil) if it does
func (c *FLDockerClient) CtrExists(ctx context.Context, containerName string) (bool, string, error) {
	containers, err := c.innerClient.ContainerList(ctx, types.ContainerListOptions{
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

// Checks if a network exists and returns (true, ID, nil) if it does
func (c *FLDockerClient) NetworkExists(ctx context.Context, networkName string) (bool, string, error) {
	nets, err := c.innerClient.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: networkName}),
	})
	if err != nil {
		return false, "", err
	}

	if len(nets) == 0 {
		return false, "", nil
	}

	return true, nets[0].ID, nil
}

// Creates a network and returns the ID
func (c *FLDockerClient) CreateNetwork(ctx context.Context, networkName string) (string, error) {
	res, err := c.innerClient.NetworkCreate(ctx, networkName, types.NetworkCreate{})
	if err != nil {
		return "", err
	}
	return res.ID, nil
}

func (c *FLDockerClient) RemoveNetwork(ctx context.Context, networkID string) error {
	return c.innerClient.NetworkRemove(ctx, networkID)
}
