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

package admin_deploy_docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/deploy"
	"github.com/funlessdev/fl-cli/pkg/docker"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Up struct {
	CoreImage   string `name:"core" short:"c" help:"core docker image to deploy" default:"${default_core_image}"`
	WorkerImage string `name:"worker" short:"w" help:"worker docker image to deploy" default:"${default_worker_image}"`
}

func (d *Up) Run(ctx context.Context, deployer deploy.DockerDeployer, logger log.FLogger) error {
	logger.Info("Deploying FunLess locally...\n")

	_ = logger.StartSpinner("Setting things up...")

	if err := setupDev(d.CoreImage, d.WorkerImage, deployer); err != nil {
		return logger.StopSpinner(err)
	}

	if err := logger.StopSpinner(deployer.CreateFLNetwork(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner(fmt.Sprintf("pulling Core image (%s) 🐋", d.CoreImage))
	if err := logger.StopSpinner(deployer.PullCoreImage(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner(fmt.Sprintf("pulling Worker image (%s) 🐋", d.WorkerImage))
	if err := logger.StopSpinner(deployer.PullWorkerImage(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("pulling Prometheus image 🐋")
	if err := logger.StopSpinner(deployer.PullPromImage(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("starting Core container 🎛️")

	if err := logger.StopSpinner(deployer.StartCore(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("starting Worker container 👷")
	if err := logger.StopSpinner(deployer.StartWorker(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("starting Prometheus container 📊")
	if err := logger.StopSpinner(deployer.StartProm(ctx)); err != nil {
		return err
	}

	logger.Info("\nDeployment complete!")
	logger.Info("You can now start using FunLess! 🎉")

	return nil
}

func setupDev(core string, worker string, deployer deploy.DockerDeployer) error {
	deployer.WithImages(core, worker)

	c, err := setupDockerClient()
	if err != nil {
		return err
	}

	deployer.WithDockerClient(c)
	return deployer.WithLogs(pkg.LocalLogsPath)
}

func setupDockerClient() (docker.DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return docker.DockerClient{}, err
	}
	return docker.NewDockerClient(cli), nil
}
