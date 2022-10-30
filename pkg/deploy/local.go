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

package deploy

import (
	"context"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/docker_utils"
	"github.com/mitchellh/go-homedir"
)

type DevDeployer interface {
	Setup(ctx context.Context, coreImg, workerImg string) error

	CreateFLNetwork(ctx context.Context) error
	PullCoreImage(ctx context.Context) error
	PullWorkerImage(ctx context.Context) error
	PullPromImage(ctx context.Context) error
	StartCore(ctx context.Context) error
	StartWorker(ctx context.Context) error
	StartProm(ctx context.Context) error

	RemoveFLNetwork(ctx context.Context) error
	RemoveCoreContainer(ctx context.Context) error
	RemoveWorkerContainer(ctx context.Context) error
}

type LocalDeployer struct {
	client   *client.Client
	logsPath string

	flNetId   string
	flNetName string

	coreImg             string
	coreContainerName   string
	workerImg           string
	workerContainerName string

	promContainerName string
}

func NewDevDeployer(coreContainerName, workerContainerName, flNetName string) DevDeployer {
	return &LocalDeployer{
		flNetName:           flNetName,
		coreContainerName:   coreContainerName,
		workerContainerName: workerContainerName,
		promContainerName:   "fl-prometheus",
	}
}

func (d *LocalDeployer) Setup(ctx context.Context, coreImg, workerImg string) error {
	d.coreImg = coreImg
	d.workerImg = workerImg

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return err
	}
	d.client = cli

	h, err := homedir.Dir()
	if err != nil {
		return err
	}
	logsPath := filepath.Join(h, "funless-logs")
	if err := os.MkdirAll(logsPath, 0755); err != nil {
		return err
	}

	d.logsPath = logsPath
	return nil
}

func (d *LocalDeployer) CreateFLNetwork(ctx context.Context) error {
	// Network for Core + Worker
	exists, id, err := docker_utils.NetExists(ctx, d.client, d.flNetName)
	if err != nil {
		return err
	}
	if exists {
		d.flNetId = id
		return nil
	}
	id, err = docker_utils.NetCreate(ctx, d.client, d.flNetName, false)
	if err != nil {
		return err
	}
	d.flNetId = id
	return nil
}

func (d *LocalDeployer) PullCoreImage(ctx context.Context) error {
	return docker_utils.PullImage(ctx, d.client, d.coreImg)
}

func (d *LocalDeployer) PullWorkerImage(ctx context.Context) error {
	return docker_utils.PullImage(ctx, d.client, d.workerImg)
}

func (d *LocalDeployer) PullPromImage(ctx context.Context) error {
	return docker_utils.PullImage(ctx, d.client, pkg.Prometheus)
}

func (d *LocalDeployer) StartCore(ctx context.Context) error {
	containerConfig := coreContainerConfig(d.coreImg)
	hostConfig := coreHostConfig(d.logsPath)
	netConf := networkConfig(d.flNetName, d.flNetId)
	configs := configs(d.coreContainerName, containerConfig, hostConfig, netConf)
	return docker_utils.RunContainer(ctx, d.client, configs)
}

func (d *LocalDeployer) StartWorker(ctx context.Context) error {
	containerConfig := workerContainerConfig(d.workerImg)
	hostConf := workerHostConfig(d.logsPath)
	netConf := networkConfig(d.flNetName, d.flNetId)
	configs := configs(d.workerContainerName, containerConfig, hostConf, netConf)
	return docker_utils.RunContainer(ctx, d.client, configs)
}

func (d *LocalDeployer) StartProm(ctx context.Context) error {
	containerConfig := promContainerConfig()
	hostConf := promHostConfig()
	netConf := networkConfig(d.flNetName, d.flNetId)
	configs := configs(d.promContainerName, containerConfig, hostConf, netConf)
	return docker_utils.RunContainer(ctx, d.client, configs)
}

func (d *LocalDeployer) RemoveFLNetwork(ctx context.Context) error {
	return docker_utils.RemoveNetwork(ctx, d.client, d.flNetName)
}

func (d *LocalDeployer) RemoveCoreContainer(ctx context.Context) error {
	return docker_utils.RemoveContainer(ctx, d.client, d.coreContainerName)
}

func (d *LocalDeployer) RemoveWorkerContainer(ctx context.Context) error {
	return docker_utils.RemoveContainer(ctx, d.client, d.workerContainerName)
}

func coreContainerConfig(coreImg string) *container.Config {
	return &container.Config{
		Image: coreImg,
		ExposedPorts: nat.PortSet{
			"4000/tcp": struct{}{},
		},
		Env:     []string{"SECRET_KEY_BASE=" + pkg.FLCoreDevSecretKey},
		Volumes: map[string]struct{}{},
	}
}
func coreHostConfig(logsPath string) *container.HostConfig {
	return &container.HostConfig{
		PortBindings: nat.PortMap{
			"4000/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "4000",
				},
			},
		},
		Mounts: []mount.Mount{
			{
				Source: logsPath,
				Target: "/tmp/funless",
				Type:   mount.TypeBind,
			},
		},
	}
}
func workerContainerConfig(workerImg string) *container.Config {
	return &container.Config{
		Image: workerImg,
	}
}
func workerHostConfig(logsPath string) *container.HostConfig {
	return &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Source: logsPath,
				Target: "/tmp/funless",
				Type:   mount.TypeBind,
			},
		},
	}
}
func promContainerConfig() *container.Config {
	return &container.Config{
		Image: pkg.Prometheus,
	}
}
func promHostConfig() *container.HostConfig {
	return &container.HostConfig{}
}
func networkConfig(networkName, networkID string) *network.NetworkingConfig {
	endpoints := make(map[string]*network.EndpointSettings, 1)
	endpoints[networkName] = &network.EndpointSettings{
		NetworkID: networkID,
	}

	return &network.NetworkingConfig{
		EndpointsConfig: endpoints,
	}
}

func configs(name string, c *container.Config, h *container.HostConfig, n *network.NetworkingConfig) docker_utils.ContainerConfigs {
	return docker_utils.ContainerConfigs{
		ContName:   name,
		Container:  c,
		Host:       h,
		Networking: n,
	}
}
