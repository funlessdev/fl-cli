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
package command

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/funlessdev/funless-cli/pkg"
	"github.com/funlessdev/funless-cli/pkg/admin"
	"github.com/funlessdev/funless-cli/pkg/log"
)

type (
	Admin struct {
		Deploy deploy `cmd:"" help:"deploy 1 core and 1 worker locally with docker containers"`
		Reset  reset  `cmd:"" help:"removes the deployment of local containers"`
	}

	deploy struct{}
	reset  struct{}
)

func (d *deploy) Run(ctx context.Context, logger log.FLogger) error {
	logger.Info("Deploying funless locally...\n")

	logger.StartSpinner("Setting things up...")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return err
	}
	deployer := admin.NewLocalDeployer(ctx, cli, "fl_net")

	if err := logger.StopSpinner(deployer.Apply(admin.SetupFLNetwork)); err != nil {
		return err
	}

	logger.StartSpinner(fmt.Sprintf("pulling Core image (%s) 📦", pkg.FLCore))
	if err := logger.StopSpinner(deployer.Apply(admin.PullCoreImage)); err != nil {
		return err
	}

	logger.StartSpinner(fmt.Sprintf("pulling Worker image (%s) 🗃", pkg.FLWorker))
	if err := logger.StopSpinner(deployer.Apply(admin.PullWorkerImage)); err != nil {
		return err
	}

	logger.StartSpinner("starting Core container 🎛️")
	if err := logger.StopSpinner(deployer.Apply(admin.StartCoreContainer)); err != nil {
		return err
	}

	logger.StartSpinner("starting Worker container 👷")
	if err := logger.StopSpinner(deployer.Apply(admin.StartWorkerContainer)); err != nil {
		return err
	}

	logger.Info("\nDeployment complete!")
	logger.Info("You can now start using Funless! 🎉")

	return err
}

func (r *reset) Run(ctx context.Context, logger log.FLogger) error {
	logger.Info("Removing local funless deployment...\n")

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return err
	}

	logger.StartSpinner("Removing Core container... ☠️")
	if err := logger.StopSpinner(admin.RemoveFLContainer(ctx, cli, "fl-core")); err != nil {
		return err
	}

	logger.StartSpinner("Removing Worker container... 🔪")
	if err := logger.StopSpinner(admin.RemoveFLContainer(ctx, cli, "fl-worker")); err != nil {
		return err
	}

	logger.StartSpinner("Removing the function containers... 🔫")
	if err := logger.StopSpinner(admin.RemoveFunctionContainers(ctx, cli)); err != nil {
		return err
	}

	logger.StartSpinner("Removing fl_net network... ✂️")
	if err := logger.StopSpinner(admin.RemoveFLNetwork(ctx, cli, "fl_net")); err != nil {
		return err
	}

	logger.Info("\nAll clear! 👍")

	return err
}
