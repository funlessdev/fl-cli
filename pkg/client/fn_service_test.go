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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/funlessdev/fl-cli/test/mocks"
	openapi "github.com/funlessdev/fl-client-sdk-go"
)

func TestFnInvoke(t *testing.T) {
	testFn := "test_fn"
	testMod := "test_mod"
	var testArgs map[string]interface{} = map[string]interface{}{"name": "Some name"}

	testCtx := context.Background()
	mockValidator := mocks.NewInputValidatorHandler(t)
	mockValidator.On("ValidateName", testFn, "function").Return(nil)
	mockValidator.On("ValidateName", testMod, "mod").Return(nil)

	t.Run("should send invoke request to server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, fmt.Sprintf("/v1/fn/%s/%s", testMod, testFn), r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]map[string]string{"Data": {"payload": "some result"}}
			jresult, _ := json.Marshal(result)
			_, _ = w.Write(jresult)
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &FnService{Client: c, InputValidatorHandler: mockValidator}

		result, err := svc.Invoke(testCtx, testFn, testMod, testArgs)

		require.NoError(t, err)
		expected := map[string]interface{}{"payload": "some result"}
		assert.Equal(t, expected, result.GetData())

		mockValidator.AssertNumberOfCalls(t, "ValidateName", 2)
		mockValidator.AssertExpectations(t)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, fmt.Sprintf("/v1/fn/%s/%s", testMod, testFn), r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"error": "some error"}
			jresult, _ := json.Marshal(result)

			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, string(jresult))
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &FnService{Client: c, InputValidatorHandler: mockValidator}

		_, err := svc.Invoke(testCtx, testFn, testMod, testArgs)

		require.Error(t, err)
		openApiError := err.(*openapi.GenericOpenAPIError)
		assert.Equal(t, "{\"error\":\"some error\"}\n", string(openApiError.Body()))
	})
}

func TestFnCreate(t *testing.T) {
	testFn := "test_fn"
	testMod := "test_mod"
	testSource, _ := filepath.Abs("../../test/fixtures/real.wasm")
	testCode, _ := os.Open(testSource)

	testCtx := context.Background()

	mockValidator := mocks.NewInputValidatorHandler(t)
	mockValidator.On("ValidateName", testFn, "function").Return(nil)
	mockValidator.On("ValidateName", testMod, "mod").Return(nil)

	t.Run("should send create request to server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, fmt.Sprintf("/v1/fn/%s", testMod), r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"Result": "some result"}
			jresult, _ := json.Marshal(result)
			_, _ = w.Write(jresult)
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &FnService{Client: c, InputValidatorHandler: mockValidator}

		err := svc.Create(testCtx, testFn, testMod, testCode)

		require.NoError(t, err)

		mockValidator.AssertNumberOfCalls(t, "ValidateName", 2)
		mockValidator.AssertExpectations(t)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, fmt.Sprintf("/v1/fn/%s", testMod), r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"error": "some error"}
			jresult, _ := json.Marshal(result)

			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, string(jresult))
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &FnService{Client: c, InputValidatorHandler: mockValidator}

		err := svc.Create(testCtx, testFn, testMod, testCode)

		require.Error(t, err)
		openApiError := err.(*openapi.GenericOpenAPIError)
		assert.Equal(t, "{\"error\":\"some error\"}\n", string(openApiError.Body()))
	})
}

func TestFnDelete(t *testing.T) {
	testFn := "test_fn"
	testMod := "test_mod"

	testCtx := context.Background()

	mockValidator := mocks.NewInputValidatorHandler(t)
	mockValidator.On("ValidateName", testFn, "function").Return(nil)
	mockValidator.On("ValidateName", testMod, "mod").Return(nil)

	t.Run("should send delete request to server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, fmt.Sprintf("/v1/fn/%s/%s", testMod, testFn), r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"Result": "some result"}
			jresult, _ := json.Marshal(result)
			_, _ = w.Write(jresult)
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &FnService{Client: c, InputValidatorHandler: mockValidator}

		err := svc.Delete(testCtx, testFn, testMod)

		require.NoError(t, err)

		mockValidator.AssertNumberOfCalls(t, "ValidateName", 2)
		mockValidator.AssertExpectations(t)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, fmt.Sprintf("/v1/fn/%s/%s", testMod, testFn), r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"error": "some error"}
			jresult, _ := json.Marshal(result)

			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, string(jresult))
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &FnService{Client: c, InputValidatorHandler: mockValidator}

		err := svc.Delete(testCtx, testFn, testMod)

		require.Error(t, err)
		openApiError := err.(*openapi.GenericOpenAPIError)
		assert.Equal(t, "{\"error\":\"some error\"}\n", string(openApiError.Body()))
	})
}
