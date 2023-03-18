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

package client

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/funlessdev/fl-cli/pkg"
	openapi "github.com/funlessdev/fl-client-sdk-go"
)

type FnHandler interface {
	Invoke(ctx context.Context, fnName string, fnMod string, fnArgs map[string]interface{}) (pkg.IvkResult, error)
	Create(ctx context.Context, fnName string, fnMod string, code *os.File) error
	Delete(ctx context.Context, fnName string, fnMod string) error
	Update(ctx context.Context, fnName string, fnMod string, code *os.File, newName string) error
}

type FnService struct {
	*Client
	InputValidatorHandler
}

var _ FnHandler = &FnService{}

func (fn *FnService) Invoke(ctx context.Context, fnName string, fnMod string, fnArgs map[string]interface{}) (pkg.IvkResult, error) {

	if err := fn.InputValidatorHandler.ValidateName(fnName, "function"); err != nil {
		return pkg.IvkResult{}, err
	}
	if err := fn.InputValidatorHandler.ValidateName(fnMod, "mod"); err != nil {
		return pkg.IvkResult{}, err
	}

	apiService := fn.Client.ApiClient.FunctionsApi
	invokeInput := openapi.InvokeInput{
		Args: fnArgs,
	}
	request := apiService.InvokeFunction(ctx, fnMod, fnName).InvokeInput(invokeInput)
	response, _, err := request.Execute()

	if err != nil {
		return pkg.IvkResult{}, pkg.ExtractError(err)
	}

	data := response.GetData()
	if data == nil {
		return pkg.IvkResult{}, errors.New("received no result")
	}

	decodedRes, err := json.Marshal(data)
	if err != nil {
		return pkg.IvkResult{}, err
	}

	return pkg.IvkResult{Result: string(decodedRes)}, nil
}

func (fn *FnService) Create(ctx context.Context, fnName string, fnMod string, code *os.File) error {

	if err := fn.InputValidatorHandler.ValidateName(fnName, "function"); err != nil {
		return err
	}
	if err := fn.InputValidatorHandler.ValidateName(fnMod, "mod"); err != nil {
		return err
	}

	apiService := fn.Client.ApiClient.FunctionsApi
	request := apiService.CreateFunction(ctx, fnMod).Name(fnName).Code(code)
	_, err := request.Execute()
	return pkg.ExtractError(err)
}

func (fn *FnService) Delete(ctx context.Context, fnName string, fnMod string) error {

	if err := fn.InputValidatorHandler.ValidateName(fnName, "function"); err != nil {
		return err
	}
	if err := fn.InputValidatorHandler.ValidateName(fnMod, "mod"); err != nil {
		return err
	}

	apiService := fn.Client.ApiClient.FunctionsApi
	request := apiService.DeleteFunction(ctx, fnMod, fnName)
	_, err := request.Execute()
	return pkg.ExtractError(err)
}

func (fn *FnService) Update(ctx context.Context, fnName string, fnMod string, code *os.File, newName string) error {

	if err := fn.InputValidatorHandler.ValidateName(fnName, "function"); err != nil {
		return err
	}
	if err := fn.InputValidatorHandler.ValidateName(newName, "new function"); err != nil {
		return err
	}

	if err := fn.InputValidatorHandler.ValidateName(fnMod, "mod"); err != nil {
		return err
	}

	apiService := fn.Client.ApiClient.FunctionsApi
	request := apiService.UpdateFunction(ctx, fnMod, fnName).Code(code).Name(newName)
	_, err := request.Execute()
	return pkg.ExtractError(err)
}
