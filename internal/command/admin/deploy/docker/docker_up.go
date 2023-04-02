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

	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/deploy"
	"github.com/funlessdev/fl-cli/pkg/homedir"
	"github.com/funlessdev/fl-cli/pkg/log"
	"gopkg.in/yaml.v2"
)

const (
	dockerComposeYmlUrl    = "https://raw.githubusercontent.com/funlessdev/fl-deploy/main/docker-compose/docker-compose.yml"
	envUrl                 = "https://raw.githubusercontent.com/funlessdev/fl-deploy/main/docker-compose/.env.example"
	prometheusConfigYmlUrl = "https://raw.githubusercontent.com/funlessdev/fl-deploy/main/docker-compose/prometheus/config.yml"
	filebeatComposeYmlUrl  = "https://raw.githubusercontent.com/funlessdev/fl-deploy/main/docker-compose/filebeat/filebeat.compose.yml"
)

type Up struct {
	CoreImage   string `name:"core" short:"c" help:"Core docker image to deploy" default:"${default_core_image}"`
	WorkerImage string `name:"worker" short:"w" help:"Worker docker image to deploy" default:"${default_worker_image}"`
}

func (f *Up) Help() string {
	return `
DESCRIPTION

	It creates a local Docker-based FunLess deployment.
	The "--core" and "--worker" flags can be used to choose a core 
	and a worker image other than the default ones.

EXAMPLES

	$ fl admin deploy docker up --core <your-core-image> --worker <your-worker-image>`
}

func (u *Up) Run(ctx context.Context, dk deploy.DockerShell, logger log.FLogger, config client.Config) error {
	logger.Info("Deploying FunLess locally...\n\n")

	cmdEnv := map[string]string{"SECRET_KEY_BASE": config.SecretKeyBase}
	ctx = context.WithValue(ctx, pkg.FLContextKey("env"), cmdEnv)

	_ = logger.StartSpinner("Setting things up...")

	composeFilePath, err := downloadFile("docker-compose.yml", dockerComposeYmlUrl)
	if err != nil {
		return logger.StopSpinner(err)
	}

	// if another core image is specified, we have to replace it in the compose file
	if err := replaceImages(u.CoreImage, u.WorkerImage); err != nil {
		return logger.StopSpinner(err)
	}

	// prometheus config file
	if err := downloadFolderFile("prometheus", "config.yml", prometheusConfigYmlUrl); err != nil {
		return logger.StopSpinner(err)
	}

	// filebeat compose file
	if err := downloadFolderFile("filebeat", "filebeat.compose.yml", filebeatComposeYmlUrl); err != nil {
		return logger.StopSpinner(err)
	}

	// .env file
	if _, err := downloadFile(".env", envUrl); err != nil {
		return logger.StopSpinner(err)
	}

	_ = logger.StopSpinner(nil)

	if err := dk.ComposeUp(ctx, composeFilePath); err != nil {
		return err
	}

	logger.Info("\nDeployment complete!\n")
	logger.Info("You can now start using FunLess! 🎉\n")

	return nil
}

func downloadFile(name, url string) (string, error) {
	// Check if already present
	if _, path, err := homedir.ReadFromConfigDir(name); err == nil {
		return path, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return homedir.WriteToConfigDir(name, content, true)

}

func downloadFolderFile(folder, file, url string) error {
	filepath := folder + "/" + file

	// Check if it's already present
	if _, _, err := homedir.ReadFromConfigDir(filepath); err == nil {
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if _, err := homedir.CreateDirInConfigDir(folder); err != nil {
		return err
	}

	_, err = homedir.WriteToConfigDir(filepath, content, true)
	return err
}

func replaceImages(core string, worker string) error {
	content, _, err := homedir.ReadFromConfigDir("docker-compose.yml")
	if err != nil {
		return errors.New("unable to read docker-compose.yml")
	}

	var composeYaml map[string]interface{}
	err = yaml.Unmarshal(content, &composeYaml)
	if err != nil {
		return err
	}

	svc := composeYaml["services"].(map[interface{}]interface{})
	svcCore := svc["core"].(map[interface{}]interface{})
	svcWorker := svc["worker"].(map[interface{}]interface{})

	svcCore["image"] = core
	svcWorker["image"] = worker

	svc["core"] = svcCore
	svc["worker"] = svcWorker
	composeYaml["services"] = svc

	newCompose, err := yaml.Marshal(composeYaml)

	if err != nil {
		return err
	}

	_, err = homedir.WriteToConfigDir("docker-compose.yml", []byte(newCompose), true)

	return err
}
