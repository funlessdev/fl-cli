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
	"os"

	openapi "github.com/funlessdev/fl-client-sdk-go"
)

type FnHandler interface {
	Invoke(ctx context.Context, fnName string, fnNamespace string, fnArgs map[string]interface{}) (openapi.FunctionInvocationSuccess, error)
	Create(ctx context.Context, fnName string, fnNamespace string, code *os.File) (openapi.FunctionCreationSuccess, error)
	Delete(ctx context.Context, fnName string, fnNamespace string) (openapi.FunctionDeletionSuccess, error)
	List(ctx context.Context, namespace string) (openapi.FunctionListSuccess, error)
}

type FnService struct {
	*Client
}

var _ FnHandler = &FnService{}

func (fn *FnService) Invoke(ctx context.Context, fnName string, fnNamespace string, fnArgs map[string]interface{}) (openapi.FunctionInvocationSuccess, error) {
	apiService := fn.Client.ApiClient.DefaultApi
	requestBody := openapi.FunctionInvocation{
		Function:  &fnName,
		Namespace: &fnNamespace,
		Args:      fnArgs,
	}
	request := apiService.V1FnInvokePost(ctx).FunctionInvocation(requestBody)
	response, _, err := apiService.V1FnInvokePostExecute(request)
	if err != nil {
		return openapi.FunctionInvocationSuccess{}, err
	}
	return *response, err
}

func (fn *FnService) Create(ctx context.Context, fnName string, fnNamespace string, code *os.File) (openapi.FunctionCreationSuccess, error) {
	apiService := fn.Client.ApiClient.DefaultApi
	request := apiService.V1FnCreatePost(ctx).Name(fnName).Namespace(fnNamespace).Code(code)
	response, _, err := apiService.V1FnCreatePostExecute(request)

	if err != nil {
		return openapi.FunctionCreationSuccess{}, err
	}
	return *response, err
}

func (fn *FnService) Delete(ctx context.Context, fnName string, fnNamespace string) (openapi.FunctionDeletionSuccess, error) {
	apiService := fn.Client.ApiClient.DefaultApi
	requestBody := openapi.FunctionDeletion{
		Name:      &fnName,
		Namespace: &fnNamespace,
	}
	request := apiService.V1FnDeleteDelete(ctx).FunctionDeletion(requestBody)
	response, _, err := apiService.V1FnDeleteDeleteExecute(request)
	if err != nil {
		return openapi.FunctionDeletionSuccess{}, err
	}
	return *response, err
}

func (fn *FnService) List(ctx context.Context, namespace string) (openapi.FunctionListSuccess, error) {
	apiService := fn.Client.ApiClient.DefaultApi
	request := apiService.V1FnListFnNamespaceGet(ctx, namespace)
	response, _, err := apiService.V1FnListFnNamespaceGetExecute(request)
	if err != nil {
		return openapi.FunctionListSuccess{}, err
	}
	return *response, err
}
