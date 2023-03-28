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

type Delete struct {
	Name string `arg:"" help:"name of the module to delete"`
}

func (c *Delete) Help() string {
	return `
DESCRIPTION

	It deletes the module with the specified name.

EXAMPLES

	$ fl mod list
		_
		<your-module-name>

	$ fl mod delete <your-module-name>

	$ fl mod list
		_
	
`
}

func (d *Delete) Run(ctx context.Context, modHandler client.ModHandler, logger log.FLogger) error {
	err := modHandler.Delete(ctx, d.Name)

	if err != nil {
		return err
	}
	logger.Infof("Successfully deleted module %s.\n", d.Name)
	return nil
}
