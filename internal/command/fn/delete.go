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

	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Delete struct {
	Name      string `arg:"" name:"name" help:"name of the function to delete"`
	Namespace string `name:"namespace" short:"n" default:"_" help:"namespace of the function to delete"`
}

func (f *Delete) Run(ctx context.Context, fnHandler client.FnHandler, logger log.FLogger) error {
	err := fnHandler.Delete(ctx, f.Name, f.Namespace)
	if err != nil {
		return extractError(err)
	}

	logger.Infof("\nSuccessfully deleted function %s/%s.", f.Namespace, f.Name)
	return nil
}
