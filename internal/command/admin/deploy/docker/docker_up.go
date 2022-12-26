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
	"io"
	"net/http"
	"strings"

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

func (u *Up) Run(ctx context.Context, dk deploy.DockerShell, logger log.FLogger) error {
	logger.Info("Deploying FunLess locally...\n")

	_ = logger.StartSpinner("Setting things up...")

	composeFilePath, err := downloadDockerCompose()
	if err != nil {
		return logger.StopSpinner(err)
	}

	// if another core image is specified, we have to replace it in the compose file
	if err := replaceImages(u.CoreImage, u.WorkerImage); err != nil {
		return logger.StopSpinner(err)
	}

	if err := downloadPrometheusConfig(); err != nil {
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

func downloadDockerCompose() (string, error) {
	// Check if it's already present
	if _, path, err := homedir.ReadFromConfigDir("docker-compose.yml"); err == nil {
		return path, nil
	}

	// Download docker-compose.yml
	resp, err := http.Get(dockerComposeYmlUrl)
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return homedir.WriteToConfigDir("docker-compose.yml", content, true)
}

func downloadPrometheusConfig() error {
	// Check if it's already present
	if _, _, err := homedir.ReadFromConfigDir("prometheus/config.yml"); err == nil {
		return nil
	}

	// Download prometheus/config.yml
	resp, err := http.Get(prometheusConfigYmlUrl)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if _, err := homedir.CreateDirInConfigDir("prometheus"); err != nil {
		return err
	}

	_, err = homedir.WriteToConfigDir("prometheus/config.yml", content, true)
	return err
}

func replaceImages(core string, worker string) error {
	if core == pkg.CoreImg && worker == pkg.WorkerImg {
		return nil
	}

	content, _, err := homedir.ReadFromConfigDir("docker-compose.yml")
	if err != nil {
		return errors.New("unable to read docker-compose.yml")
	}

	newCompose := string(content)
	if core != pkg.CoreImg {
		newCompose = strings.Replace(string(content), pkg.CoreImg, core, 1)
	}
	if worker != pkg.WorkerImg {
		newCompose = strings.Replace(newCompose, pkg.WorkerImg, worker, 1)
	}

	_, err = homedir.WriteToConfigDir("docker-compose.yml", []byte(newCompose), true)
	return err
}
