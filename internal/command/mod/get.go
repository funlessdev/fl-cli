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

package mod

import (
	"context"

	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Get struct {
	Name  string `arg:"" help:"name of the module"`
	Count bool   `name:"count" short:"c" default:"false" help:"return number of results"`
}

func (g *Get) Run(ctx context.Context, modHandler client.ModHandler, logger log.FLogger) error {
	res, err := modHandler.Get(ctx, g.Name)
	if err != nil {
		return extractError(err)
	}

	data := res.GetData()
	name := data.Name
	functions := data.Functions

	if err != nil {
		return extractError(err)
	}
	logger.Infof("Module: %s\n", *name)
	logger.Info("Functions:\n")

	for _, v := range functions {
		logger.Info(*v.Name + "\n")
	}

	if g.Count {
		logger.Infof("Count: %d\n", len(functions))
	}

	return nil
}
