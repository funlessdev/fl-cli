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
package deploy

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func RemoveFLNetwork(ctx context.Context, client *client.Client, flNetName string) error {
	exists, net, err := flNetExists(ctx, client, flNetName)

	if err != nil {
		return err
	}
	if exists {
		return client.NetworkRemove(ctx, net.ID)
	}
	return nil
}

func RemoveFLContainer(ctx context.Context, client *client.Client, name string) error {
	exists, container, err := flContainerExists(ctx, client, name)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	timeout := time.Duration(5 * time.Second)
	if err := client.ContainerStop(ctx, container.ID, &timeout); err != nil {
		return err
	}

	return client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
}

func RemoveFunctionContainers(ctx context.Context, client *client.Client) error {
	containers, err := functionContainersList(ctx, client)
	if err != nil {
		return err
	}
	for _, container := range containers {
		timeout := time.Duration(2 * time.Second)
		if err := client.ContainerStop(ctx, container.ID, &timeout); err != nil {
			return err
		}
		if err := client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
			return err
		}

	}
	return nil
}
