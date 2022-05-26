// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package main

import (
	"fmt"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/funlessdev/funless-cli/pkg/client"
	"github.com/funlessdev/funless-cli/pkg/command"
)

type CLI struct {
	Fn    command.Fn    `cmd:"" help:"todo fn subcommand help"`
	Admin command.Admin `cmd:"" help:"todo admin subcommand help"`
}

func main() {
	cli := CLI{}

	flConfig := client.Config{Host: "http://localhost:8080"}
	flClient, err := client.NewClient(http.DefaultClient, flConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fnSvc := &client.FnService{Client: flClient}

	ctx := kong.Parse(&cli,
		kong.Name("fl"),
		kong.Description("Funless CLI - fl"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:             true,
			NoExpandSubcommands: true,
		}),
		kong.BindTo(fnSvc, (*client.FnHandler)(nil)),
	)

	ctx.FatalIfErrorf(ctx.Run())
}
