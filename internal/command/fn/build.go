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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/build"
	"github.com/funlessdev/fl-cli/pkg/docker"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Build struct {
	Name        string `arg:"" help:"The name of the function"`
	Source      string `arg:"" type:"existingdir" help:"Path of the source directory"`
	Destination string `short:"d" type:"path" help:"Path where the compiled wasm file will be saved" default:"."`
	Language    string `short:"l" enum:"rust,js" required:"" help:"Programming language of the function"`
}

func (c *Build) Help() string {
	return `
DESCRIPTION

	It creates wasm for function specified in source.
	It must be use the flag "--language" to specify the language of the funcion.
	The possible value is one of from the following list.

		[rust, js]

	The "--destination" flag can be used to choose a destination directory other than the default one. 

EXAMPLES
	
	$ fl fn build <your-function-name> <your-function-source> --language=<lang-from-enum> --destination=<your-destination-directory>
`

}
func (b *Build) Run(ctx context.Context, builder build.DockerBuilder, logger log.FLogger) error {
	logger.Info(fmt.Sprintf("Building %s into a wasm binary...\n\n", b.Name))

	_ = logger.StartSpinner("Setting up...")
	if err := logger.StopSpinner(setupBuilder(builder, b.Language, b.Destination)); err != nil {
		return err
	}

	_ = logger.StartSpinner(fmt.Sprintf("Checking directory %s files for language %s... 🔍", b.Source, b.Language))
	if err := logger.StopSpinner(checkMustContainFiles(b.Language, b.Source)); err != nil {
		return err
	}

	_ = logger.StartSpinner(fmt.Sprintf("Pulling %s builder image (%s) 📦", pkg.SupportedLanguages[b.Language].Name, pkg.SupportedLanguages[b.Language].BuilderImage))
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

	logger.Info(fmt.Sprintf("\nSuccessfully built %s.wasm 🥳🥳\n", b.Name))
	return nil
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

func checkMustContainFiles(lang, source string) error {
	language, ok := pkg.SupportedLanguages[lang]
	if !ok {
		return errors.New("unsupported language")
	}
	for _, f := range language.MustContainFiles {
		path := filepath.Join(source, f)
		_, err := os.Stat(path)
		if err != nil && os.IsNotExist(err) {
			return fmt.Errorf("necessary file %s not found in path %s", f, source)
		}
	}
	return nil
}
