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

package fn

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"

	openapi "github.com/funlessdev/fl-client-sdk-go"
)

func TestFnInvoke(t *testing.T) {
	testResult := "test-res"
	testFn := "test-fn"
	testMod := "test-mod"
	testArgs := map[string]string{"name": "Some name"}
	testJArgs := "{\"name\":\"Some name\"}"
	testParsedJArgs := map[string]interface{}{"name": "Some name"}
	testCtx := context.Background()
	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).Build()

	t.Run("should use FnService.Invoke to invoke functions", func(t *testing.T) {
		cmd := Invoke{
			Name:   testFn,
			Module: testMod,
			Args:   map[string]string{},
		}

		mockFnHandler := mocks.NewFnHandler(t)
		mockFnHandler.On("Invoke", testCtx, testFn, testMod, map[string]interface{}{}).Return(pkg.IvkResult{Result: testResult}, nil)

		err := cmd.Run(testCtx, mockFnHandler, testLogger)
		require.NoError(t, err)
		mockFnHandler.AssertCalled(t, "Invoke", testCtx, testFn, testMod, map[string]interface{}{})
		mockFnHandler.AssertNumberOfCalls(t, "Invoke", 1)
		mockFnHandler.AssertExpectations(t)
	})

	t.Run("should correctly print result", func(t *testing.T) {
		cmd := Invoke{
			Name:   testFn,
			Module: testMod,
			Args:   map[string]string{},
		}

		mockFnHandler := mocks.NewFnHandler(t)
		mockFnHandler.On("Invoke", testCtx, testFn, testMod, map[string]interface{}{}).Return(pkg.IvkResult{Result: testResult}, nil)

		var outbuf bytes.Buffer
		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).Build()

		err := cmd.Run(testCtx, mockFnHandler, bufLogger)

		require.NoError(t, err)
		assert.Equal(t, testResult, (&outbuf).String())
		mockFnHandler.AssertExpectations(t)
	})

	t.Run("should correctly parse and forward keyword args", func(t *testing.T) {
		cmd := Invoke{
			Name:   testFn,
			Module: testMod,
			Args:   testArgs,
		}

		mockArgs := make(map[string]interface{}, len(testArgs))
		for k, v := range testArgs {
			mockArgs[k] = v
		}

		mockFnHandler := mocks.NewFnHandler(t)
		mockFnHandler.On("Invoke", testCtx, testFn, testMod, mockArgs).Return(pkg.IvkResult{Result: testResult}, nil)

		err := cmd.Run(testCtx, mockFnHandler, testLogger)
		require.NoError(t, err)
		mockFnHandler.AssertCalled(t, "Invoke", testCtx, testFn, testMod, mockArgs)
		mockFnHandler.AssertNumberOfCalls(t, "Invoke", 1)
		mockFnHandler.AssertExpectations(t)
	})

	t.Run("should correctly parse and forward json args", func(t *testing.T) {
		cmd := Invoke{
			Name:     testFn,
			Module:   testMod,
			JsonArgs: testJArgs,
		}

		mockFnHandler := mocks.NewFnHandler(t)
		mockFnHandler.On("Invoke", testCtx, testFn, testMod, testParsedJArgs).Return(
			pkg.IvkResult{Result: testResult}, nil)

		err := cmd.Run(testCtx, mockFnHandler, testLogger)
		require.NoError(t, err)
		mockFnHandler.AssertCalled(t, "Invoke", testCtx, testFn, testMod, testParsedJArgs)
		mockFnHandler.AssertNumberOfCalls(t, "Invoke", 1)
		mockFnHandler.AssertExpectations(t)
	})

	t.Run("should return error if invalid invoke request", func(t *testing.T) {
		// missing module
		cmd := Invoke{
			Name: testFn,
		}

		mockFnHandler := mocks.NewFnHandler(t)
		e := &openapi.GenericOpenAPIError{}
		mockFnHandler.On("Invoke", testCtx, testFn, "", mock.Anything).Return(
			pkg.IvkResult{}, e)

		err := cmd.Run(testCtx, mockFnHandler, testLogger)
		require.Error(t, err)
	})
}
