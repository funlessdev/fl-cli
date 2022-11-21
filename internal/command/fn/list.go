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
	"errors"

	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type List struct {
	Namespace string `arg:"" name:"namespace" default:"_" help:"namespace of the functions to list"`
	Count     bool   `name:"count" short:"c" default:"false" help:"return number of results"`
}

func (f *List) Run(ctx context.Context, fnHandler client.FnHandler, logger log.FLogger) error {
	res, err := fnHandler.List(ctx, f.Namespace)
	if err != nil {
		return extractError(err)
	}

	if res.Result != nil {
		// TODO: extract list items and print
		if err != nil {
			return err
		}
		for _, v := range res.Result {
			logger.Info(v)
		}
		if f.Count {
			logger.Infof("Count: %d\n", len(res.Result))
		}
	} else {
		return errors.New("received nil result")
	}

	return nil
}
