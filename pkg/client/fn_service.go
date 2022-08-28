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
package client

import (
	"context"

	swagger "github.com/funlessdev/fl-client-sdk-go"
)

//					 default
// https://${HOST}/{NAMESPACE}/fn/hello
// https://${HOST}/{NAMESPACE}/ev/event1
// https://${HOST}/{NAMESPACE}/pkg/{PACKAGE}/fn/hello

type FnHandler interface {
	Invoke(ctx context.Context, fnName string, fnNamespace string, fnArgs interface{}) (swagger.FunctionInvocationSuccess, error)
	Create(ctx context.Context, fnName string, fnNamespace string, code string, language string) (swagger.FunctionCreationSuccess, error)
	Delete(ctx context.Context, fnName string, fnNamespace string) (swagger.FunctionDeletionSuccess, error)
}

type FnService struct {
	*Client
}

var _ FnHandler = &FnService{}

func (fn *FnService) Invoke(ctx context.Context, fnName string, fnNamespace string, fnArgs interface{}) (swagger.FunctionInvocationSuccess, error) {
	apiService := fn.Client.ApiClient.DefaultApi
	response, _, err := apiService.InvokePost(ctx, swagger.FunctionInvocation{
		Function:  fnName,
		Namespace: fnNamespace,
		Args:      &fnArgs,
	})
	if err != nil {
		return swagger.FunctionInvocationSuccess{}, err
	}
	return response, err
}

func (fn *FnService) Create(ctx context.Context, fnName string, fnNamespace string, code string, language string) (swagger.FunctionCreationSuccess, error) {
	apiService := fn.Client.ApiClient.DefaultApi
	response, _, err := apiService.CreatePost(ctx, swagger.FunctionCreation{
		Name:      fnName,
		Namespace: fnNamespace,
		Code:      code,
		Image:     language,
	})
	if err != nil {
		return swagger.FunctionCreationSuccess{}, err
	}
	return response, err
}

func (fn *FnService) Delete(ctx context.Context, fnName string, fnNamespace string) (swagger.FunctionDeletionSuccess, error) {
	apiService := fn.Client.ApiClient.DefaultApi
	response, _, err := apiService.DeletePost(ctx, swagger.FunctionDeletion{
		Name:      fnName,
		Namespace: fnNamespace,
	})
	if err != nil {
		return swagger.FunctionDeletionSuccess{}, err
	}
	return response, err
}
