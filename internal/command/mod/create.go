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

	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Create struct {
	Name string `arg:"" help:"Name of the module to create"`
}

func (c *Create) Help() string {
	return `
DESCRIPRION

	It creates a new module with the specified name.
	Module name must respect the pattern [a-zA-Z0-9_]*, 
	i.e., it must be an alphanumeric string with underscores allowed.

EXAMPLES

	$ fl mod list
		_
		
	$ fl mod create <your-module-name>

	$ fl mod list
		_
		<your-module-name>
`

}

func (c *Create) Run(ctx context.Context, modHandler client.ModHandler, logger log.FLogger, parent *Mod) error {
	ctx = context.WithValue(ctx, pkg.FLContextKey("api_host"), parent.Host)
	err := modHandler.Create(ctx, c.Name)

	if err != nil {
		return err
	}
	logger.Infof("Successfully created module %s.\n", c.Name)
	return nil
}
