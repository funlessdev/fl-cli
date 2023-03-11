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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/funlessdev/fl-cli/test/mocks"
	openapi "github.com/funlessdev/fl-client-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testMod string = "test_mod"

func setupHttpServer(t *testing.T, expectedPath string, expectedhttpMethod string, result interface{}, status int) *Client {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedhttpMethod, r.Method)
		assert.Equal(t, expectedPath, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		jresult, _ := json.Marshal(result)
		w.WriteHeader(status)
		_, _ = w.Write(jresult)
	}))
	t.Cleanup(func() {
		server.Close()
	})
	c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
	return c
}

func TestModGet(t *testing.T) {
	testCtx := context.Background()

	mockValidator := mocks.NewInputValidatorHandler(t)
	mockValidator.On("ValidateName", testMod, "mod").Return(nil)

	t.Run("should send get request to server", func(t *testing.T) {
		res := map[string]map[string]interface{}{"Data": {"Name": testMod, "Functions": []string{}}}
		client := setupHttpServer(t, fmt.Sprintf("/v1/fn/%s", testMod), http.MethodGet, res, http.StatusOK)

		svc := &ModService{Client: client, InputValidatorHandler: mockValidator}

		result, err := svc.Get(testCtx, testMod)

		require.NoError(t, err)

		tmod := testMod
		expected := *openapi.NewSingleModuleResult()
		expected.Data = openapi.NewSingleModuleResultData()
		expected.Data.Name = &tmod
		expected.Data.Functions = []openapi.ModuleNameModule{}

		assert.Equal(t, expected, result)

		mockValidator.AssertNumberOfCalls(t, "ValidateName", 1)
		mockValidator.AssertExpectations(t)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		res := map[string]string{"error": "some error"}
		client := setupHttpServer(t, fmt.Sprintf("/v1/fn/%s", testMod), http.MethodGet, res, http.StatusNotFound)
		svc := &ModService{Client: client, InputValidatorHandler: mockValidator}
		_, err := svc.Get(testCtx, testMod)
		require.Error(t, err)
		openApiError := err.(*openapi.GenericOpenAPIError)
		assert.Equal(t, "{\"error\":\"some error\"}", string(openApiError.Body()))
	})
}

func TestModCreate(t *testing.T) {
	testCtx := context.Background()
	mockValidator := mocks.NewInputValidatorHandler(t)
	mockValidator.On("ValidateName", testMod, "mod").Return(nil)

	t.Run("should send create request to server", func(t *testing.T) {
		client := setupHttpServer(t, "/v1/fn", http.MethodPost, nil, http.StatusOK)
		svc := &ModService{Client: client, InputValidatorHandler: mockValidator}
		err := svc.Create(testCtx, testMod)
		require.NoError(t, err)
		mockValidator.AssertNumberOfCalls(t, "ValidateName", 1)
		mockValidator.AssertExpectations(t)
	})

	t.Run("should return error if input is invalid", func(t *testing.T) {
		errMockValidator := mocks.NewInputValidatorHandler(t)
		errMockValidator.On("ValidateName", testMod, "mod").Return(errors.New("invalid error"))
		svc := &ModService{Client: nil, InputValidatorHandler: errMockValidator}
		err := svc.Create(testCtx, testMod)
		require.Error(t, err)
		require.Equal(t, "invalid error", err.Error())
		errMockValidator.AssertNumberOfCalls(t, "ValidateName", 1)
		errMockValidator.AssertExpectations(t)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		res := map[string]string{"error": "some error"}
		client := setupHttpServer(t, "/v1/fn", http.MethodPost, res, http.StatusNotFound)
		svc := &ModService{Client: client, InputValidatorHandler: mockValidator}
		err := svc.Create(testCtx, testMod)
		require.Error(t, err)
		openApiError := err.(*openapi.GenericOpenAPIError)
		assert.Equal(t, "{\"error\":\"some error\"}", string(openApiError.Body()))
	})
}

func TestModDelete(t *testing.T) {
	testCtx := context.Background()

	mockValidator := mocks.NewInputValidatorHandler(t)
	mockValidator.On("ValidateName", testMod, "mod").Return(nil)

	t.Run("should send delete request to server", func(t *testing.T) {
		client := setupHttpServer(t, fmt.Sprintf("/v1/fn/%s", testMod), http.MethodDelete, nil, http.StatusOK)
		svc := &ModService{Client: client, InputValidatorHandler: mockValidator}

		err := svc.Delete(testCtx, testMod)
		require.NoError(t, err)

		mockValidator.AssertNumberOfCalls(t, "ValidateName", 1)
		mockValidator.AssertExpectations(t)
	})

	t.Run("should return error if input is invalid", func(t *testing.T) {
		errMockValidator := mocks.NewInputValidatorHandler(t)
		errMockValidator.On("ValidateName", testMod, "mod").Return(errors.New("invalid error"))
		svc := &ModService{Client: nil, InputValidatorHandler: errMockValidator}
		err := svc.Delete(testCtx, testMod)
		require.Error(t, err)
		require.Equal(t, "invalid error", err.Error())
		errMockValidator.AssertNumberOfCalls(t, "ValidateName", 1)
		errMockValidator.AssertExpectations(t)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		res := map[string]string{"error": "some error"}
		client := setupHttpServer(t, fmt.Sprintf("/v1/fn/%s", testMod), http.MethodDelete, res, http.StatusNotFound)
		svc := &ModService{Client: client, InputValidatorHandler: mockValidator}
		err := svc.Delete(testCtx, testMod)

		require.Error(t, err)
		openApiError := err.(*openapi.GenericOpenAPIError)
		assert.Equal(t, "{\"error\":\"some error\"}", string(openApiError.Body()))
	})
}

func TestModUpdate(t *testing.T) {
	testNewMod := "test_mod_2"

	testCtx := context.Background()

	mockValidator := mocks.NewInputValidatorHandler(t)
	mockValidator.On("ValidateName", testMod, "mod").Return(nil)
	mockValidator.On("ValidateName", testNewMod, "new mod").Return(nil)

	t.Run("should send update request to server", func(t *testing.T) {
		client := setupHttpServer(t, fmt.Sprintf("/v1/fn/%s", testMod), http.MethodPut, nil, http.StatusOK)
		svc := &ModService{Client: client, InputValidatorHandler: mockValidator}
		err := svc.Update(testCtx, testMod, testNewMod)
		require.NoError(t, err)
		mockValidator.AssertNumberOfCalls(t, "ValidateName", 2)
		mockValidator.AssertExpectations(t)
	})

	t.Run("should return error if input is invalid", func(t *testing.T) {
		errMockValidator := mocks.NewInputValidatorHandler(t)
		errMockValidator.On("ValidateName", testMod, "mod").Return(errors.New("invalid error"))
		svc := &ModService{Client: nil, InputValidatorHandler: errMockValidator}
		err := svc.Update(testCtx, testMod, testNewMod)
		require.Error(t, err)
		require.Equal(t, "invalid error", err.Error())
		errMockValidator.AssertNumberOfCalls(t, "ValidateName", 1)
		errMockValidator.AssertExpectations(t)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		res := map[string]string{"error": "some error"}
		client := setupHttpServer(t, fmt.Sprintf("/v1/fn/%s", testMod), http.MethodPut, res, http.StatusNotFound)
		svc := &ModService{Client: client, InputValidatorHandler: mockValidator}
		err := svc.Update(testCtx, testMod, testNewMod)
		require.Error(t, err)
		openApiError := err.(*openapi.GenericOpenAPIError)
		assert.Equal(t, "{\"error\":\"some error\"}", string(openApiError.Body()))
	})
}

func TestModList(t *testing.T) {
	testCtx := context.Background()
	f1 := "f1"
	f2 := "f2"
	f3 := "f3"

	t.Run("should send list request to server", func(t *testing.T) {
		res := map[string][]map[string]string{"Data": {{"name": f1}, {"name": f2}, {"name": f3}}}
		client := setupHttpServer(t, "/v1/fn", http.MethodGet, res, http.StatusOK)
		svc := &ModService{Client: client}
		result, err := svc.List(testCtx)
		require.NoError(t, err)
		expected := *openapi.NewModuleNamesResult()
		expected.Data = []openapi.ModuleNameModule{{Name: &f1}, {Name: &f2}, {Name: &f3}}
		assert.Equal(t, expected, result)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		res := map[string]string{"error": "some error"}
		client := setupHttpServer(t, "/v1/fn", http.MethodGet, res, http.StatusNotFound)
		svc := &ModService{Client: client}
		_, err := svc.List(testCtx)
		require.Error(t, err)
		openApiError := err.(*openapi.GenericOpenAPIError)
		assert.Equal(t, "{\"error\":\"some error\"}", string(openApiError.Body()))
	})

}
