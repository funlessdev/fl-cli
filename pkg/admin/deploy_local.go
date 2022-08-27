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
	"os"
	"regexp"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/funlessdev/funless-cli/pkg"
)

type LocalDeployer struct {
	ctx       context.Context
	client    *client.Client
	flNetId   string
	flNetName string
}

func NewLocalDeployer(ctx context.Context, client *client.Client, flNetName string) *LocalDeployer {
	return &LocalDeployer{
		ctx:       ctx,
		client:    client,
		flNetName: flNetName,
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

func SetupFLNetwork(d *LocalDeployer) error {
	exists, net, err := flNetExists(d.ctx, d.client, d.flNetName)

	if err != nil {
		return err
	}
	if exists {

		d.flNetId = net.ID
		return nil
	}
	id, err := flNetCreate(d.ctx, d.client, d.flNetName)
	d.flNetId = id
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

	return startContainer(d.ctx, d.client, configs)
}

func StartWorkerContainer(d *LocalDeployer) error {

	dockerHost, exists := os.LookupEnv("DOCKER_HOST")
	if !exists || dockerHost == "" {
		dockerHost = "/var/run/docker.sock"
	} else {
		r, _ := regexp.Compile("^((unix|tcp|http)://)")
		dockerHost = r.ReplaceAllString(dockerHost, "")
	}

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
	return startContainer(d.ctx, d.client, configs)
}
