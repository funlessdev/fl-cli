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

package admin_user

import (
	"bytes"
	"context"
	"testing"

	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	testName := "testName"
	testToken := "some_token"
	mockResult := pkg.UserNameToken{Name: testName, Token: testToken}

	ctx := context.Background()
	var outbuf bytes.Buffer
	bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).DisableAnimation().Build()
	mockUserHandler := mocks.NewUserHandler(t)
	mockUserHandler.On("Create", ctx, "userA").Return(mockResult, nil)

	cmd := CreateUser{Name: "userA"}
	err := cmd.Run(ctx, mockUserHandler, bufLogger)
	require.NoError(t, err)
	mockUserHandler.AssertCalled(t, "Create", ctx, "userA")
	mockUserHandler.AssertNumberOfCalls(t, "Create", 1)
	mockUserHandler.AssertExpectations(t)

	require.Contains(t, outbuf.String(), testName)
	require.Contains(t, outbuf.String(), testToken)
}

func TestListUsers(t *testing.T) {
	mockResult := pkg.UserNamesList{Names: []string{"userA", "userB"}}
	ctx := context.Background()
	var outbuf bytes.Buffer
	bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).DisableAnimation().Build()
	mockUserHandler := mocks.NewUserHandler(t)
	mockUserHandler.On("List", ctx).Return(mockResult, nil)

	cmd := ListUsers{}
	err := cmd.Run(ctx, mockUserHandler, bufLogger)
	require.NoError(t, err)
	mockUserHandler.AssertCalled(t, "List", ctx)
	mockUserHandler.AssertNumberOfCalls(t, "List", 1)
	mockUserHandler.AssertExpectations(t)

	require.Contains(t, outbuf.String(), "userA")
	require.Contains(t, outbuf.String(), "userB")
}
