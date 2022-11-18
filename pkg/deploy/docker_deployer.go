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
	"github.com/docker/go-connections/nat"
	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/docker"
	"github.com/mitchellh/go-homedir"
)

type DockerDeployer interface {
	WithImages(coreImg, workerImg string)
	WithDockerClient(cli docker.DockerClient)
	WithLogs(path string) error

	CreateFLNetwork(ctx context.Context) error
	PullCoreImage(ctx context.Context) error
	PullWorkerImage(ctx context.Context) error
	PullPromImage(ctx context.Context) error
	StartCore(ctx context.Context) error
	StartWorker(ctx context.Context) error
	StartProm(ctx context.Context) error
}

type FLDockerDeployer struct {
	flDocker docker.DockerClient

	logsPath            string
	flNetId             string
	flNetName           string
	coreImg             string
	coreContainerName   string
	workerImg           string
	workerContainerName string
	promContainerName   string
}

func NewDockerDeployer(flNetName, coreCtrName, workerCtrName, promCtrName string) DockerDeployer {
	return &FLDockerDeployer{
		flNetName:           flNetName,
		coreContainerName:   coreCtrName,
		workerContainerName: workerCtrName,
		promContainerName:   promCtrName,
	}
}

func (d *FLDockerDeployer) WithImages(coreImg, workerImg string) {
	d.coreImg = coreImg
	d.workerImg = workerImg
}

func (d *FLDockerDeployer) WithDockerClient(cli docker.DockerClient) {
	d.flDocker = cli
}

func (d *FLDockerDeployer) WithLogs(path string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	logsPath := filepath.Join(home, path)
	err = os.MkdirAll(logsPath, os.ModePerm)
	if err != nil {
		return err
	}
	d.logsPath = logsPath
	return nil
}

func (d *FLDockerDeployer) CreateFLNetwork(ctx context.Context) error {
	// Network for Core + Worker
	exists, id, err := d.flDocker.NetworkExists(ctx, d.flNetName)
	if err != nil {
		return err
	}
	if exists {
		d.flNetId = id
		return nil
	}
	id, err = d.flDocker.CreateNetwork(ctx, d.flNetName)
	if err != nil {
		return err
	}
	d.flNetId = id
	return nil
}

func (d *FLDockerDeployer) PullCoreImage(ctx context.Context) error {
	return d.pull(ctx, d.coreImg)
}

func (d *FLDockerDeployer) PullWorkerImage(ctx context.Context) error {
	return d.pull(ctx, d.workerImg)
}

func (d *FLDockerDeployer) PullPromImage(ctx context.Context) error {
	return d.pull(ctx, pkg.PrometheusImg)
}

func (d *FLDockerDeployer) StartCore(ctx context.Context) error {
	containerConfig := coreContainerConfig(d.coreImg)
	hostConfig := coreHostConfig(d.logsPath)
	netConf := networkConfig(d.flNetName, d.flNetId)
	configs := configs(d.coreContainerName, containerConfig, hostConfig, netConf)
	return d.flDocker.RunAsync(ctx, configs)
}

func (d *FLDockerDeployer) StartWorker(ctx context.Context) error {
	containerConfig := workerContainerConfig(d.workerImg)
	hostConf := workerHostConfig(d.logsPath)
	netConf := networkConfig(d.flNetName, d.flNetId)
	configs := configs(d.workerContainerName, containerConfig, hostConf, netConf)
	return d.flDocker.RunAsync(ctx, configs)
}

func (d *FLDockerDeployer) StartProm(ctx context.Context) error {
	containerConfig := promContainerConfig()
	hostConf := promHostConfig()
	netConf := networkConfig(d.flNetName, d.flNetId)
	configs := configs(d.promContainerName, containerConfig, hostConf, netConf)
	return d.flDocker.RunAsync(ctx, configs)
}

func (d *FLDockerDeployer) pull(ctx context.Context, img string) error {
	exists, err := d.flDocker.ImageExists(ctx, img)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return d.flDocker.Pull(ctx, img)

}

func coreContainerConfig(coreImg string) *container.Config {
	return &container.Config{
		Image: coreImg,
		ExposedPorts: nat.PortSet{
			"4000/tcp": struct{}{},
		},
		Env:     []string{"SECRET_KEY_BASE=" + pkg.CoreDevSecretKey},
		Volumes: map[string]struct{}{},
	}
}
func coreHostConfig(logsPath string) *container.HostConfig {
	return &container.HostConfig{
		PortBindings: nat.PortMap{
			"4000/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
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
		Image: pkg.PrometheusImg,
	}
}
func promHostConfig() *container.HostConfig {
	return &container.HostConfig{
		PortBindings: nat.PortMap{
			"9090/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "9090",
				},
			},
		},
	}
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

func configs(name string, c *container.Config, h *container.HostConfig, n *network.NetworkingConfig) docker.ContainerConfigs {
	return docker.ContainerConfigs{
		ContName:   name,
		Container:  c,
		Host:       h,
		Networking: n,
	}
}
