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

package fn

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/build"
	"github.com/funlessdev/fl-cli/pkg/docker"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Build struct {
	Name        string `arg:"" help:"the name of the function"`
	Source      string `arg:"" short:"s" required:"" xor:"dir-file,dir-build" type:"existingdir" help:"path of the source directory"`
	Destination string `short:"d" type:"path" help:"path where the compiled wasm file will be saved" default:"."`
	Language    string `short:"l" enum:"rust,js" required:"" help:"programming language of the function"`
}

func (b *Build) Run(ctx context.Context, builder build.DockerBuilder, logger log.FLogger) error {
	logger.Info(fmt.Sprintf("Building %s into a wasm binary...\n", b.Name))

	_ = logger.StartSpinner("Setting up...")
	if err := logger.StopSpinner(setupBuilder(builder, b.Language, b.Destination)); err != nil {
		return err
	}
	_ = logger.StartSpinner(fmt.Sprintf("Pulling %s builder image (%s) 📦", langNames[b.Language], pkg.FLBuilderImages[b.Language]))
	if err := logger.StopSpinner(builder.PullBuilderImage(ctx)); err != nil {
		return err
	}
	_ = logger.StartSpinner("Building source... 🛠️")
	if err := builder.BuildSource(ctx, b.Source); err != nil {
		return logger.StopSpinner(err)
	}
	if err := logger.StopSpinner(builder.RenameCodeWasm(b.Name)); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("\nSuccessfully built function at %s/%s.wasm 🥳🥳", b.Destination, b.Name))
	return nil
}

var langNames = map[string]string{
	"js":   "Javascript",
	"rust": "Rust",
}

func setupBuilder(builder build.DockerBuilder, lang, out string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return err
	}
	flDocker := docker.NewDockerClient(cli)
	if err != nil {
		return err
	}

	dest := filepath.Clean(out)
	if err = builder.Setup(flDocker, lang, dest); err != nil {
		return err
	}
	return nil
}
