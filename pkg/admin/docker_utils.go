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
package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type configuration struct {
	container  *container.Config
	host       *container.HostConfig
	networking *network.NetworkingConfig
}

func pullFLImage(ctx context.Context, c *client.Client, image string) error {
	if err := pullImage(ctx, c, image); err != nil {
		return err
	}
	return nil
}

func pullImage(ctx context.Context, c *client.Client, image string) error {
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

func flNetExists(ctx context.Context, client *client.Client, netName string) (bool, types.NetworkResource, error) {
	nets, err := client.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: netName}),
	})
	if err != nil {
		return false, types.NetworkResource{}, err
	}

	if len(nets) == 0 {
		return false, types.NetworkResource{}, nil
	}

	return true, nets[0], nil
}

func flNetCreate(ctx context.Context, client *client.Client, netName string) (string, error) {
	res, err := client.NetworkCreate(ctx, netName, types.NetworkCreate{})
	if err != nil {
		return "", err
	}
	if res.Warning != "" {
		fmt.Printf("Warning creating fl_net network: %s\n", res.Warning)
	}
	return res.ID, nil
}

func startContainer(ctx context.Context, c *client.Client, configs configuration, containerName string) error {
	resp, err := c.ContainerCreate(ctx, configs.container, configs.host, configs.networking, nil, containerName)

	if err != nil {
		return err
	}

	if err := c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func buildNetworkConfig(networkName, networkID string) network.NetworkingConfig {
	endpoints := make(map[string]*network.EndpointSettings, 1)
	endpoints[networkName] = &network.EndpointSettings{
		NetworkID: networkID,
	}

	return network.NetworkingConfig{
		EndpointsConfig: endpoints,
	}
}

func getDockerHost() string {
	dockerHost, exists := os.LookupEnv("DOCKER_HOST")
	if !exists || dockerHost == "" {
		dockerHost = "/var/run/docker.sock"
	} else {
		r, _ := regexp.Compile("^((unix|tcp|http)://)")
		dockerHost = r.ReplaceAllString(dockerHost, "")
	}
	return dockerHost
}

func flContainerExists(ctx context.Context, c *client.Client, containerName string) (bool, types.Container, error) {
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

func functionContainersList(ctx context.Context, c *client.Client) ([]types.Container, error) {

	// match all containers name containing funless suffix
	filter := filters.NewArgs(filters.KeyValuePair{Key: "name", Value: "funless"})
	filter.Contains("funless")

	containers, err := c.ContainerList(ctx, types.ContainerListOptions{
		Filters: filter,
	})

	if err != nil {
		return nil, err
	}

	return containers, nil
}
