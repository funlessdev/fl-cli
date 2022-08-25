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
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/funlessdev/funless-cli/pkg"
)

type LocalDeployer struct {
	ctx    context.Context
	client *client.Client
	fl_net string
}

func NewLocalDeployer(ctx context.Context, client *client.Client) *LocalDeployer {
	return &LocalDeployer{
		ctx:    ctx,
		client: client,
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
	exists, net, err := flNetExists(d.ctx, d.client)

	if err != nil {
		return err
	}
	if exists {

		d.fl_net = net.ID
		return nil
	}
	id, err := flNetCreate(d.ctx, d.client)
	d.fl_net = id
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

	return startContainer(d.ctx, d.client, containerConfig, hostConfig)
}

func StartWorkerContainer(d *LocalDeployer) error {

	containerConfig := &container.Config{
		Image: pkg.FLWorker,
	}
	return startContainer(d.ctx, d.client, containerConfig, nil)
}
