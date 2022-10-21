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
)

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
		return errors.New("No corresponding builder image found for the given language")
	}
	b.builderImg = image

	containerName, exists := pkg.FLRuntimeNames[language]
	if !exists {
		return errors.New("No corresponding container name found for the given language")
	}
	b.builderContainerName = containerName

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return err
	}
	b.client = cli

	outPath, _ := filepath.Abs(b.outPath)
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

	containerConfig := &container.Config{
		Image:   b.builderImg,
		Volumes: map[string]struct{}{},
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Source:   absPath,
				Target:   "/lib_fl/",
				ReadOnly: true,
				Type:     mount.TypeBind,
			},
			{
				Source: b.outPath,
				Target: "/out_wasm",
				Type:   mount.TypeBind,
			},
		},
		AutoRemove: true,
	}
	return runContainer(ctx, b.client, hostConfig, containerConfig, b.builderContainerName)
}

func (b *WasmBuilder) PullBuilderImage(ctx context.Context) error {
	return pullImage(ctx, b.client, b.builderImg)
}
