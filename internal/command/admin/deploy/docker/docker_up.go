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
	"io"
	"net/http"
	"path/filepath"

	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/deploy"
	"github.com/funlessdev/fl-cli/pkg/homedir"
	"github.com/funlessdev/fl-cli/pkg/log"
)

const (
	dockerComposeYmlUrl    = "https://raw.githubusercontent.com/funlessdev/fl-deploy/main/docker-compose/docker-compose.yml"
	prometheusConfigYmlUrl = "https://raw.githubusercontent.com/funlessdev/fl-deploy/main/docker-compose/prometheus/config.yml"
)

type Up struct {
	CoreImage   string `name:"core" short:"c" help:"core docker image to deploy" default:"${default_core_image}"`
	WorkerImage string `name:"worker" short:"w" help:"worker docker image to deploy" default:"${default_worker_image}"`
}

func (d *Up) Run(ctx context.Context, dk deploy.DockerShell, logger log.FLogger) error {
	logger.Info("Deploying FunLess locally...\n")

	_ = logger.StartSpinner("Setting things up...")

	composeFilePath, err := getFileInConfigDir(dockerComposeYmlUrl, "docker-compose.yml")
	if err != nil {
		return logger.StopSpinner(err)
	}

	if _, err := getFileInConfigDir(prometheusConfigYmlUrl, "prometheus/config.yml"); err != nil {
		return logger.StopSpinner(err)
	}
	_ = logger.StopSpinner(nil)

	if err := dk.ComposeUp(composeFilePath); err != nil {
		return err
	}

	logger.Info("\nDeployment complete!")
	logger.Info("You can now start using FunLess! 🎉")

	return nil
}

var getFileInConfigDir = func(url string, file string) (string, error) {
	// Try to read from config dir
	_, path, err := homedir.ReadFromConfigDir(file)
	if err == nil {
		return path, nil
	}

	// if file doesn't exist or unreadable, download it
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// if we are in a sub folder, create it
	parentDir := filepath.Dir(file)
	if parentDir != pkg.ConfigDir {
		if _, err := homedir.CreateDirInConfigDir(parentDir); err != nil {
			return "", err
		}
	}

	return homedir.WriteToConfigDir(file, content, true)
}
