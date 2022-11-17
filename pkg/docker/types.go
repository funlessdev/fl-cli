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

package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

type ContainerConfigs struct {
	ContName   string
	Container  *container.Config
	Host       *container.HostConfig
	Networking *network.NetworkingConfig
}

type DockerClient interface {
	ImageExists(ctx context.Context, image string) (bool, error)
	Pull(ctx context.Context, image string) error

	CtrExists(ctx context.Context, containerName string) (bool, string, error)
	RunAndWait(ctx context.Context, conf ContainerConfigs) error
	RunAsync(ctx context.Context, conf ContainerConfigs) error
	RemoveCtr(ctx context.Context, containerID string) error

	NetworkExists(ctx context.Context, networkName string) (bool, string, error)
	CreateNetwork(ctx context.Context, networkName string) (string, error)
	RemoveNetwork(ctx context.Context, networkID string) error
}
