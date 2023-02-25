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
	"os"
	"path/filepath"

	apiAppsV1 "k8s.io/api/apps/v1"
	apiBatchV1 "k8s.io/api/batch/v1"
	apiCoreV1 "k8s.io/api/core/v1"
	apiRbacV1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesDeployer interface {
	WithConfig(config string) error

	CreateNamespace(ctx context.Context) error
	CreateSvcAccount(ctx context.Context) error
	CreateRole(ctx context.Context) error
	CreateRoleBinding(ctx context.Context) error
	CreatePrometheusConfigMap(ctx context.Context) error
	DeployPrometheus(ctx context.Context) error
	DeployPrometheusService(ctx context.Context) error
	DeployPostgres(ctx context.Context) error
	DeployPostgresService(ctx context.Context) error
	StartInitPostgres(ctx context.Context) error
	DeployCore(ctx context.Context) error
	DeployCoreService(ctx context.Context) error
	DeployWorker(ctx context.Context) error
}

type FLKubernetesDeployer struct {
	kubernetesClientSet kubernetes.Interface

	namespace string
}

func getYAMLContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	yml, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return yml, nil
}

func NewKubernetesDeployer() KubernetesDeployer {
	return &FLKubernetesDeployer{namespace: "fl"}
}

func (k *FLKubernetesDeployer) WithConfig(config string) error {
	if config == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		config = filepath.Join(home, ".kube", "config")
	}

	kConfig, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		return err
	}

	clientSet, err := kubernetes.NewForConfig(kConfig)
	if err != nil {
		return err
	}

	k.kubernetesClientSet = clientSet
	return nil
}

func (k *FLKubernetesDeployer) CreateNamespace(ctx context.Context) error {
	yml, err := getYAMLContent("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/namespace.yml")
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "Namespace", APIVersion: "v1"}
	obj, err := ParseKubernetesYAML(yml, &apiCoreV1.Namespace{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	ns := obj.(*apiCoreV1.Namespace)

	_, err = k.kubernetesClientSet.CoreV1().Namespaces().Create(ctx, ns, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) CreateSvcAccount(ctx context.Context) error {
	yml, err := getYAMLContent("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/svc-account.yml")
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "ServiceAccount", APIVersion: "v1"}
	obj, err := ParseKubernetesYAML(yml, &apiCoreV1.ServiceAccount{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	svc := obj.(*apiCoreV1.ServiceAccount)

	_, err = k.kubernetesClientSet.CoreV1().ServiceAccounts(k.namespace).Create(ctx, svc, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) CreateRole(ctx context.Context) error {
	yml, err := getYAMLContent("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/svc-account.yml")
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "Role", APIVersion: "rbac.authorization.k8s.io/v1"}
	obj, err := ParseKubernetesYAML(yml, &apiRbacV1.Role{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	role := obj.(*apiRbacV1.Role)

	_, err = k.kubernetesClientSet.RbacV1().Roles(k.namespace).Create(ctx, role, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) CreateRoleBinding(ctx context.Context) error {
	yml, err := getYAMLContent("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/svc-account.yml")
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "RoleBinding", APIVersion: "rbac.authorization.k8s.io/v1"}
	obj, err := ParseKubernetesYAML(yml, &apiRbacV1.RoleBinding{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	roleBind := obj.(*apiRbacV1.RoleBinding)

	_, err = k.kubernetesClientSet.RbacV1().RoleBindings(k.namespace).Create(ctx, roleBind, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) CreatePrometheusConfigMap(ctx context.Context) error {
	yml, err := getYAMLContent("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/prometheus-cm.yml")
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"}
	obj, err := ParseKubernetesYAML(yml, &apiCoreV1.ConfigMap{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	configMap := obj.(*apiCoreV1.ConfigMap)

	_, err = k.kubernetesClientSet.CoreV1().ConfigMaps(k.namespace).Create(ctx, configMap, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) DeployPrometheus(ctx context.Context) error {
	yml, err := getYAMLContent("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/prometheus.yml")
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"}
	obj, err := ParseKubernetesYAML(yml, &apiAppsV1.Deployment{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	deployment := obj.(*apiAppsV1.Deployment)

	_, err = k.kubernetesClientSet.AppsV1().Deployments(k.namespace).Create(ctx, deployment, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) DeployPrometheusService(ctx context.Context) error {
	yml, err := getYAMLContent("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/prometheus.yml")
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "Service", APIVersion: "v1"}
	obj, err := ParseKubernetesYAML(yml, &apiCoreV1.Service{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	service := obj.(*apiCoreV1.Service)

	_, err = k.kubernetesClientSet.CoreV1().Services(k.namespace).Create(ctx, service, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) DeployPostgres(ctx context.Context) error {
	yml, err := getYAMLContent("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/postgres.yml")
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"}
	obj, err := ParseKubernetesYAML(yml, &apiAppsV1.Deployment{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	deployment := obj.(*apiAppsV1.Deployment)

	_, err = k.kubernetesClientSet.AppsV1().Deployments(k.namespace).Create(ctx, deployment, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) DeployPostgresService(ctx context.Context) error {
	yml, err := getYAMLContent("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/postgres.yml")
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "Service", APIVersion: "v1"}
	obj, err := ParseKubernetesYAML(yml, &apiCoreV1.Service{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	service := obj.(*apiCoreV1.Service)

	_, err = k.kubernetesClientSet.CoreV1().Services(k.namespace).Create(ctx, service, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) StartInitPostgres(ctx context.Context) error {
	yml, err := getYAMLContent("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/postgres.yml")
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "Job", APIVersion: "batch/v1"}
	obj, err := ParseKubernetesYAML(yml, &apiBatchV1.Job{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	job := obj.(*apiBatchV1.Job)

	_, err = k.kubernetesClientSet.BatchV1().Jobs(k.namespace).Create(ctx, job, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) DeployCore(ctx context.Context) error {
	yml, err := getYAMLContent("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/core.yml")
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"}
	obj, err := ParseKubernetesYAML(yml, &apiAppsV1.Deployment{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	deployment := obj.(*apiAppsV1.Deployment)

	_, err = k.kubernetesClientSet.AppsV1().Deployments(k.namespace).Create(ctx, deployment, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) DeployCoreService(ctx context.Context) error {
	yml, err := getYAMLContent("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/core.yml")
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "Service", APIVersion: "v1"}
	obj, err := ParseKubernetesYAML(yml, &apiCoreV1.Service{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	service := obj.(*apiCoreV1.Service)

	_, err = k.kubernetesClientSet.CoreV1().Services(k.namespace).Create(ctx, service, v1.CreateOptions{})

	return err
}

func (k *FLKubernetesDeployer) DeployWorker(ctx context.Context) error {
	yml, err := getYAMLContent("https://raw.githubusercontent.com/funlessdev/fl-deploy/main/kind/worker.yml")
	if err != nil {
		return err
	}

	typeMeta := v1.TypeMeta{Kind: "DaemonSet", APIVersion: "apps/v1"}
	obj, err := ParseKubernetesYAML(yml, &apiAppsV1.DaemonSet{TypeMeta: typeMeta})
	if err != nil {
		return err
	}

	daemonSet := obj.(*apiAppsV1.DaemonSet)

	_, err = k.kubernetesClientSet.AppsV1().DaemonSets(k.namespace).Create(ctx, daemonSet, v1.CreateOptions{})

	return err
}
