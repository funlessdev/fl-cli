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

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/docker/docker/client"
	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This file contains the integration tests for the docker package.
// It tests the ImageHandler, ContainerHandler and NetworkHandler concrete implementations.

func TestImageHandling(t *testing.T) {
	// If the environment variable is not set, we skip the test. (NOTE: with "run test" in vscode you're not passing the env var)
	runIntegration := os.Getenv("INTEGRATION_TESTS")
	if runIntegration == "" {
		t.Skip("set INTEGRATION_TESTS (optionally with DOCKER_HOST) to run this test")
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	require.NoError(t, err)
	ctx := context.TODO()
	flDocker := docker.NewDockerClient(dockerClient)

	t.Run("ImageExists should return false when an image is not present", func(t *testing.T) {
		exists, err := flDocker.ImageExists(ctx, "should_not_have/this_image:for_sure")
		require.NoError(t, err)
		require.False(t, exists, "image should not exist")
	})

	t.Run("ImageExists should returns true after Pull on CoreImg", func(t *testing.T) {
		t.Logf("DEBUG: Pulling image %s for real! It might take some time...", pkg.CoreImg)
		err := flDocker.Pull(ctx, pkg.CoreImg)
		require.NoError(t, err)

		exists, err := flDocker.ImageExists(ctx, pkg.CoreImg)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("Pull on already present image should be alright", func(t *testing.T) {
		err := flDocker.Pull(ctx, pkg.CoreImg)
		require.NoError(t, err)
	})
}

// func TestContainerHandler(t *testing.T) {
// 	// If the environment variable is not set, we skip the test. (NOTE: with "run test" in vscode you're not passing the env var)
// 	runIntegration := os.Getenv("INTEGRATION_TESTS")
// 	if runIntegration == "" {
// 		t.Skip("set INTEGRATION_TESTS (optionally with DOCKER_HOST) to run this test")
// 	}

// 	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
// 	require.NoError(t, err)
// 	ctx := context.TODO()
// 	flDocker := docker.NewDockerClient(dockerClient)

// 	t.Run("Exists should return false when a container is not present", func(t *testing.T) {
// 		exists, id, err := flDocker.CtrExists(ctx, "should_not_have_this_container")
// 		require.NoError(t, err)
// 		require.Empty(t, id)
// 		require.False(t, exists, "container should not exist")
// 	})

// 	t.Run("RunAsync should return an error if the image is not present", func(t *testing.T) {
// 		conf := docker.ContainerConfigs{
// 			ContName: "test_container",
// 			Container: &container.Config{
// 				Image: "should_not_have/this_image:for_sure",
// 			},
// 		}
// 		err := flDocker.RunAsync(ctx, conf)
// 		require.Error(t, err)
// 	})

// 	t.Run("RunAndWait should return an error if the image is not present", func(t *testing.T) {
// 		conf := docker.ContainerConfigs{
// 			ContName: "test_container",
// 			Container: &container.Config{
// 				Image: "should_not_have/this_image:for_sure",
// 			},
// 		}
// 		err := flDocker.RunAndWait(ctx, conf)
// 		require.Error(t, err)
// 	})

// 	t.Run("RunAndWait should return nil and remove the container when success", func(t *testing.T) {
// 		t.Log("DEBUG: Pulling image hello-world! It might take some time...")
// 		_ = flDocker.Pull(ctx, "hello-world:latest")

// 		contName := "test_container"

// 		conf := docker.ContainerConfigs{
// 			ContName: contName,
// 			Container: &container.Config{
// 				Image: "hello-world:latest",
// 			},
// 		}
// 		t.Log("DEBUG: running container hello-world")
// 		err := flDocker.RunAndWait(ctx, conf)
// 		require.NoError(t, err)

// 		t.Log("DEBUG: checking if container hello-world is still present")
// 		exists, id, err := flDocker.CtrExists(ctx, contName)
// 		assert.NoError(t, err)
// 		assert.False(t, exists)
// 		assert.Empty(t, id)
// 	})

// 	t.Run("RunAsync should return nil immediately and the container stays up when success", func(t *testing.T) {
// 		t.Log("DEBUG: Pulling Prometheys image! It might take some time...")
// 		_ = flDocker.Pull(ctx, pkg.PrometheusImg)

// 		contName := "test_prom_cont"

// 		conf := docker.ContainerConfigs{
// 			ContName: contName,
// 			Container: &container.Config{
// 				Image: pkg.PrometheusImg,
// 			},
// 		}
// 		t.Log("DEBUG: running core container")
// 		err := flDocker.RunAsync(ctx, conf)
// 		assert.NoError(t, err)

// 		t.Log("DEBUG: checking if container is still present")
// 		exists, id, err := flDocker.CtrExists(ctx, contName)
// 		assert.NoError(t, err)
// 		assert.True(t, exists)
// 		assert.NotEmpty(t, id)

// 		t.Log("DEBUG: removing container")
// 		err = flDocker.RemoveCtr(ctx, contName)
// 		assert.NoError(t, err)

// 		t.Log("DEBUG: checking if container is still present")
// 		exists, id, err = flDocker.CtrExists(ctx, contName)
// 		assert.NoError(t, err)
// 		assert.False(t, exists)
// 		assert.Empty(t, id)
// 	})
// }

func TestNetworkHandler(t *testing.T) {
	// If the environment variable is not set, we skip the test. (NOTE: with "run test" in vscode you're not passing the env var)
	runIntegration := os.Getenv("INTEGRATION_TESTS")
	if runIntegration == "" {
		t.Skip("set INTEGRATION_TESTS (optionally with DOCKER_HOST) to run this test")
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	require.NoError(t, err)
	ctx := context.TODO()
	flDocker := docker.NewDockerClient(dockerClient)

	t.Run("Exists should return false when a network is not present", func(t *testing.T) {
		exists, id, err := flDocker.NetworkExists(ctx, "should_not_have_this_network")
		require.NoError(t, err)
		require.False(t, exists, "network should not exist")
		require.Empty(t, id)
	})

	t.Run("Create a network, check it exists and remove it", func(t *testing.T) {
		networkName := "test_network"
		netId, err := flDocker.CreateNetwork(ctx, networkName)
		require.NoError(t, err)
		require.NotEmpty(t, netId)

		exists, id, err := flDocker.NetworkExists(ctx, networkName)
		assert.NoError(t, err)
		assert.True(t, exists, "network should exist")
		assert.Equal(t, netId, id)

		err = flDocker.RemoveNetwork(ctx, networkName)
		require.NoError(t, err)

		exists, id, err = flDocker.NetworkExists(ctx, networkName)
		require.NoError(t, err)
		require.False(t, exists, "network should not exist")
		require.Empty(t, id)
	})
}
