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

package admin_deploy

import (
	docker "github.com/funlessdev/fl-cli/internal/command/admin/deploy/docker"
	kubernetes "github.com/funlessdev/fl-cli/internal/command/admin/deploy/kubernetes"
)

type Deploy struct {
	Docker     deploy_docker     `cmd:"" name:"docker" aliases:"d" help:"Deploy locally with 1 core and 1 worker docker containers"`
	Kubernetes deploy_kubernetes `cmd:"" name:"kubernetes" aliases:"k,k8s" help:"Deploy on an existing kubernetes cluster"`
}

type deploy_docker struct {
	Up   docker.Up   `cmd:"" name:"up" aliases:"u" help:"Spin up Docker-based FunLess deployment"`
	Down docker.Down `cmd:"" name:"down" aliases:"d" help:"Tear down Docker-based FunLess deployment"`
}

type deploy_kubernetes struct {
	Up   kubernetes.Up   `cmd:"" name:"up" aliases:"u" help:"Spin up Kubernetes-based FunLess deployment"`
	Down kubernetes.Down `cmd:"" name:"down" aliases:"d" help:"Tear down Kubernetes-based FunLess deployment"`
}

func (f *Deploy) Help() string {
	return `
The deployment consists of a few main services:

	FunLess Core: the orchestrator of the platform;
	FunLess Worker(s): the functions executor;
	Prometheus: metrics collector (also used for scheduling by the Core);
	PostgreSQL DB: the database for functions/modules/users.

`
}
