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

package admin

import (
	"context"
	"fmt"

	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/deploy"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type dev struct {
	CoreImage   string `name:"core" short:"c" help:"core docker image to deploy"`
	WorkerImage string `name:"worker" short:"w" help:"worker docker image to deploy"`
}

func (d *dev) Run(ctx context.Context, deployer deploy.DockerDeployer, logger log.FLogger) error {
	logger.Info("Deploying funless locally...\n")

	_ = logger.StartSpinner("Setting things up...")

	if d.CoreImage == "" {
		d.CoreImage = pkg.FLCore
	}

	if d.WorkerImage == "" {
		d.WorkerImage = pkg.FLWorker
	}

	if err := deployer.Setup(ctx); err != nil {
		return logger.StopSpinner(err)
	}

	if err := logger.StopSpinner(deployer.CreateFLNetworks(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner(fmt.Sprintf("pulling Core image (%s) 📦", d.CoreImage))
	if err := logger.StopSpinner(deployer.PullCoreImage(ctx, d.CoreImage)); err != nil {
		return err
	}

	_ = logger.StartSpinner(fmt.Sprintf("pulling Worker image (%s) 🗃", d.WorkerImage))
	if err := logger.StopSpinner(deployer.PullWorkerImage(ctx, d.WorkerImage)); err != nil {
		return err
	}

	_ = logger.StartSpinner("starting Core container 🎛️")
	if err := logger.StopSpinner(deployer.StartCore(ctx, d.CoreImage)); err != nil {
		return err
	}

	_ = logger.StartSpinner("starting Worker container 👷")
	if err := logger.StopSpinner(deployer.StartWorker(ctx, d.WorkerImage)); err != nil {
		return err
	}

	logger.Info("\nDeployment complete!")
	logger.Info("You can now start using Funless! 🎉")

	return nil
}
