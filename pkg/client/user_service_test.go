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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	testCtx := context.Background()
	testName := "testName"

	t.Run("should return name,token", func(t *testing.T) {
		res := map[string]map[string]string{"data": {"name": "some_name", "token": "some_token"}}
		client := setupHttpServer(t, "/v1/admin/subjects", http.MethodPost, res, http.StatusOK)
		svc := &UserService{Client: client}
		result, err := svc.Create(testCtx, testName)
		require.NoError(t, err)
		require.Equal(t, "some_name", result.Name)
		require.Equal(t, "some_token", result.Token)
	})
	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		res := map[string]map[string]string{"errors": {"detail": "some error"}}
		client := setupHttpServer(t, "/v1/admin/subjects", http.MethodPost, res, http.StatusNotFound)
		svc := &UserService{Client: client}
		_, err := svc.Create(testCtx, testName)
		require.Error(t, err)
		assert.Equal(t, "some error", err.Error())
	})
}

func TestList(t *testing.T) {
	testCtx := context.Background()

	t.Run("should return list of names", func(t *testing.T) {
		res := map[string][]map[string]string{"data": {{"name": "name"}, {"name": "name2"}}}
		client := setupHttpServer(t, "/v1/admin/subjects", http.MethodGet, res, http.StatusOK)
		svc := &UserService{Client: client}
		result, err := svc.List(testCtx)
		require.NoError(t, err)

		require.Equal(t, "name", result.Names[0])
		require.Equal(t, "name2", result.Names[1])
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		res := map[string]map[string]string{"errors": {"detail": "some error"}}
		client := setupHttpServer(t, "/v1/admin/subjects", http.MethodGet, res, http.StatusNotFound)
		svc := &UserService{Client: client}
		_, err := svc.List(testCtx)
		require.Error(t, err)
		assert.Equal(t, "some error", err.Error())
	})
}
