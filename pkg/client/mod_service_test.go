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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/funlessdev/fl-cli/test/mocks"
	openapi "github.com/funlessdev/fl-client-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModGet(t *testing.T) {
	testMod := "test_mod"

	testCtx := context.Background()

	mockValidator := mocks.NewInputValidatorHandler(t)
	mockValidator.On("ValidateName", testMod, "mod").Return(nil)

	t.Run("should send get request to server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, fmt.Sprintf("/v1/fn/%s", testMod), r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]map[string]interface{}{"Data": {"Name": testMod, "Functions": []string{}}}
			jresult, _ := json.Marshal(result)
			_, _ = w.Write(jresult)
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &ModService{Client: c, InputValidatorHandler: mockValidator}

		result, err := svc.Get(testCtx, testMod)

		require.NoError(t, err)
		expected := *openapi.NewSingleModuleResult()
		expected.Data = openapi.NewSingleModuleResultData()
		expected.Data.Name = &testMod
		expected.Data.Functions = []openapi.ModuleNameModule{}

		assert.Equal(t, expected, result)

		mockValidator.AssertNumberOfCalls(t, "ValidateName", 1)
		mockValidator.AssertExpectations(t)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, fmt.Sprintf("/v1/fn/%s", testMod), r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"error": "some error"}
			jresult, _ := json.Marshal(result)

			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, string(jresult))
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &ModService{Client: c, InputValidatorHandler: mockValidator}

		_, err := svc.Get(testCtx, testMod)

		require.Error(t, err)
		openApiError := err.(*openapi.GenericOpenAPIError)
		assert.Equal(t, "{\"error\":\"some error\"}\n", string(openApiError.Body()))
	})
}

func TestModCreate(t *testing.T) {
	testMod := "test_mod"

	testCtx := context.Background()

	mockValidator := mocks.NewInputValidatorHandler(t)
	mockValidator.On("ValidateName", testMod, "mod").Return(nil)

	t.Run("should send create request to server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/v1/fn", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			jresult, _ := json.Marshal(nil)
			_, _ = w.Write(jresult)
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &ModService{Client: c, InputValidatorHandler: mockValidator}

		err := svc.Create(testCtx, testMod)

		require.NoError(t, err)

		mockValidator.AssertNumberOfCalls(t, "ValidateName", 1)
		mockValidator.AssertExpectations(t)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/v1/fn", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"error": "some error"}
			jresult, _ := json.Marshal(result)

			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, string(jresult))
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &ModService{Client: c, InputValidatorHandler: mockValidator}

		err := svc.Create(testCtx, testMod)

		require.Error(t, err)
		openApiError := err.(*openapi.GenericOpenAPIError)
		assert.Equal(t, "{\"error\":\"some error\"}\n", string(openApiError.Body()))
	})
}

func TestModDelete(t *testing.T) {
	testMod := "test_mod"

	testCtx := context.Background()

	mockValidator := mocks.NewInputValidatorHandler(t)
	mockValidator.On("ValidateName", testMod, "mod").Return(nil)

	t.Run("should send delete request to server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, fmt.Sprintf("/v1/fn/%s", testMod), r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			jresult, _ := json.Marshal(nil)
			_, _ = w.Write(jresult)
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &ModService{Client: c, InputValidatorHandler: mockValidator}

		err := svc.Delete(testCtx, testMod)

		require.NoError(t, err)

		mockValidator.AssertNumberOfCalls(t, "ValidateName", 1)
		mockValidator.AssertExpectations(t)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, fmt.Sprintf("/v1/fn/%s", testMod), r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"error": "some error"}
			jresult, _ := json.Marshal(result)

			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, string(jresult))
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &ModService{Client: c, InputValidatorHandler: mockValidator}

		err := svc.Delete(testCtx, testMod)

		require.Error(t, err)
		openApiError := err.(*openapi.GenericOpenAPIError)
		assert.Equal(t, "{\"error\":\"some error\"}\n", string(openApiError.Body()))
	})
}

func TestModUpdate(t *testing.T) {
	testMod := "test_mod"
	testNewMod := "test_mod_2"

	testCtx := context.Background()

	mockValidator := mocks.NewInputValidatorHandler(t)
	mockValidator.On("ValidateName", testMod, "mod").Return(nil)
	mockValidator.On("ValidateName", testNewMod, "new mod").Return(nil)

	t.Run("should send update request to server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPut, r.Method)
			assert.Equal(t, fmt.Sprintf("/v1/fn/%s", testMod), r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			jresult, _ := json.Marshal(nil)
			_, _ = w.Write(jresult)
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &ModService{Client: c, InputValidatorHandler: mockValidator}

		err := svc.Update(testCtx, testMod, testNewMod)

		require.NoError(t, err)

		mockValidator.AssertNumberOfCalls(t, "ValidateName", 2)
		mockValidator.AssertExpectations(t)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPut, r.Method)
			assert.Equal(t, fmt.Sprintf("/v1/fn/%s", testMod), r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"error": "some error"}
			jresult, _ := json.Marshal(result)

			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, string(jresult))
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &ModService{Client: c, InputValidatorHandler: mockValidator}

		err := svc.Update(testCtx, testMod, testNewMod)

		require.Error(t, err)
		openApiError := err.(*openapi.GenericOpenAPIError)
		assert.Equal(t, "{\"error\":\"some error\"}\n", string(openApiError.Body()))
	})
}

func TestModList(t *testing.T) {
	testCtx := context.Background()
	f1 := "f1"
	f2 := "f2"
	f3 := "f3"

	t.Run("should send list request to server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/v1/fn", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string][]map[string]string{"Data": {{"name": f1}, {"name": f2}, {"name": f3}}}
			jresult, _ := json.Marshal(result)
			_, _ = w.Write(jresult)
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &ModService{Client: c}

		result, err := svc.List(testCtx)

		require.NoError(t, err)
		expected := *openapi.NewModuleNamesResult()
		expected.Data = []openapi.ModuleNameModule{{Name: &f1}, {Name: &f2}, {Name: &f3}}
		assert.Equal(t, expected, result)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/v1/fn", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"error": "some error"}
			jresult, _ := json.Marshal(result)

			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, string(jresult))
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &ModService{Client: c}

		_, err := svc.List(testCtx)

		require.Error(t, err)
		openApiError := err.(*openapi.GenericOpenAPIError)
		assert.Equal(t, "{\"error\":\"some error\"}\n", string(openApiError.Body()))
	})

}
