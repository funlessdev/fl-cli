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
package deploy

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
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/funlessdev/fl-cli/pkg"
)

type LocalDeployer struct {
	client           *client.Client
	flNetId          string
	flRuntimeNetId   string
	flNetName        string
	flRuntimeNetName string
}

func NewLocalDeployer(flNetName string, flRuntimeNetName string) *LocalDeployer {
	return &LocalDeployer{
		flNetName:        flNetName,
		flRuntimeNetName: flRuntimeNetName,
	}
}

func (d *LocalDeployer) SetupClient(ctx context.Context) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return err
	}
	d.client = cli
	return nil
}

func (d *LocalDeployer) SetupFLNetworks(ctx context.Context) error {
	// Network for Core + Worker
	exists, net, err := flNetExists(ctx, d.client, d.flNetName)
	if err != nil {
		return err
	}
	if exists {

		d.flNetId = net.ID
		return nil
	}
	id, err := flNetCreate(ctx, d.client, d.flNetName, false)
	if err != nil {
		return err
	}
	d.flNetId = id

	// Network for Worker + Runtimes
	exists, net, err = flNetExists(ctx, d.client, d.flRuntimeNetName)
	if err != nil {
		return err
	}
	if exists {

		d.flRuntimeNetId = net.ID
		return nil
	}
	runtimeId, err := flNetCreate(ctx, d.client, d.flRuntimeNetName, true)
	d.flRuntimeNetId = runtimeId

	return err
}

func (d *LocalDeployer) PullCoreImage(ctx context.Context) error {
	return pullFLImage(ctx, d.client, pkg.FLCore)
}

func (d *LocalDeployer) PullWorkerImage(ctx context.Context) error {
	return pullFLImage(ctx, d.client, pkg.FLWorker)
}

func (d *LocalDeployer) StartCore(ctx context.Context) error {

	containerConfig := &container.Config{
		Image: pkg.FLCore,
		ExposedPorts: nat.PortSet{
			"4001/tcp": struct{}{},
		},
		Volumes: map[string]struct{}{},
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
		Mounts: []mount.Mount{
			{
				Source: "/home/giusdp/funless-logs/",
				Target: "/tmp/funless",
				Type:   mount.TypeBind,
			},
		},
	}

	netConf := buildNetworkConfig(d.flNetName, d.flNetId)

	configs := configuration{
		container:  containerConfig,
		host:       hostConfig,
		networking: &netConf,
	}

	return startCoreContainer(ctx, d.client, configs, "fl-core")
}

func (d *LocalDeployer) StartWorker(ctx context.Context) error {

	dockerHost := getDockerHost()

	containerConfig := &container.Config{
		Image: pkg.FLWorker,
		Env:   []string{"RUNTIME_NETWORK=" + d.flNetName},
	}

	hostConf := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Source: dockerHost,
				Target: "/var/run/docker-host.sock",
				Type:   mount.TypeBind,
			},
			{
				Source: "/home/giusdp/funless-logs/",
				Target: "/tmp/funless",
				Type:   mount.TypeBind,
			},
		},
	}

	netConf := buildNetworkConfig(d.flNetName, d.flNetId)

	configs := configuration{
		container:  containerConfig,
		host:       hostConf,
		networking: &netConf,
	}
	return startWorkerContainer(ctx, d.client, configs, "fl-worker", d.flRuntimeNetId)
}

type configuration struct {
	container  *container.Config
	host       *container.HostConfig
	networking *network.NetworkingConfig
}

func pullFLImage(ctx context.Context, c *client.Client, image string) error {
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

func flNetCreate(ctx context.Context, client *client.Client, netName string, internal bool) (string, error) {
	res, err := client.NetworkCreate(ctx, netName, types.NetworkCreate{Internal: internal})
	if err != nil {
		return "", err
	}
	if res.Warning != "" {
		fmt.Printf("Warning creating fl_net network: %s\n", res.Warning)
	}
	return res.ID, nil
}

func startCoreContainer(ctx context.Context, c *client.Client, configs configuration, containerName string) error {
	resp, err := c.ContainerCreate(ctx, configs.container, configs.host, configs.networking, nil, containerName)

	if err != nil {
		return err
	}

	if err := c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func startWorkerContainer(ctx context.Context, c *client.Client, configs configuration, containerName, runtimeNetId string) error {
	resp, err := c.ContainerCreate(ctx, configs.container, configs.host, configs.networking, nil, containerName)

	if err != nil {
		return err
	}

	runtimeNetSettings := &network.EndpointSettings{
		NetworkID: runtimeNetId,
	}

	if err := c.NetworkConnect(ctx, runtimeNetId, resp.ID, runtimeNetSettings); err != nil {
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
