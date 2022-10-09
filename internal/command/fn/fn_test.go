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
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"testing/fstest"

	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/funlessdev/fl-cli/test/mocks"
	swagger "github.com/funlessdev/fl-client-sdk-go"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestFnInvoke(t *testing.T) {
	testResult := map[string]interface{}{"payload": "Hi"}
	testFn := "test-fn"
	testNs := "test-ns"
	testArgs := map[string]string{"name": "Some name"}
	testJArgs := "{\"name\":\"Some name\"}"
	testParsedJArgs := map[string]interface{}{"name": "Some name"}
	testCtx := context.Background()
	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).Build()

	t.Run("should use FnService.Invoke to invoke functions", func(t *testing.T) {
		cmd := Invoke{
			Name:      testFn,
			Namespace: testNs,
			Args:      map[string]string{},
		}
		mockInvoker := mocks.NewFnHandler(t)

		mockInvoker.On("Invoke", testCtx, testFn, testNs, map[string]interface{}{}).Return(swagger.FunctionInvocationSuccess{Result: testResult}, nil)

		err := cmd.Run(testCtx, mockInvoker, testLogger)
		require.NoError(t, err)
		mockInvoker.AssertCalled(t, "Invoke", testCtx, testFn, testNs, map[string]interface{}{})
		mockInvoker.AssertNumberOfCalls(t, "Invoke", 1)
		mockInvoker.AssertExpectations(t)
	})

	t.Run("should correctly print result", func(t *testing.T) {
		cmd := Invoke{
			Name:      testFn,
			Namespace: testNs,
			Args:      map[string]string{},
		}
		mockInvoker := mocks.NewFnHandler(t)

		mockInvoker.On("Invoke", testCtx, testFn, testNs, map[string]interface{}{}).Return(swagger.FunctionInvocationSuccess{Result: testResult}, nil)

		var outbuf bytes.Buffer
		var testOutput, _ = json.Marshal(testResult)
		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).Build()

		err := cmd.Run(testCtx, mockInvoker, bufLogger)

		require.NoError(t, err)
		assert.Equal(t, string(testOutput)+"\n", (&outbuf).String())
		mockInvoker.AssertExpectations(t)
	})

	t.Run("should correctly parse and forward keyword args", func(t *testing.T) {
		cmd := Invoke{
			Name:      testFn,
			Namespace: testNs,
			Args:      testArgs,
		}
		mockInvoker := mocks.NewFnHandler(t)

		mockArgs := make(map[string]interface{}, len(testArgs))
		for k, v := range testArgs {
			mockArgs[k] = v
		}
		mockInvoker.On("Invoke", testCtx, testFn, testNs, mockArgs).Return(swagger.FunctionInvocationSuccess{Result: testResult}, nil)

		err := cmd.Run(testCtx, mockInvoker, testLogger)
		require.NoError(t, err)
		mockInvoker.AssertCalled(t, "Invoke", testCtx, testFn, testNs, mockArgs)
		mockInvoker.AssertNumberOfCalls(t, "Invoke", 1)
		mockInvoker.AssertExpectations(t)
	})

	t.Run("should correctly parse and forward json args", func(t *testing.T) {
		cmd := Invoke{
			Name:      testFn,
			Namespace: testNs,
			JsonArgs:  testJArgs,
		}
		mockInvoker := mocks.NewFnHandler(t)

		mockInvoker.On("Invoke", testCtx, testFn, testNs, testParsedJArgs).Return(swagger.FunctionInvocationSuccess{Result: testResult}, nil)

		err := cmd.Run(testCtx, mockInvoker, testLogger)
		require.NoError(t, err)
		mockInvoker.AssertCalled(t, "Invoke", testCtx, testFn, testNs, testParsedJArgs)
		mockInvoker.AssertNumberOfCalls(t, "Invoke", 1)
		mockInvoker.AssertExpectations(t)
	})

	t.Run("should return error if invalid invoke request", func(t *testing.T) {
		cmd := Invoke{
			Name: testFn,
		}
		mockInvoker := mocks.NewFnHandler(t)
		mockInvoker.On("Invoke", testCtx, testFn, "", map[string]interface{}{}).Return(swagger.FunctionInvocationSuccess{}, fmt.Errorf("some error in FnService.Invoke"))

		err := cmd.Run(testCtx, mockInvoker, testLogger)
		require.Error(t, err)
	})

}

func TestFnCreate(t *testing.T) {
	testResult := "test-fn"
	testFn := "test-fn"
	testNs := "test-ns"
	testSource := "test.js"
	testLanguage := "nodejs"
	testCode := []byte("console.log('Something')")
	testFS := fstest.MapFS{
		"test.js": {
			Data: testCode,
		},
	}
	testCtx := context.Background()
	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).Build()

	t.Run("should use FnService.Create to create functions", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			Source:    testSource,
			Language:  testLanguage,
			FS:        testFS,
		}
		mockInvoker := mocks.NewFnHandler(t)

		mockInvoker.On("Create", testCtx, testFn, testNs, string(testCode), testLanguage).Return(swagger.FunctionCreationSuccess{Result: &testResult}, nil)

		err := cmd.Run(testCtx, mockInvoker, testLogger)
		require.NoError(t, err)
		mockInvoker.AssertCalled(t, "Create", testCtx, testFn, testNs, string(testCode), testLanguage)
		mockInvoker.AssertNumberOfCalls(t, "Create", 1)
		mockInvoker.AssertExpectations(t)
	})

	t.Run("should correctly print result", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			Source:    testSource,
			Language:  testLanguage,
			FS:        testFS,
		}
		mockInvoker := mocks.NewFnHandler(t)

		mockInvoker.On("Create", testCtx, testFn, testNs, string(testCode), testLanguage).Return(swagger.FunctionCreationSuccess{Result: &testResult}, nil)

		var outbuf bytes.Buffer

		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).Build()

		err := cmd.Run(testCtx, mockInvoker, bufLogger)

		require.NoError(t, err)
		assert.Equal(t, testResult+"\n", (&outbuf).String())
		mockInvoker.AssertExpectations(t)
	})

	t.Run("should return error if invalid create request", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			Source:    testSource,
			Language:  testLanguage,
			FS:        testFS,
		}
		mockInvoker := mocks.NewFnHandler(t)
		mockInvoker.On("Create", testCtx, testFn, testNs, string(testCode), testLanguage).Return(swagger.FunctionCreationSuccess{}, fmt.Errorf("some error in FnService.Invoke"))

		err := cmd.Run(testCtx, mockInvoker, testLogger)
		require.Error(t, err)
	})

	t.Run("should return error if source file does not exist", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			Source:    "no_file",
			Language:  testLanguage,
			FS:        testFS,
		}
		mockInvoker := mocks.NewFnHandler(t)

		err := cmd.Run(testCtx, mockInvoker, testLogger)
		require.Error(t, err)
	})

}

func TestFnDelete(t *testing.T) {
	testResult := "test-fn"
	testFn := "test-fn"
	testNs := "test-ns"
	testCtx := context.Background()
	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).Build()

	t.Run("should use FnService.Delete to delete functions", func(t *testing.T) {
		cmd := Delete{
			Name:      testFn,
			Namespace: testNs,
		}
		mockInvoker := mocks.NewFnHandler(t)

		mockInvoker.On("Delete", testCtx, testFn, testNs).Return(swagger.FunctionDeletionSuccess{Result: &testResult}, nil)

		err := cmd.Run(testCtx, mockInvoker, testLogger)
		require.NoError(t, err)
		mockInvoker.AssertCalled(t, "Delete", testCtx, testFn, testNs)
		mockInvoker.AssertNumberOfCalls(t, "Delete", 1)
		mockInvoker.AssertExpectations(t)
	})
	t.Run("should correctly print result", func(t *testing.T) {
		cmd := Delete{
			Name:      testFn,
			Namespace: testNs,
		}
		mockInvoker := mocks.NewFnHandler(t)

		mockInvoker.On("Delete", testCtx, testFn, testNs).Return(swagger.FunctionDeletionSuccess{Result: &testResult}, nil)

		var outbuf bytes.Buffer
		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).Build()

		err := cmd.Run(testCtx, mockInvoker, bufLogger)

		require.NoError(t, err)
		assert.Equal(t, testResult+"\n", (&outbuf).String())
		mockInvoker.AssertExpectations(t)
	})

	t.Run("should return error if invalid delete request", func(t *testing.T) {
		cmd := Delete{
			Name:      testFn,
			Namespace: testNs,
		}
		mockInvoker := mocks.NewFnHandler(t)
		mockInvoker.On("Delete", testCtx, testFn, testNs).Return(swagger.FunctionDeletionSuccess{}, fmt.Errorf("some error in FnService.Invoke"))

		err := cmd.Run(testCtx, mockInvoker, testLogger)
		require.Error(t, err)
	})

}
