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
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/docker/docker/client"
	"github.com/funlessdev/fl-cli/internal/command/admin"
	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/deploy"
	"github.com/funlessdev/fl-cli/pkg/docker"
	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

var flNet = "test-fl-net"
var coreName = "fl-core-test"
var workerName = "fl-worker-test"
var promName = "fl-prom-test"

func buildDockerClient(t *testing.T) docker.DockerClient {
	t.Helper()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	assert.NoError(t, err)
	return docker.NewDockerClient(cli)
}

func TestAdminDevRun(t *testing.T) {
	runIntegration := os.Getenv("INTEGRATION_TESTS")
	if runIntegration == "" {
		t.Skip("set INTEGRATION_TESTS (optionally with DOCKER_HOST) to run this test")
	}

	flDocker := buildDockerClient(t)
	admCmd := admin.Admin{}
	admCmd.Dev.CoreImage = pkg.CoreImg
	admCmd.Dev.WorkerImage = pkg.WorkerImg

	deployer := deploy.NewDockerDeployer(flNet, coreName, workerName, promName)
	remover := deploy.NewDockerRemover(flNet, coreName, workerName, promName)

	b := log.NewLoggerBuilder()
	var outbuf bytes.Buffer
	logger, _ := b.WithDebug(true).WithWriter(&outbuf).Build()

	ctx := context.Background()

	t.Run("should successfully deploy and remove funless when no errors occurr", func(t *testing.T) {
		err := admCmd.Dev.Run(ctx, deployer, logger)
		assert.NoError(t, err)

		assertContainer(t, flDocker, coreName)
		assertContainer(t, flDocker, workerName)
		assertContainer(t, flDocker, promName)
		assertNetwork(t, flDocker, flNet)

		err = admCmd.Reset.Run(ctx, remover, logger)
		assert.NoError(t, err)

		assertNetworkRemoved(t, flDocker, flNet)
		assertContainerRemoved(t, flDocker, coreName)
		assertContainerRemoved(t, flDocker, workerName)
		assertContainerRemoved(t, flDocker, promName)
	})

	t.Run("should successfully deploy without creating networks when they already exist", func(t *testing.T) {

		_ = deployer.CreateFLNetwork(ctx)
		assertNetwork(t, flDocker, flNet)

		err := admCmd.Dev.Run(ctx, deployer, logger)
		assert.NoError(t, err)

		assertContainer(t, flDocker, coreName)
		assertContainer(t, flDocker, workerName)
		assertContainer(t, flDocker, promName)
		assertNetwork(t, flDocker, flNet)

		err = admCmd.Reset.Run(ctx, remover, logger)
		assert.NoError(t, err)

		assertNetworkRemoved(t, flDocker, flNet)
		assertContainerRemoved(t, flDocker, coreName)
		assertContainerRemoved(t, flDocker, workerName)
		assertContainerRemoved(t, flDocker, promName)
	})

	t.Run("should fail when core is already running", func(t *testing.T) {
		_ = deployer.CreateFLNetwork(ctx)
		_ = deployer.PullCoreImage(ctx)
		_ = deployer.StartCore(ctx)
		assertContainer(t, flDocker, coreName)

		err := admCmd.Dev.Run(ctx, deployer, logger)
		assert.Error(t, err)

		err = admCmd.Reset.Run(ctx, remover, logger)
		assert.NoError(t, err)

		assertNetworkRemoved(t, flDocker, flNet)
		assertContainerRemoved(t, flDocker, coreName)
		assertContainerRemoved(t, flDocker, workerName)
		assertContainerRemoved(t, flDocker, promName)
	})

	t.Run("should create ~/funless-logs folder when successfully deployed", func(t *testing.T) {
		logFolder, err := homedir.Expand("~/funless-logs")
		assert.NoError(t, err)

		os.RemoveAll(logFolder) // cleanup folder from previous test runs

		err = admCmd.Dev.Run(ctx, deployer, logger)
		assert.NoError(t, err)

		assert.DirExists(t, logFolder)
		files, err := os.ReadDir(logFolder)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(files), 0)

		err = admCmd.Reset.Run(ctx, remover, logger)
		assert.NoError(t, err)

		assertNetworkRemoved(t, flDocker, flNet)
		assertContainerRemoved(t, flDocker, coreName)
		assertContainerRemoved(t, flDocker, workerName)
		assertContainerRemoved(t, flDocker, promName)

		err = os.RemoveAll(logFolder)
		assert.NoError(t, err)
		assert.NoDirExists(t, logFolder)
	})
}

func assertContainer(t *testing.T, flDocker docker.DockerClient, name string) {
	t.Helper()

	ctx := context.TODO()
	exists, id, err := flDocker.CtrExists(ctx, name)

	assert.True(t, exists)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func assertNetwork(t *testing.T, flDocker docker.DockerClient, name string) {
	t.Helper()

	ctx := context.TODO()
	exists, id, err := flDocker.NetworkExists(ctx, name)

	assert.True(t, exists)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func assertNetworkRemoved(t *testing.T, flDocker docker.DockerClient, flNetName string) {
	t.Helper()
	ctx := context.TODO()
	exists, _, err := flDocker.NetworkExists(ctx, flNetName)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func assertContainerRemoved(t *testing.T, flDocker docker.DockerClient, containerName string) {
	t.Helper()
	ctx := context.TODO()
	exists, _, err := flDocker.CtrExists(ctx, containerName)
	assert.NoError(t, err)
	assert.False(t, exists)
}
