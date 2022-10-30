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

package build

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/docker_utils"
)

type DockerBuilder interface {
	Setup(ctx context.Context, language string, outDir string) error
	PullBuilderImage(ctx context.Context) error
	BuildSource(ctx context.Context, srcPath string) error
}

type WasmBuilder struct {
	client               *client.Client
	builderImg           string
	builderContainerName string
	outPath              string
}

var _ DockerBuilder = &WasmBuilder{}

func NewWasmBuilder() *WasmBuilder {
	return &WasmBuilder{}
}

func (b *WasmBuilder) Setup(ctx context.Context, language string, outDir string) error {
	image, exists := pkg.FLRuntimes[language]
	if !exists {
		return errors.New("no corresponding builder image found for the given language")
	}
	b.builderImg = image

	containerName, exists := pkg.FLRuntimeNames[language]
	if !exists {
		return errors.New("no corresponding container name found for the given language")
	}
	b.builderContainerName = containerName

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return err
	}
	b.client = cli

	outPath, _ := filepath.Abs(outDir)
	b.outPath = outPath

	err = os.MkdirAll(outPath, 0700)
	if err != nil {
		return err
	}

	return nil
}

func (b *WasmBuilder) BuildSource(ctx context.Context, srcPath string) error {
	absPath, err := filepath.Abs(srcPath)
	if err != nil {
		return err
	}

	containerConfig := builderContainerConfig(b.builderImg)
	hostConfig := builderHostConfig(absPath, b.outPath)

	configs := docker_utils.ContainerConfigs{
		ContName:   b.builderContainerName,
		Container:  containerConfig,
		Host:       hostConfig,
		Networking: nil,
	}

	return docker_utils.RunAndWaitContainer(ctx, b.client, configs)
}

func (b *WasmBuilder) PullBuilderImage(ctx context.Context) error {
	return docker_utils.PullImage(ctx, b.client, b.builderImg)
}

func builderContainerConfig(builderImg string) *container.Config {
	return &container.Config{
		Image:   builderImg,
		Volumes: map[string]struct{}{},
	}
}

func builderHostConfig(absPath, outPath string) *container.HostConfig {
	return &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Source:   absPath,
				Target:   "/lib_fl/",
				ReadOnly: true,
				Type:     mount.TypeBind,
			},
			{
				Source: outPath,
				Target: "/out_wasm",
				Type:   mount.TypeBind,
			},
		},
		AutoRemove: true,
	}
}
