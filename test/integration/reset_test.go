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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/funlessdev/fl-cli/internal/command/admin"
	"github.com/funlessdev/fl-cli/pkg/deploy"
	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestAdminResetRun(t *testing.T) {
	runIntegration := os.Getenv("INTEGRATION_TESTS")
	if runIntegration == "" {
		t.Skip("set INTEGRATION_TESTS (optionally with DOCKER_HOST) to run this test")
	}

	admCmd := admin.Admin{Dev: struct{}{}, Reset: struct{}{}}

	coreName := "fl-core-test"
	workerName := "fl-worker-test"
	flNetName := "fl-net-test"
	flRuntimeName := "fl-runtime-net-test"
	localDeployer := deploy.NewLocalDeployer(coreName, workerName, flNetName, flRuntimeName)

	b := log.NewLoggerBuilder()
	var outbuf bytes.Buffer
	logger, _ := b.WithDebug(true).WithWriter(&outbuf).Build()

	ctx := context.Background()

	err := admCmd.Dev.Run(ctx, localDeployer, logger)
	assert.NoError(t, err)

	err = admCmd.Reset.Run(ctx, localDeployer, logger)
	assert.NoError(t, err)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	assert.NoError(t, err)

	assertNetworksRemoved(t, ctx, cli, flNetName, flRuntimeName)
	assertContainersRemoved(t, ctx, cli, coreName, workerName)

}

func assertNetworksRemoved(t *testing.T, ctx context.Context, cli *client.Client, flNetName, flRuntimeName string) {
	t.Helper()

	nets, _ := cli.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: flNetName}),
	})
	assert.Empty(t, nets)

	nets, _ = cli.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: flRuntimeName}),
	})
	assert.Empty(t, nets)
}

func assertContainersRemoved(t *testing.T, ctx context.Context, cli *client.Client, coreName, workerName string) {
	t.Helper()

	containers, _ := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: coreName}),
	})
	assert.Empty(t, containers)

	containers, _ = cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: workerName}),
	})
	assert.Empty(t, containers)
	filter := filters.NewArgs(filters.KeyValuePair{Key: "name", Value: "funless"})
	filter.Contains("funless")

	containers, _ = cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filter,
	})

	assert.Empty(t, containers)
}
