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

package deploy

import (
	"context"

	"k8s.io/client-go/kubernetes"
)

type KubernetesDeployer interface {
	WithClientSet(cs kubernetes.Clientset)
	WithImages(coreImg, workerImg, promImg string)

	CreateNamespace(ctx context.Context)
	CreateSvcAccount(ctx context.Context)
	CreateRole(ctx context.Context)
	CreateRoleBinding(ctx context.Context)
	DeployCore(ctx context.Context)
	DeployWorker(ctx context.Context)
	DeployPrometheus(ctx context.Context)
}

type FLKubernetesDeployer struct {
	kubernetesClientSet kubernetes.Clientset

	namespace       string
	svcAccountName  string
	roleName        string
	roleBindingName string

	coreName       string
	workerName     string
	prometheusName string
}

func NewKubernetesDeployer() KubernetesDeployer {
	return &FLKubernetesDeployer{}
}

func (k *FLKubernetesDeployer) WithClientSet(cs kubernetes.Clientset) {

}

func (k *FLKubernetesDeployer) WithImages(coreImg, workerImg, promImg string) {

}

func (k *FLKubernetesDeployer) CreateNamespace(ctx context.Context) {

}

func (k *FLKubernetesDeployer) CreateSvcAccount(ctx context.Context) {

}

func (k *FLKubernetesDeployer) CreateRole(ctx context.Context) {

}

func (k *FLKubernetesDeployer) CreateRoleBinding(ctx context.Context) {

}

func (k *FLKubernetesDeployer) DeployCore(ctx context.Context) {

}

func (k *FLKubernetesDeployer) DeployWorker(ctx context.Context) {

}

func (k *FLKubernetesDeployer) DeployPrometheus(ctx context.Context) {

}
