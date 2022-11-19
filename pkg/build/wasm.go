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
	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/docker"
)

var builderNames = map[string]string{
	"js":   "fl-js-builder",
	"rust": "fl-rust-builder",
}

type DockerBuilder interface {
	Setup(client docker.DockerClient, language string, dest string) error
	PullBuilderImage(ctx context.Context) error
	BuildSource(ctx context.Context, srcPath string) error
	RenameCodeWasm(name string) error
}

type WasmBuilder struct {
	flDocker             docker.DockerClient
	builderImg           string
	builderContainerName string
	outPath              string
}

func NewWasmBuilder() DockerBuilder {
	return &WasmBuilder{}
}

func (b *WasmBuilder) Setup(flDocker docker.DockerClient, language string, dest string) error {
	dest, err := filepath.Abs(dest)
	if err != nil {
		return err
	}

	image, exists := pkg.FLBuilderImages[language]
	if !exists {
		return errors.New("no corresponding builder image found for the given language")
	}
	b.builderImg = image

	containerName, exists := builderNames[language]
	if !exists {
		return errors.New("no corresponding builder name found for the given language")
	}
	b.builderContainerName = containerName
	b.flDocker = flDocker
	b.outPath = dest

	err = os.MkdirAll(dest, 0700)
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

	configs := docker.ContainerConfigs{
		ContName:   b.builderContainerName,
		Container:  containerConfig,
		Host:       hostConfig,
		Networking: nil,
	}

	return b.flDocker.RunAndWait(ctx, configs)
}

func (b *WasmBuilder) PullBuilderImage(ctx context.Context) error {
	exists, err := b.flDocker.ImageExists(ctx, b.builderImg)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return b.flDocker.Pull(ctx, b.builderImg)
}

func (b *WasmBuilder) RenameCodeWasm(name string) error {
	return os.Rename(filepath.Join(b.outPath, "code.wasm"), filepath.Join(b.outPath, name+".wasm"))
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
