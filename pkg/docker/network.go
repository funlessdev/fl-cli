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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type FLNetworkHandler struct{}

func NewFLNetworkHandler() NetworkHandler {
	return &FLNetworkHandler{}
}

// Checks if a network exists and returns (true, ID, nil) if it does
func (nh *FLNetworkHandler) Exists(ctx context.Context, dockerClient *client.Client, networkName string) (bool, string, error) {
	nets, err := dockerClient.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: networkName}),
	})
	if err != nil {
		return false, "", err
	}

	if len(nets) == 0 {
		return false, "", nil
	}

	return true, nets[0].ID, nil
}

// Creates a network and returns the ID
func (nh *FLNetworkHandler) Create(ctx context.Context, dockerClient *client.Client, networkName string) (string, error) {
	res, err := dockerClient.NetworkCreate(ctx, networkName, types.NetworkCreate{})
	if err != nil {
		return "", err
	}
	return res.ID, nil
}

func (nh *FLNetworkHandler) Remove(ctx context.Context, dockerClient *client.Client, networkID string) error {
	return dockerClient.NetworkRemove(ctx, networkID)
}
