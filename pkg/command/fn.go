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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	swagger "github.com/funlessdev/fl-client-sdk-go"
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
		Name      string        `arg:"" name:"name" help:"name of the function to create"`
		Namespace string        `name:"namespace" short:"n" help:"namespace of the function to create"`
		Source    string        `name:"source" required:"" short:"s" help:"path of the source file"`
		Language  string        `name:"language" required:"" short:"l" help:"programming language of the function"`
		FS        fs.ReadFileFS `kong:"-"`
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

func (f *Invoke) Run(ctx context.Context, invoker client.FnHandler, writer io.Writer) error {
	var args interface{}
	if f.Args != nil {
		args = f.Args
	} else if f.JsonArgs != "" {
		err := json.Unmarshal([]byte(f.JsonArgs), &args)
		if err != nil {
			return err
		}
	}
	res, err := invoker.Invoke(ctx, f.Name, f.Namespace, args)
	if err != nil {
		return extractError(err)
	}

	if res.Result != nil {
		decodedRes, err := json.Marshal(*res.Result)
		if err != nil {
			return err
		}
		fmt.Fprintln(writer, string(decodedRes))
	} else {
		return errors.New("Received nil result")
	}

	return nil
}

func (f *Create) Run(ctx context.Context, invoker client.FnHandler, writer io.Writer) error {
	var code []byte
	var err error

	if f.FS != nil {
		code, err = fs.ReadFile(f.FS, f.Source)
	} else {
		code, err = os.ReadFile(f.Source)
	}

	if err != nil {
		return err
	}

	res, err := invoker.Create(ctx, f.Name, f.Namespace, string(code), f.Language)
	if err != nil {
		return extractError(err)
	}

	fmt.Fprintln(writer, res.Result)
	return nil
}

func (f *Delete) Run(ctx context.Context, invoker client.FnHandler, writer io.Writer) error {
	res, err := invoker.Delete(ctx, f.Name, f.Namespace)
	if err != nil {
		return extractError(err)
	}

	fmt.Fprintln(writer, res.Result)
	return nil
}

func extractError(err error) error {
	swaggerError, ok_sw := err.(swagger.GenericSwaggerError)
	if ok_sw {
		switch swaggerError.Model().(type) {
		case swagger.FunctionCreationError:
			specificError := swaggerError.Model().(swagger.FunctionCreationError)
			return errors.New(specificError.Error_)
		case swagger.FunctionDeletionError:
			specificError := swaggerError.Model().(swagger.FunctionDeletionError)
			return errors.New(specificError.Error_)
		case swagger.FunctionInvocationError:
			specificError := swaggerError.Model().(swagger.FunctionInvocationError)
			return errors.New(specificError.Error_)
		}
	}
	return err
}
