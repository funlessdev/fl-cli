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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/funlessdev/fl-cli/pkg"
	"github.com/mitchellh/go-homedir"
)

type LocalDeployer struct {
	client   *client.Client
	logsPath string

	flNetId   string
	flNetName string

	flRuntimeNetId   string
	flRuntimeNetName string

	coreImg             string
	coreContainerName   string
	workerImg           string
	workerContainerName string
}

func NewLocalDeployer(coreContainerName, workerContainerName, flNetName, flRuntimeNetName string) *LocalDeployer {
	return &LocalDeployer{
		flNetName:        flNetName,
		flRuntimeNetName: flRuntimeNetName,

		coreContainerName:   coreContainerName,
		workerContainerName: workerContainerName,
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

func (d *LocalDeployer) CreateFLNetworks(ctx context.Context) error {
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
	return pullFLImage(ctx, d.client, d.coreImg)
}

func (d *LocalDeployer) PullWorkerImage(ctx context.Context) error {
	return pullFLImage(ctx, d.client, d.workerImg)
}

func (d *LocalDeployer) StartCore(ctx context.Context) error {

	containerConfig := &container.Config{
		Image: d.coreImg,
		ExposedPorts: nat.PortSet{
			"4001/tcp": struct{}{},
		},
		Env:     []string{"SECRET_KEY_BASE=" + pkg.FLCoreDevSecretKey},
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
				Source: d.logsPath,
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

	return startCoreContainer(ctx, d.client, configs, d.coreContainerName)
}

func (d *LocalDeployer) StartWorker(ctx context.Context) error {

	dockerHost := getDockerHost()

	containerConfig := &container.Config{
		Image: d.workerImg,
		Env:   []string{"RUNTIME_NETWORK=" + d.flRuntimeNetName},
	}

	hostConf := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Source: dockerHost,
				Target: "/var/run/docker-host.sock",
				Type:   mount.TypeBind,
			},
			{
				Source: d.logsPath,
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
	return startWorkerContainer(ctx, d.client, configs, d.workerContainerName, d.flRuntimeNetId)
}

func (d *LocalDeployer) RemoveFLNetworks(ctx context.Context) error {
	if err := removeNetwork(ctx, d.client, d.flNetName); err != nil {
		return err
	}
	return removeNetwork(ctx, d.client, d.flRuntimeNetName)
}

func (d *LocalDeployer) RemoveCoreContainer(ctx context.Context) error {
	return removeContainer(ctx, d.client, d.coreContainerName)
}

func (d *LocalDeployer) RemoveWorkerContainer(ctx context.Context) error {
	return removeContainer(ctx, d.client, d.workerContainerName)
}

func (d *LocalDeployer) RemoveFunctionContainers(ctx context.Context) error {
	containers, err := functionContainersList(ctx, d.client)
	if err != nil {
		return err
	}

	var removalErr error = nil
	for _, container := range containers {
		if err := d.client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			removalErr = err
		}
	}
	return removalErr
}
