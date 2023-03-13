// Copyright 2023 Giuseppe De Palma, Matteo Trentin
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

type UserHandler interface {
	Create(ctx context.Context, name string) (pkg.UserNameToken, error)
	List(ctx context.Context) error
}

type UserService struct {
	*Client
}

var _ UserHandler = &UserService{}

func (u *UserService) Create(ctx context.Context, name string) (pkg.UserNameToken, error) {
	apiService := u.Client.ApiClient.SubjectsApi

	requestBody := openapi.SubjectName{
		Subject: &openapi.SubjectNameSubject{
			Name: &name,
		},
	}
	res, _, err := apiService.CreateSubject(ctx).SubjectName(requestBody).Execute()

	if err != nil {
		return pkg.UserNameToken{}, pkg.ExtractError(err)
	}

	data := res.GetData()
	return pkg.UserNameToken{Name: *data.Name, Token: *data.Token}, err
}

func (u *UserService) List(ctx context.Context) error {
	return nil
}
