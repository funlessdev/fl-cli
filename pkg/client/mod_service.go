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

	openapi "github.com/funlessdev/fl-client-sdk-go"
)

type ModHandler interface {
	Get(ctx context.Context, modName string) (openapi.SingleModuleResult, error)
	Create(ctx context.Context, modName string) error
	Delete(ctx context.Context, modName string) error
	Update(ctx context.Context, modName string, newName string) error
	List(ctx context.Context) (openapi.ModuleNamesResult, error)
}

type ModService struct {
	*Client
	InputValidatorHandler
}

var _ ModHandler = &ModService{}

func (fn *ModService) Get(ctx context.Context, modName string) (openapi.SingleModuleResult, error) {

	if err := fn.InputValidatorHandler.ValidateName(modName, "mod"); err != nil {
		return *openapi.NewSingleModuleResult(), err
	}

	apiService := fn.Client.ApiClient.ModulesApi
	request := apiService.ShowModuleByName(ctx, modName)
	response, _, err := request.Execute()
	if err != nil {
		return *openapi.NewSingleModuleResult(), err
	} else {
		return *response, nil
	}
}

func (fn *ModService) Create(ctx context.Context, modName string) error {

	if err := fn.InputValidatorHandler.ValidateName(modName, "mod"); err != nil {
		return err
	}

	apiService := fn.Client.ApiClient.ModulesApi

	requestBody := openapi.ModuleName{
		Module: &openapi.SubjectNameSubject{
			Name: &modName,
		},
	}

	_, err := apiService.CreateModule(ctx).ModuleName(requestBody).Execute()
	return err
}

func (fn *ModService) Delete(ctx context.Context, modName string) error {

	if err := fn.InputValidatorHandler.ValidateName(modName, "mod"); err != nil {
		return err
	}

	apiService := fn.Client.ApiClient.ModulesApi
	_, err := apiService.DeleteModule(ctx, modName).Execute()
	return err
}

func (fn *ModService) Update(ctx context.Context, modName string, newName string) error {

	if err := fn.InputValidatorHandler.ValidateName(modName, "mod"); err != nil {
		return err
	}
	if err := fn.InputValidatorHandler.ValidateName(newName, "new mod"); err != nil {
		return err
	}

	apiService := fn.Client.ApiClient.ModulesApi
	requestBody := openapi.ModuleName{
		Module: &openapi.SubjectNameSubject{
			Name: &newName,
		},
	}
	request := apiService.UpdateModule(ctx, modName).ModuleName2(requestBody)
	_, err := request.Execute()
	return err
}

func (fn *ModService) List(ctx context.Context) (openapi.ModuleNamesResult, error) {
	apiService := fn.Client.ApiClient.ModulesApi
	response, _, err := apiService.ListModules(ctx).Execute()
	if err != nil {
		return *openapi.NewModuleNamesResult(), err
	} else {
		return *response, nil
	}
}
