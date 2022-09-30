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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoke(t *testing.T) {
	testFn := "test-fn"
	testNs := "test-ns"
	var testArgs interface{} = map[string]string{"name": "Some name"}

	testCtx := context.Background()

	t.Run("should send invoke request to server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/invoke", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"Result": "some result"}
			jresult, _ := json.Marshal(result)
			_, _ = w.Write(jresult)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &FnService{Client: c}

		_, err := svc.Invoke(testCtx, testFn, testNs, testArgs)

		require.NoError(t, err)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/invoke", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"error": "some error"}
			jresult, _ := json.Marshal(result)

			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, string(jresult))
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &FnService{Client: c}

		_, err := svc.Invoke(testCtx, testFn, testNs, testArgs)

		require.Error(t, err)
	})
}

func TestCreate(t *testing.T) {
	testFn := "test-fn"
	testNs := "test-ns"
	testLanguage := "nodejs"
	testCode := "console.log('Something')"

	testCtx := context.Background()
	t.Run("should send create request to server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/create", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"Result": "some result"}
			jresult, _ := json.Marshal(result)
			_, _ = w.Write(jresult)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &FnService{Client: c}

		_, err := svc.Create(testCtx, testFn, testNs, testCode, testLanguage)

		require.NoError(t, err)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/create", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"error": "some error"}
			jresult, _ := json.Marshal(result)

			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, string(jresult))
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &FnService{Client: c}

		_, err := svc.Create(testCtx, testFn, testNs, testCode, testLanguage)

		require.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	testFn := "test-fn"
	testNs := "test-ns"

	testCtx := context.Background()

	t.Run("should send delete request to server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/delete", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"Result": "some result"}
			jresult, _ := json.Marshal(result)
			_, _ = w.Write(jresult)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &FnService{Client: c}

		_, err := svc.Delete(testCtx, testFn, testNs)

		require.NoError(t, err)
	})

	t.Run("should return error if request encounters an HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/delete", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			result := map[string]string{"error": "some error"}
			jresult, _ := json.Marshal(result)

			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, string(jresult))
		}))
		defer server.Close()

		c, _ := NewClient(http.DefaultClient, Config{Host: server.URL})
		svc := &FnService{Client: c}

		_, err := svc.Delete(testCtx, testFn, testNs)

		require.Error(t, err)
	})
}
