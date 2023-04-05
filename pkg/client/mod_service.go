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

func (mod *ModService) injectAPIToken() {
	if mod.Client != nil {
		apiToken := mod.Client.Config.APIToken
		apiConfig := mod.Client.ApiClient.GetConfig()
		apiConfig.DefaultHeader["Authorization"] = "Bearer " + apiToken
	}
}

func (mod *ModService) injectHost(ctx context.Context) {
	overrideHost, ok := ctx.Value(pkg.FLContextKey("api_host")).(string)
	if ok && overrideHost != "" {
		apiConfig := mod.Client.ApiClient.GetConfig()
		apiConfig.Host = overrideHost
	}
}

func (mod *ModService) Get(ctx context.Context, modName string) (pkg.SingleModule, error) {

	mod.injectHost(ctx)
	mod.injectAPIToken()

	if err := mod.InputValidatorHandler.ValidateName(modName, "mod"); err != nil {
		return pkg.SingleModule{}, err
	}

	apiService := mod.Client.ApiClient.ModulesApi
	request := apiService.ShowModuleByName(ctx, modName)
	response, _, err := request.Execute()
	if err != nil {
		return pkg.SingleModule{}, pkg.ExtractError(err)
	}
	data := response.GetData()
	name := data.Name

	var functions []string
	for _, mod := range data.Functions {
		functions = append(functions, *mod.Name)
	}

	return pkg.SingleModule{
		Name:      *name,
		Functions: functions,
	}, nil

}

func (mod *ModService) Create(ctx context.Context, modName string) error {

	mod.injectHost(ctx)
	mod.injectAPIToken()

	if err := mod.InputValidatorHandler.ValidateName(modName, "mod"); err != nil {
		return err
	}

	apiService := mod.Client.ApiClient.ModulesApi

	requestBody := openapi.ModuleName{
		Module: &openapi.SubjectNameSubject{
			Name: &modName,
		},
	}

	_, err := apiService.CreateModule(ctx).ModuleName(requestBody).Execute()
	return pkg.ExtractError(err)
}

func (mod *ModService) Delete(ctx context.Context, modName string) error {

	mod.injectHost(ctx)
	mod.injectAPIToken()

	if err := mod.InputValidatorHandler.ValidateName(modName, "mod"); err != nil {
		return err
	}

	apiService := mod.Client.ApiClient.ModulesApi
	_, err := apiService.DeleteModule(ctx, modName).Execute()
	return pkg.ExtractError(err)
}

func (mod *ModService) Update(ctx context.Context, modName string, newName string) error {

	mod.injectHost(ctx)
	mod.injectAPIToken()

	if err := mod.InputValidatorHandler.ValidateName(modName, "mod"); err != nil {
		return err
	}
	if err := mod.InputValidatorHandler.ValidateName(newName, "new mod"); err != nil {
		return err
	}

	apiService := mod.Client.ApiClient.ModulesApi
	requestBody := openapi.ModuleName{
		Module: &openapi.SubjectNameSubject{
			Name: &newName,
		},
	}
	request := apiService.UpdateModule(ctx, modName).ModuleName2(requestBody)
	_, err := request.Execute()
	return pkg.ExtractError(err)
}

func (mod *ModService) List(ctx context.Context) (pkg.ModuleNameList, error) {

	mod.injectHost(ctx)
	mod.injectAPIToken()

	apiService := mod.Client.ApiClient.ModulesApi
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
