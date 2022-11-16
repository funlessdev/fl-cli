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
	"github.com/docker/docker/client"
)

type ContainerConfigs struct {
	ContName   string
	Container  *container.Config
	Host       *container.HostConfig
	Networking *network.NetworkingConfig
}

type ImageHandler interface {
	Exists(ctx context.Context, dockerClient *client.Client, image string) (bool, error)
	Pull(ctx context.Context, dockerClient *client.Client, image string) error
}

type ContainerHandler interface {
	Exists(ctx context.Context, dockerClient *client.Client, containerName string) (bool, string, error)
	RunAndWait(ctx context.Context, dockerClient *client.Client, conf ContainerConfigs) error
	RunAsync(ctx context.Context, dockerClient *client.Client, conf ContainerConfigs) error
	Remove(ctx context.Context, dockerClient *client.Client, containerID string) error
}

type NetworkHandler interface {
	Exists(ctx context.Context, dockerClient *client.Client, networkName string) (bool, string, error)
	Create(ctx context.Context, dockerClient *client.Client, networkName string) (string, error)
	Remove(ctx context.Context, dockerClient *client.Client, networkID string) error
}
