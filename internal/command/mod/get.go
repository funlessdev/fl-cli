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
	Name  string `arg:"" help:"Name of the module"`
	Count bool   `name:"count" short:"c" default:"false" help:"Return number of results"`
}

func (c *Get) Help() string {
	return `
DESCRIPTION

	List all the functions and informations about the specified module.
	The "--count" flag can be used to return the number of results

EXAMPLES
	
	$ fl mod get <your-module-name> --count
`
}

func (g *Get) Run(ctx context.Context, modHandler client.ModHandler, logger log.FLogger) error {
	res, err := modHandler.Get(ctx, g.Name)
	if err != nil {
		return err
	}

	logger.Infof("Module: %s\n", res.Name)
	logger.Info("Functions:\n")

	for _, name := range res.Functions {
		logger.Info(name + "\n")
	}

	if g.Count {
		logger.Infof("Count: %d\n", len(res.Functions))
	}

	return nil
}
