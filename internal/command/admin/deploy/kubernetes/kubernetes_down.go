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

package admin_deploy_kubernetes

import (
	"context"

	"github.com/funlessdev/fl-cli/pkg/deploy"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Down struct {
	KubeConfig string `name:"kubeconfig" short:"k" help:"Absolute path to the kubeconfig file"`
}

func (f *Down) Help() string {
	return `
DESCRIPTION

	It removes a Kubernetes FunLess deployment.
	The "--kubeconfig" flag can be used to specify the absolute path 
	to the kubeconfig file.

EXAMPLES

	$ fl admin deploy kubernetes down --kubeconfig <your-kubeconfig-path>`
}

func (k *Down) Run(ctx context.Context, remover deploy.KubernetesRemover, logger log.FLogger) error {
	logger.Info("Removing Kubernetes FunLess deployment...\n\n")

	_ = logger.StartSpinner("Setting things up...")
	if err := logger.StopSpinner(remover.WithConfig(k.KubeConfig)); err != nil {
		return err
	}

	_ = logger.StartSpinner("Removing Namespace...")
	if err := logger.StopSpinner(remover.RemoveNamespace(ctx)); err != nil {
		return err
	}

	logger.Info("\nAll clear!\n")

	return nil
}
