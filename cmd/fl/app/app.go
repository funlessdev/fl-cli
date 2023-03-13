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

package app

import (
	"context"
	"net/http"
	"time"

	"github.com/alecthomas/kong"
	"github.com/funlessdev/fl-cli/internal/command/admin"
	"github.com/funlessdev/fl-cli/internal/command/fn"
	"github.com/funlessdev/fl-cli/internal/command/mod"
	"github.com/funlessdev/fl-cli/internal/command/template"
	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/build"
	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/deploy"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type CLI struct {
	Fn       fn.Fn             `cmd:"" help:"create, delete and manage functions"`
	Mod      mod.Mod           `cmd:"" help:"create, delete and manage modules"`
	Admin    admin.Admin       `cmd:"" aliases:"a" help:"deploy and manage the platform"`
	Template template.Template `cmd:"" help:"pull function templates"`

	Version kong.VersionFlag `short:"v" cmd:"" passthrough:"" help:"show fl version"`
}

func ParseCMD(version string) (*kong.Context, error) {
	cli := CLI{}
	ctx := context.Background()

	logger, err := buildLogger()
	dockerShell := buildDockerShell()

	kubernetesDeployer := deploy.NewKubernetesDeployer()
	kubernetesRemover := deploy.NewKubernetesRemover()

	wasmBuilder := build.NewWasmBuilder()

	if err != nil {
		return nil, err
	}

	flConfig := client.Config{Host: "http://localhost:4000"}
	flClient, err := client.NewClient(http.DefaultClient, flConfig)
	validator := client.InputValidator{}
	if err != nil {
		return nil, err
	}
	fnSvc := &client.FnService{Client: flClient, InputValidatorHandler: &validator}
	modSvc := &client.ModService{Client: flClient, InputValidatorHandler: &validator}
	userSvc := &client.UserService{Client: flClient}

	kong_ctx := kong.Parse(&cli,
		kong.Name("fl"),
		kong.Description("Funless CLI - fl"),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:             true,
			NoExpandSubcommands: true,
			Summary:             false,
			FlagsLast:           true,
		}),
		kong.BindTo(ctx, (*context.Context)(nil)),
		kong.BindTo(fnSvc, (*client.FnHandler)(nil)),
		kong.BindTo(modSvc, (*client.ModHandler)(nil)),
		kong.BindTo(logger, (*log.FLogger)(nil)),
		kong.BindTo(dockerShell, (*deploy.DockerShell)(nil)),
		kong.BindTo(kubernetesDeployer, (*deploy.KubernetesDeployer)(nil)),
		kong.BindTo(kubernetesRemover, (*deploy.KubernetesRemover)(nil)),
		kong.BindTo(wasmBuilder, (*build.DockerBuilder)(nil)),
		kong.BindTo(userSvc, (*client.UserHandler)(nil)),
		kong.Vars{
			"version":              version,
			"default_core_image":   pkg.CoreImg,
			"default_worker_image": pkg.WorkerImg,
		},
		kong.UsageOnError(),
	)
	return kong_ctx, nil
}

func Run(kong_ctx *kong.Context) error {
	return kong_ctx.Run()
}

func buildLogger() (log.FLogger, error) {
	b := log.NewLoggerBuilder()
	logger, err := b.WithDebug(true).SpinnerFrequency(150 * time.Millisecond).SpinnerCharSet(59).Build()
	return logger, err
}

func buildDockerShell() deploy.DockerShell {
	return &deploy.FLDockerShell{}
}
