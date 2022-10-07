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

import "context"

type DockerDeployer interface {
	Setup(ctx context.Context, coreImg, workerImg string) error

	CreateFLNetworks(ctx context.Context) error
	PullCoreImage(ctx context.Context) error
	PullWorkerImage(ctx context.Context) error
	StartCore(ctx context.Context) error
	StartWorker(ctx context.Context) error

	RemoveFLNetworks(ctx context.Context) error
	RemoveCoreContainer(context.Context) error
	RemoveWorkerContainer(context.Context) error
	RemoveFunctionContainers(ctx context.Context) error
}
