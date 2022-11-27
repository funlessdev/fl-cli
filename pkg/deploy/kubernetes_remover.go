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
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesRemover interface {
	WithConfig(config string) error

	RemoveNamespace(ctx context.Context) error
}

type FLKubernetesRemover struct {
	kubernetesClientSet kubernetes.Interface

	namespace string
}

func NewKubernetesRemover() KubernetesRemover {
	return &FLKubernetesRemover{namespace: "fl"}
}

func (k *FLKubernetesRemover) WithConfig(config string) error {
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

func (k *FLKubernetesRemover) RemoveNamespace(ctx context.Context) error {
	selector := fmt.Sprintf("metadata.name=%s", k.namespace)
	watcher, err := k.kubernetesClientSet.CoreV1().Namespaces().Watch(ctx, v1.ListOptions{FieldSelector: selector})
	if err != nil {
		return err
	}

	err = k.kubernetesClientSet.CoreV1().Namespaces().Delete(ctx, k.namespace, v1.DeleteOptions{})
	if err != nil {
		return err
	}

	for {
		event := <-watcher.ResultChan()
		if event.Type == watch.Deleted {
			break
		}
	}
	return err
}
