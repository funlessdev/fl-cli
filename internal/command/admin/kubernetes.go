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

	"github.com/funlessdev/fl-cli/pkg/deploy"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type k8s struct {
	KubeConfig string `name:"kubeconfig" short:"k" help:"absolute path to the kubeconfig file"`
}

func (k *k8s) Run(ctx context.Context, deployer deploy.KubernetesDeployer, logger log.FLogger) error {
	logger.Info("Deploying FunLess on Kubernetes...\n")

	_ = logger.StartSpinner("Setting things up...")
	if err := logger.StopSpinner(deployer.WithConfig(k.KubeConfig)); err != nil {
		return err
	}

	_ = logger.StartSpinner("Creating Namespace...")
	if err := logger.StopSpinner(deployer.CreateNamespace(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("Creating ServiceAccount...")
	if err := logger.StopSpinner(deployer.CreateSvcAccount(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("Creating Role...")
	if err := logger.StopSpinner(deployer.CreateRole(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("Creating RoleBinding...")
	if err := logger.StopSpinner(deployer.CreateRoleBinding(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("Creating Prometheus ConfigMap...")
	if err := logger.StopSpinner(deployer.CreatePrometheusConfigMap(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("Deploying Prometheus...")
	if err := logger.StopSpinner(deployer.DeployPrometheus(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("Deploying Prometheus Service...")
	if err := logger.StopSpinner(deployer.DeployPrometheusService(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("Deploying Core...")
	if err := logger.StopSpinner(deployer.DeployCore(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("Deploying Core Service...")
	if err := logger.StopSpinner(deployer.DeployCoreService(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("Deploying Workers...")
	if err := logger.StopSpinner(deployer.DeployWorker(ctx)); err != nil {
		return err
	}

	logger.Info("\nDeployment complete!")
	logger.Info("You can now start using FunLess! 🎉")

	return nil
}
