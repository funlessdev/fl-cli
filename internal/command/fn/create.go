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
	"os"
	"path/filepath"

	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/build"
	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Create struct {
	Name     string `arg:"" help:"name of the function to create"`
	Source   string `arg:"" type:"existingdir" help:"path of the source directory"`
	Module   string `short:"m" default:"_" help:"module of the function to create"`
	Language string `short:"l" required:"" enum:"rust,js" help:"programming language of the function"`
}

func (c *Create) Run(ctx context.Context, builder build.DockerBuilder, fnHandler client.FnHandler, logger log.FLogger) error {
	logger.Infof("Creating %s function...\n\n", c.Name)

	_ = logger.StartSpinner("Building function...🏗 ️")
	dest, err := os.MkdirTemp("", "funless-bin")
	if err != nil {
		return logger.StopSpinner(err)
	}
	defer os.RemoveAll(dest)

	if err := setupBuilder(builder, c.Language, dest); err != nil {
		return logger.StopSpinner(err)
	}
	if err := checkMustContainFiles(c.Language, c.Source); err != nil {
		return logger.StopSpinner(err)
	}
	if err := builder.PullBuilderImage(ctx); err != nil {
		return logger.StopSpinner(err)
	}
	if err := builder.BuildSource(ctx, c.Source); err != nil {
		return logger.StopSpinner(err)
	}
	_ = logger.StopSpinner(nil)

	_ = logger.StartSpinner("Uploading function... 📮")
	code, err := openWasmFile(filepath.Join(dest, "code.wasm"))
	if err != nil {
		return logger.StopSpinner(err)
	}

	err = fnHandler.Create(ctx, c.Name, c.Module, code)
	if err != nil {
		return logger.StopSpinner(pkg.ExtractError(err))
	}
	_ = logger.StopSpinner(nil)

	logger.Info(fmt.Sprintf("\nSuccessfully created function %s/%s.\n", c.Module, c.Name))
	return nil
}
