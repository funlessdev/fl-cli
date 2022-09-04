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
	"errors"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/funlessdev/fl-cli/pkg"
)

type LocalDeployer struct {
	ctx              context.Context
	client           *client.Client
	flNetId          string
	flRuntimeNetId   string
	flNetName        string
	flRuntimeNetName string
}

func NewLocalDeployer(ctx context.Context, client *client.Client, flNetName string, flRuntimeNetName string) *LocalDeployer {
	return &LocalDeployer{
		ctx:              ctx,
		client:           client,
		flNetName:        flNetName,
		flRuntimeNetName: flRuntimeNetName,
	}
}

// Apply performs the steps to deploy funless locally
func (d *LocalDeployer) Apply(f func(*LocalDeployer) error) error {
	if err := f(d); err != nil {
		if client.IsErrConnectionFailed(err) {
			return errors.New("could not connect to docker, please make sure docker is running and accessible")
		}
		return err
	}
	return nil
}

func SetupFLNetworks(d *LocalDeployer) error {
	// Network for Core + Worker
	exists, net, err := flNetExists(d.ctx, d.client, d.flNetName)
	if err != nil {
		return err
	}
	if exists {

		d.flNetId = net.ID
		return nil
	}
	id, err := flNetCreate(d.ctx, d.client, d.flNetName, false)
	if err != nil {
		return err
	}
	d.flNetId = id

	// Network for Worker + Runtimes
	exists, net, err = flNetExists(d.ctx, d.client, d.flRuntimeNetName)
	if err != nil {
		return err
	}
	if exists {

		d.flRuntimeNetId = net.ID
		return nil
	}
	runtimeId, err := flNetCreate(d.ctx, d.client, d.flRuntimeNetName, true)
	d.flRuntimeNetId = runtimeId

	return err
}

func PullCoreImage(d *LocalDeployer) error {
	return pullFLImage(d.ctx, d.client, pkg.FLCore)
}

func PullWorkerImage(d *LocalDeployer) error {
	return pullFLImage(d.ctx, d.client, pkg.FLWorker)
}

func StartCoreContainer(d *LocalDeployer) error {
	containerConfig := &container.Config{
		Image: pkg.FLCore,
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

	netConf := buildNetworkConfig(d.flNetName, d.flNetId)

	configs := configuration{
		container:  containerConfig,
		host:       hostConfig,
		networking: &netConf,
	}

	return startCoreContainer(d.ctx, d.client, configs, "fl-core")
}

func StartWorkerContainer(d *LocalDeployer) error {

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
		},
	}

	netConf := buildNetworkConfig(d.flNetName, d.flNetId)

	configs := configuration{
		container:  containerConfig,
		host:       hostConf,
		networking: &netConf,
	}
	return startWorkerContainer(d.ctx, d.client, configs, "fl-worker", d.flRuntimeNetId)
}
