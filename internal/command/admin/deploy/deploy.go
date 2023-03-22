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
	Docker     deploy_docker     `cmd:"" name:"docker" aliases:"d" help:"deploy locally with 1 core and 1 worker docker containers"`
	Kubernetes deploy_kubernetes `cmd:"" name:"kubernetes" aliases:"k,k8s" help:"deploy on an existing kubernetes cluster"`
}

type deploy_docker struct {
	Up   docker.Up   `cmd:"" name:"up" aliases:"u" help:"spin up Docker-based FunLess deployment"`
	Down docker.Down `cmd:"" name:"down" aliases:"d" help:"tear down Docker-based FunLess deployment"`
}

type deploy_kubernetes struct {
	Up   kubernetes.Up   `cmd:"" name:"up" aliases:"u" help:"spin up Kubernetes-based FunLess deployment"`
	Down kubernetes.Down `cmd:"" name:"down" aliases:"d" help:"tear down Kubernetes-based FunLess deployment"`
}

func (f *Deploy) Help() string {
	return `Below is a description of the architecture.

## CORE
	
At the heart of the platform there is the Core component, which is the orchestrator of the platform.
It manages functions and modules using a Postgres database behind.
When an invocation request arrives, it acts as a scheduler to pick one of the available Worker componets,
and it then sends the request to the worker to invoke the function (with the wasm binary if missing).
	
The core code resides in the apps/core folder, it is a Phoenix application that exposes a json REST API.
	
## WORKER

The Worker is the actual functions executor. 
It makes use of Rust NIFs to run user-defined functions via the wasmtime runtime. 
The worker makes use of a cache to avoid requesting the same function multiple times, 
and it is able to run multiple functions in parallel. 
When the core sends an invocation request, the worker first tries to find the function in the cache, 
if it is not present it requests back to the core the wasm binary. 
Then it executes the function and sends back the result to the core.

The worker code resides in the apps/worker folder.

## PROMETHEUS

We are using Prometheus to collect metrics from the platform. 
Besides collecting metrics from the Core and Worker, 
it is used by the Core to access the metrics of the Workers to make scheduling decisions.

## POSTGRES

We are using Postgres as the platform database, used by the Core to store functions and modules.
`
}

func (f *deploy_docker) Help() string {
	return "Group of commands for managing a local docker deployment with 1 core and 1 worker"
}
