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

package admin_deploy_docker

import (
	"context"
	"errors"
	"os"

	"github.com/funlessdev/fl-cli/pkg/deploy"
	"github.com/funlessdev/fl-cli/pkg/homedir"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Down struct{}

func (r *Down) Run(ctx context.Context, dk deploy.DockerShell, logger log.FLogger) error {
	logger.Info("Removing local FunLess deployment...\n")

	_, composeFilePath, err := homedir.ReadFromConfigDir("docker-compose.yml")
	if err != nil {
		return errors.New("unable to read docker-compose.yml file")
	}
	defer os.Remove(composeFilePath)

	err = dk.ComposeDown(composeFilePath)
	if err != nil {
		return err
	}

	logger.Info("\nAll clear! 👍")

	return nil
}
