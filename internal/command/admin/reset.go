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

	"github.com/funlessdev/fl-cli/pkg/deploy"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type reset struct{}

func (r *reset) Run(ctx context.Context, deployer deploy.DevDeployer, logger log.FLogger) error {
	logger.Info("Removing local FunLess deployment...\n")

	if err := deployer.Setup(ctx, "", ""); err != nil {
		return err
	}

	_ = logger.StartSpinner("Removing Core container... ☠️")
	if err := logger.StopSpinner(deployer.RemoveCoreContainer(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("Removing Worker container... 🔪")
	if err := logger.StopSpinner(deployer.RemoveWorkerContainer(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("Removing Prometheus container... ⚰️")
	if err := logger.StopSpinner(deployer.RemovePromContainer(ctx)); err != nil {
		return err
	}

	_ = logger.StartSpinner("Removing fl network... ✂️")
	if err := logger.StopSpinner(deployer.RemoveFLNetwork(ctx)); err != nil {
		return err
	}

	logger.Info("\nAll clear! 👍")

	return nil
}
