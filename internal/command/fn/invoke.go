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
	"encoding/json"

	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Invoke struct {
	Name     string            `arg:"" name:"name" help:"Name of the function to invoke"`
	Module   string            `name:"module" short:"n" default:"_" help:"Module of the function to invoke"`
	Args     map[string]string `name:"args" short:"a" help:"Arguments of the function to invoke" xor:"args"`
	JsonArgs string            `name:"json" short:"j" help:"Json encoded arguments of the function to invoke; overrides args" xor:"args"`
}

func (c *Invoke) Help() string {
	return `
DESCRIPTION

	It invokes the function with the specified name.
	The "--module" flag can be used to choose a module other than the default one.
	The "--args" and "--json" flags can be used to pass parameters to functions. 
	"--json" flag overrides "--args"
	
EXAMPLES
	
	$ fl fn invoke <your-function-name> --module=<your-module-name>
`
}

func (f *Invoke) Run(ctx context.Context, fnHandler client.FnHandler, logger log.FLogger) error {
	args := make(map[string]interface{}, len(f.Args))
	if f.Args != nil {
		for k, v := range f.Args {
			args[k] = v
		}
	} else if f.JsonArgs != "" {
		err := json.Unmarshal([]byte(f.JsonArgs), &args)
		if err != nil {
			return err
		}
	}
	res, err := fnHandler.Invoke(ctx, f.Name, f.Module, args)
	if err != nil {
		return err
	}

	logger.Info(res.Result)

	return nil
}
