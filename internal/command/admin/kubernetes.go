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
	"os"
	"path/filepath"

	"github.com/funlessdev/fl-cli/pkg/deploy"
	"github.com/funlessdev/fl-cli/pkg/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type k8s struct {
	CoreImage       string `name:"core" short:"c" help:"core docker image to deploy" default:"${default_core_image}"`
	WorkerImage     string `name:"worker" short:"w" help:"worker docker image to deploy" default:"${default_worker_image}"`
	PrometheusImage string `name:"prometheus" short:"p" help:"prometheus docker image to deploy" default:"${default_prometheus_image}"`
	KubeConfig      string `name:"kubeconfig" short:"k" help:"absolute path to the kubeconfig file"`
}

func (k *k8s) Run(ctx context.Context, deployer deploy.DockerDeployer, logger log.FLogger) error {
	logger.Info("Deploying FunLess on Kubernetes...\n")

	_ = logger.StartSpinner("Setting things up...")

	_, err := setupKubernetesClientSet(k.KubeConfig)

	if err != nil {
		return err
	}

	logger.Info("\nDeployment complete!")
	logger.Info("You can now start using FunLess! 🎉")

	return nil
}

func setupDeployer() {

}

func setupKubernetesClientSet(config string) (kubernetes.Interface, error) {
	if config == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		config = filepath.Join(home, ".kube", "config")
	}

	kConfig, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(kConfig)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}
