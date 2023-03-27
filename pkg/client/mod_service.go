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

	"github.com/funlessdev/fl-cli/pkg"
	openapi "github.com/funlessdev/fl-client-sdk-go"
)

type ModHandler interface {
	Get(ctx context.Context, modName string) (pkg.SingleModule, error)
	Create(ctx context.Context, modName string) error
	Delete(ctx context.Context, modName string) error
	Update(ctx context.Context, modName string, newName string) error
	List(ctx context.Context) (pkg.ModuleNameList, error)
}

type ModService struct {
	*Client
	InputValidatorHandler
}

var _ ModHandler = &ModService{}

func (fn *ModService) setAPIToken() {
	if fn.Client != nil {
		apiToken := fn.Client.Config.APIToken
		apiConfig := fn.Client.ApiClient.GetConfig()
		apiConfig.DefaultHeader["Authorization"] = "Bearer " + apiToken
	}
}

func (fn *ModService) Get(ctx context.Context, modName string) (pkg.SingleModule, error) {

	fn.setAPIToken()

	if err := fn.InputValidatorHandler.ValidateName(modName, "mod"); err != nil {
		return pkg.SingleModule{}, err
	}

	apiService := fn.Client.ApiClient.ModulesApi
	request := apiService.ShowModuleByName(ctx, modName)
	response, _, err := request.Execute()
	if err != nil {
		return pkg.SingleModule{}, pkg.ExtractError(err)
	}
	data := response.GetData()
	name := data.Name

	var functions []string
	for _, fn := range data.Functions {
		functions = append(functions, *fn.Name)
	}

	return pkg.SingleModule{
		Name:      *name,
		Functions: functions,
	}, nil

}

func (fn *ModService) Create(ctx context.Context, modName string) error {

	fn.setAPIToken()

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
	return pkg.ExtractError(err)
}

func (fn *ModService) Delete(ctx context.Context, modName string) error {

	fn.setAPIToken()

	if err := fn.InputValidatorHandler.ValidateName(modName, "mod"); err != nil {
		return err
	}

	apiService := fn.Client.ApiClient.ModulesApi
	_, err := apiService.DeleteModule(ctx, modName).Execute()
	return pkg.ExtractError(err)
}

func (fn *ModService) Update(ctx context.Context, modName string, newName string) error {

	fn.setAPIToken()

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
	return pkg.ExtractError(err)
}

func (fn *ModService) List(ctx context.Context) (pkg.ModuleNameList, error) {

	fn.setAPIToken()

	apiService := fn.Client.ApiClient.ModulesApi
	response, _, err := apiService.ListModules(ctx).Execute()
	if err != nil {
		return pkg.ModuleNameList{}, pkg.ExtractError(err)
	}

	var modules []string
	for _, mod := range response.GetData() {
		modules = append(modules, *mod.Name)
	}

	return pkg.ModuleNameList{Names: modules}, nil
}
