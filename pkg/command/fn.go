// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/funlessdev/funless-cli/pkg/client"
)

// TODO: fix tests
type (
	Fn struct {
		Invoke Invoke `cmd:"" help:"todo fn invoke help"`
		Create Create `cmd:"" help:"todo fn create help"`
		Delete Delete `cmd:"" help:"todo fn delete help"`
	}

	Create struct {
		Name      string `arg:"" name:"name" help:"name of the function to create"`
		Namespace string `name:"namespace" short:"n" help:"namespace of the function to create"`
		Source    string `name:"source" required:"" short:"s" help:"path of the source file"`
		Language  string `name:"language" required:"" short:"l" help:"programming language of the function"`
	}

	Invoke struct {
		Name      string            `arg:"" name:"name" help:"name of the function to invoke"`
		Namespace string            `name:"namespace" short:"n" help:"namespace of the function to invoke"`
		Args      map[string]string `name:"args" short:"a" help:"arguments of the function to invoke" xor:"args"`
		JsonArgs  string            `name:"json" short:"j" help:"json encoded arguments of the function to invoke; overrides args" xor:"args"`
	}

	Delete struct {
		Name      string `arg:"" name:"name" help:"name of the function to delete"`
		Namespace string `name:"namespace" short:"n" help:"namespace of the function to delete"`
	}
)

func (f *Invoke) Run(invoker client.FnHandler) error {
	var args interface{}
	if f.Args != nil {
		args = f.Args
	} else if f.JsonArgs != "" {
		err := json.Unmarshal([]byte(f.JsonArgs), &args)
		if err != nil {
			return err
		}
	}
	res, err := invoker.Invoke(f.Name, f.Namespace, args)
	if err != nil {
		return err
	}

	if res.Result != nil {
		decodedRes, err := json.Marshal(*res.Result)
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", string(decodedRes))
	} else {
		return errors.New("Received nil result")
	}

	return nil
}

func (f *Create) Run(invoker client.FnHandler) error {
	code, err := os.ReadFile(f.Source)
	if err != nil {
		return err
	}

	res, err := invoker.Create(f.Name, f.Namespace, string(code), f.Language)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", res.Result)
	return nil
}

func (f *Delete) Run(invoker client.FnHandler) error {
	res, err := invoker.Delete(f.Name, f.Namespace)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", res.Result)
	return nil
}
