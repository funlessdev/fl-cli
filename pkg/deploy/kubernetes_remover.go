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

type KubernetesRemover interface {
	WithClientSet(cs kubernetes.Clientset)

	RemoveNamespace(ctx context.Context)
	RemoveSvcAccount(ctx context.Context)
	RemoveRole(ctx context.Context)
	RemoveRoleBinding(ctx context.Context)
	RemoveCore(ctx context.Context)
	RemoveWorker(ctx context.Context)
	RemovePrometheus(ctx context.Context)
}

type FlKubernetesRemover struct {
	kubernetesClientSet kubernetes.Clientset

	namespace       string
	svcAccountName  string
	roleName        string
	roleBindingName string

	coreName       string
	workerName     string
	prometheusName string
}

func NewKubernetesRemover() KubernetesRemover {
	return &FlKubernetesRemover{}
}

func (k *FlKubernetesRemover) WithClientSet(cs kubernetes.Clientset) {

}

func (k *FlKubernetesRemover) RemoveNamespace(ctx context.Context) {

}

func (k *FlKubernetesRemover) RemoveSvcAccount(ctx context.Context) {

}

func (k *FlKubernetesRemover) RemoveRole(ctx context.Context) {

}

func (k *FlKubernetesRemover) RemoveRoleBinding(ctx context.Context) {

}

func (k *FlKubernetesRemover) RemoveCore(ctx context.Context) {

}

func (k *FlKubernetesRemover) RemoveWorker(ctx context.Context) {

}

func (k *FlKubernetesRemover) RemovePrometheus(ctx context.Context) {

}
