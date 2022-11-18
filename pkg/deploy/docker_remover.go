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

	"github.com/funlessdev/fl-cli/pkg/docker"
)

type DockerRemover interface {
	WithDockerClient(cli docker.DockerClient)
	RemoveFLNetwork(ctx context.Context) error
	RemoveCoreContainer(ctx context.Context) error
	RemoveWorkerContainer(ctx context.Context) error
	RemovePromContainer(ctx context.Context) error
}

type FLDockerRemover struct {
	flDocker docker.DockerClient

	flNetName           string
	coreContainerName   string
	workerContainerName string
	promContainerName   string
}

func NewDockerRemover(flNetName, coreCtrName, workerCtrName, promCtrName string) DockerRemover {
	return &FLDockerRemover{
		flNetName:           flNetName,
		coreContainerName:   coreCtrName,
		workerContainerName: workerCtrName,
		promContainerName:   promCtrName,
	}
}

func (r *FLDockerRemover) WithDockerClient(cli docker.DockerClient) {
	r.flDocker = cli
}

func (r *FLDockerRemover) RemoveFLNetwork(ctx context.Context) error {
	exists, id, err := r.flDocker.NetworkExists(ctx, r.flNetName)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	return r.flDocker.RemoveNetwork(ctx, id)
}

func (r *FLDockerRemover) RemoveCoreContainer(ctx context.Context) error {
	return r.remove(ctx, r.coreContainerName)
}

func (r *FLDockerRemover) RemoveWorkerContainer(ctx context.Context) error {
	return r.remove(ctx, r.workerContainerName)
}

func (r *FLDockerRemover) RemovePromContainer(ctx context.Context) error {
	return r.remove(ctx, r.promContainerName)
}

func (r *FLDockerRemover) remove(ctx context.Context, ctr string) error {
	exists, id, err := r.flDocker.CtrExists(ctx, ctr)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	return r.flDocker.RemoveCtr(ctx, id)
}
