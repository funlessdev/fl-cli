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

package docker_utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type ContainerConfigs struct {
	ContName   string
	Container  *container.Config
	Host       *container.HostConfig
	Networking *network.NetworkingConfig
}

func ImageExistsLocally(ctx context.Context, c *client.Client, image string) (bool, error) {
	_, _, err := c.ImageInspectWithRaw(ctx, image)
	notFound := client.IsErrNotFound(err)

	/* notFound being false means the error is something else; we still return false as we can't be sure the image actually exists */
	if err != nil && !notFound {
		return false, err
	}

	if notFound {
		return false, nil
	}

	return true, nil
}

func PullImage(ctx context.Context, c *client.Client, image string) error {
	exists, err := ImageExistsLocally(ctx, c, image)
	if exists {
		return nil
	}
	if err != nil {
		return err
	}

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

func RunAndWaitContainer(ctx context.Context, c *client.Client, config ContainerConfigs) error {
	resp, err := c.ContainerCreate(ctx, config.Container, config.Host, config.Networking, nil, config.ContName)
	if err != nil {
		return err
	}

	if err := c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	okC, errC := c.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	for {
		select {
		case <-okC:
			return nil
		case err = <-errC:
			return err
		}
	}
}

func RunContainer(ctx context.Context, c *client.Client, configs ContainerConfigs) error {
	resp, err := c.ContainerCreate(ctx, configs.Container, configs.Host, configs.Networking, nil, configs.ContName)
	if err != nil {
		return err
	}

	if err := c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func NetExists(ctx context.Context, client *client.Client, netName string) (bool, string, error) {
	nets, err := client.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: netName}),
	})
	if err != nil {
		return false, "", err
	}

	if len(nets) == 0 {
		return false, "", nil
	}

	return true, nets[0].ID, nil
}

func NetCreate(ctx context.Context, client *client.Client, netName string, internal bool) (string, error) {
	res, err := client.NetworkCreate(ctx, netName, types.NetworkCreate{Internal: internal})
	if err != nil {
		return "", err
	}
	return res.ID, nil
}

func RemoveContainer(ctx context.Context, c *client.Client, name string) error {
	exists, container, err := containerExists(ctx, c, name)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	return c.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true})
}

func RemoveNetwork(ctx context.Context, c *client.Client, netName string) error {
	exists, id, err := NetExists(ctx, c, netName)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	return c.NetworkRemove(ctx, id)
}

func containerExists(ctx context.Context, c *client.Client, containerName string) (bool, types.Container, error) {
	containers, err := c.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName}),
	})
	if err != nil {
		return false, types.Container{}, err
	}

	if len(containers) == 0 {
		return false, types.Container{}, nil
	}

	return true, containers[0], nil
}
