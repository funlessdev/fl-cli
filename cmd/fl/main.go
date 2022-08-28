// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/funlessdev/funless-cli/pkg/client"
	"github.com/funlessdev/funless-cli/pkg/command"
	"github.com/funlessdev/funless-cli/pkg/log"
)

// CLIVersion holds the current version, to be set by the build with go build -ldflags "-X main.FLVersion=<version>"
var FLVersion = "vX.Y.Z.build"

type CLI struct {
	Fn    command.Fn    `cmd:"" help:"todo fn subcommand help"`
	Admin command.Admin `cmd:"" help:"todo admin subcommand help"`

	Version kong.VersionFlag `short:"v" cmd:"" passthrough:"" help:"show fl version"`
}

func main() {
	cli := CLI{}
	ctx := context.Background()

	logger, err := buildLogger()

	writer := os.Stdout

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	flConfig := client.Config{Host: "http://localhost:4001"}
	flClient, err := client.NewClient(http.DefaultClient, flConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fnSvc := &client.FnService{Client: flClient}

	kong_ctx := kong.Parse(&cli,
		kong.Name("fl"),
		kong.Description("Funless CLI - fl"),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:             true,
			NoExpandSubcommands: true,
			Summary:             true,
			FlagsLast:           true,
		}),
		kong.BindTo(ctx, (*context.Context)(nil)),
		kong.BindTo(fnSvc, (*client.FnHandler)(nil)),
		kong.BindTo(logger, (*log.FLogger)(nil)),
		kong.BindTo(writer, (*io.Writer)(nil)),
		kong.Vars{
			"version": FLVersion,
		},
		kong.UsageOnError(),
	)

	kong_ctx.FatalIfErrorf(kong_ctx.Run())
}

func buildLogger() (log.FLogger, error) {
	b := log.NewLoggerBuilder()
	logger, err := b.WithDebug(true).SpinnerFrequency(150 * time.Millisecond).SpinnerCharSet(59).Build()
	return logger, err
}
