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
	"io"
	"net/http"

	"github.com/funlessdev/fl-cli/pkg/fl_k8s"
	apiCoreV1 "k8s.io/api/core/v1"
	apiRbacV1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubernetesDeployer interface {
	WithClientSet(cs kubernetes.Interface)

	CreateNamespace(ctx context.Context) error
	CreateSvcAccount(ctx context.Context) error
	CreateRole(ctx context.Context) error
	CreateRoleBinding(ctx context.Context) error
	DeployCore(ctx context.Context) error
	DeployWorker(ctx context.Context) error
	DeployPrometheus(ctx context.Context) error
}

type FLKubernetesDeployer struct {
	kubernetesClientSet kubernetes.Interface

	namespace string
}

func NewKubernetesDeployer() KubernetesDeployer {
	return &FLKubernetesDeployer{namespace: "fl"}
}

func (k *FLKubernetesDeployer) WithClientSet(cs kubernetes.Interface) {
	k.kubernetesClientSet = cs
}

func (k *FLKubernetesDeployer) CreateNamespace(ctx context.Context) error {
	resp, err := http.Get("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/namespace.yml")
	if err != nil {
		return err
	}

	yml, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "Namespace", APIVersion: "v1"}
	obj, err := fl_k8s.ParseKubernetesYAML(yml, &apiCoreV1.Namespace{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	ns := obj.(*apiCoreV1.Namespace)

	_, err = k.kubernetesClientSet.CoreV1().Namespaces().Create(ctx, ns, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) CreateSvcAccount(ctx context.Context) error {
	resp, err := http.Get("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/svc-account.yml")
	if err != nil {
		return err
	}

	yml, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "ServiceAccount", APIVersion: "v1"}
	obj, err := fl_k8s.ParseKubernetesYAML(yml, &apiCoreV1.ServiceAccount{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	svc := obj.(*apiCoreV1.ServiceAccount)

	_, err = k.kubernetesClientSet.CoreV1().ServiceAccounts(k.namespace).Create(ctx, svc, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) CreateRole(ctx context.Context) error {
	resp, err := http.Get("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/svc-account.yml")
	if err != nil {
		return err
	}

	yml, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "Role", APIVersion: "rbac.authorization.k8s.io/v1"}
	obj, err := fl_k8s.ParseKubernetesYAML(yml, &apiRbacV1.Role{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	role := obj.(*apiRbacV1.Role)

	_, err = k.kubernetesClientSet.RbacV1().Roles(k.namespace).Create(ctx, role, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) CreateRoleBinding(ctx context.Context) error {
	resp, err := http.Get("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/svc-account.yml")
	if err != nil {
		return err
	}

	yml, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "RoleBinding", APIVersion: "rbac.authorization.k8s.io/v1"}
	obj, err := fl_k8s.ParseKubernetesYAML(yml, &apiRbacV1.RoleBinding{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	roleBind := obj.(*apiRbacV1.RoleBinding)

	_, err = k.kubernetesClientSet.RbacV1().RoleBindings(k.namespace).Create(ctx, roleBind, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) DeployCore(ctx context.Context) error {
	return nil
}

func (k *FLKubernetesDeployer) DeployWorker(ctx context.Context) error {
	return nil
}

func (k *FLKubernetesDeployer) DeployPrometheus(ctx context.Context) error {
	return nil
}
