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

	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Delete struct {
	Name   string `arg:"" name:"name" help:"name of the function to delete"`
	Module string `name:"module" short:"m" default:"_" help:"module of the function to delete"`
}

func (f *Delete) Run(ctx context.Context, fnHandler client.FnHandler, logger log.FLogger, parent *Fn) error {

	ctx = context.WithValue(ctx, pkg.FLContextKey("api_host"), parent.Host)

	err := fnHandler.Delete(ctx, f.Name, f.Module)
	if err != nil {
		return err
	}
	logger.Infof("\nSuccessfully deleted function %s/%s.\n", f.Module, f.Name)
	return nil
}
