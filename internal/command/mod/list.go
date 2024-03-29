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

type List struct {
	Count bool `name:"count" short:"c" default:"false" help:"Return number of results"`
}

func (c *List) Help() string {
	return `
DESCRIPTION

	List all modules.
	The "--count" flag can be used to return the number of results.

EXAMPLES
	
	$ fl mod list --count
`
}

func (l *List) Run(ctx context.Context, modHandler client.ModHandler, logger log.FLogger, parent *Mod) error {
	ctx = context.WithValue(ctx, pkg.FLContextKey("api_host"), parent.Host)
	res, err := modHandler.List(ctx)

	if err != nil {
		return err
	}
	for _, name := range res.Names {
		logger.Info(name + "\n")
	}
	if l.Count {
		logger.Infof("Count: %d\n", len(res.Names))
	}
	return nil
}
