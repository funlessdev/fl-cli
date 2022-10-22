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

package build

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func imageExistsLocally(ctx context.Context, c *client.Client, image string) (bool, error) {
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

func pullImage(ctx context.Context, c *client.Client, image string) error {
	exists, err := imageExistsLocally(ctx, c, image)
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

func runContainer(ctx context.Context, c *client.Client, hostConfig *container.HostConfig, containerConfig *container.Config, containerName string) error {
	resp, err := c.ContainerCreate(ctx, containerConfig, hostConfig, &network.NetworkingConfig{}, nil, containerName)

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
