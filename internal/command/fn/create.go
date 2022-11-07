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
	"path"
	"strings"

	"github.com/funlessdev/fl-cli/pkg/build"
	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Create struct {
	Name       string `arg:"" name:"name" help:"name of the function to create"`
	Namespace  string `name:"namespace" short:"n" default:"_" help:"namespace of the function to create"`
	SourceDir  string `name:"source-dir" group:"build" short:"d" required:"" xor:"dir-file,dir-build" type:"existingdir" help:"path of the source directory"`
	SourceFile string `name:"source-file" short:"f" required:"" xor:"dir-file" type:"existingFile" help:"path of the source file"`
	OutDir     string `name:"out-dir" short:"o" xor:"out-build" type:"existingdir" help:"path where the compiled code file will be saved"`
	NoBuild    bool   `name:"no-build" short:"b" xor:"dir-build,out-build" help:"upload the file as-is, without building it"`
	Language   string `name:"language" group:"build" short:"l" enum:"js,rust" help:"programming language of the function"`
}

func (f *Create) Run(ctx context.Context, builder build.DockerBuilder, fnHandler client.FnHandler, logger log.FLogger) error {
	code, err := f.obtainCode(ctx, builder, logger)
	if err != nil {
		return err
	}

	res, err := fnHandler.Create(ctx, f.Name, f.Namespace, code)
	if err != nil {
		return extractError(err)
	}

	logger.Info(*res.Result)
	return nil
}

func (f *Create) obtainCode(ctx context.Context, builder build.DockerBuilder, logger log.FLogger) (*os.File, error) {
	if f.SourceDir != "" {
		return f.buildFromDir(ctx, builder, logger)
	} else if f.NoBuild {
		return f.openWasmFile()
	} else {
		//NOTE: build single file => not implemented
		return nil, errors.New("building from a single file is not yet implemented")
	}
}

func (f *Create) buildFromDir(ctx context.Context, builder build.DockerBuilder, logger log.FLogger) (*os.File, error) {
	if f.OutDir == "" {
		/* can't use default, as outDir is also in a xor group */
		f.OutDir = "./out_wasm/"
	}
	logger.Info("Building the given function using fl-runtimes...\n")

	_ = logger.StartSpinner("Setting up...")
	if build_err := logger.StopSpinner(builder.Setup(ctx, f.Language, f.OutDir)); build_err != nil {
		return nil, build_err
	}

	_ = logger.StartSpinner(fmt.Sprintf("Pulling builder image for %s 📦", f.Language))
	if build_err := logger.StopSpinner(builder.PullBuilderImage(ctx)); build_err != nil {
		return nil, build_err
	}
	_ = logger.StartSpinner("Building source using builder image 🛠️")
	if build_err := logger.StopSpinner(builder.BuildSource(ctx, f.SourceDir)); build_err != nil {
		return nil, build_err
	}

	return os.Open(path.Join(f.OutDir, "./code.wasm"))
}

func (f *Create) openWasmFile() (*os.File, error) {
	if !strings.HasSuffix(f.SourceFile, ".wasm") {
		return nil, errors.New("a file with the .wasm extension must be passed")
	}

	code, err := os.Open(f.SourceFile)
	if err != nil {
		return nil, err
	}
	stat, err := code.Stat()

	if err != nil {
		return nil, err
	}
	if stat.Size() == 0 {
		return nil, errors.New("passing an empty file as source")
	}
	return code, nil
}
