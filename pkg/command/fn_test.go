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
package command

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/funlessdev/funless-cli/pkg/client"
	"github.com/funlessdev/funless-cli/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFn(t *testing.T) {
	testfn := "test-fn"

	t.Run("should use FnService.Invoke to invoke functions", func(t *testing.T) {
		cmd := Fn{Name: testfn}
		mockInvoker := mocks.NewFnHandler(t)

		mockInvoker.On("Invoke", testfn).Return(&http.Response{}, nil)

		err := cmd.Run(mockInvoker)
		require.NoError(t, err)
		mockInvoker.AssertCalled(t, "Invoke", testfn)
		mockInvoker.AssertNumberOfCalls(t, "Invoke", 1)
		mockInvoker.AssertExpectations(t)
	})

	t.Run("should return error if invalid invoke request", func(t *testing.T) {
		cmd := Fn{Name: testfn}
		mockInvoker := mocks.NewFnHandler(t)
		mockInvoker.On("Invoke", testfn).Return(nil, fmt.Errorf("some error in FnService.Invoke"))

		err := cmd.Run(mockInvoker)
		require.Error(t, err)
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/_/fn/"+testfn, r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, _ := client.NewClient(http.DefaultClient, client.Config{Host: server.URL})
	svc := &client.FnService{Client: c}

	t.Run("should send invoke request to server", func(t *testing.T) {
		cmd := Fn{Name: testfn}
		err := cmd.Run(svc)
		require.NoError(t, err)
	})

}
