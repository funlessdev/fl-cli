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
	"fmt"
	"os"
	"testing"

	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"

	openapi "github.com/funlessdev/fl-client-sdk-go"
)

func TestFnList(t *testing.T) {
	testNs := "test-ns"
	testCtx := context.Background()
	testResult := []string{"f1", "f2"}
	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).Build()

	t.Run("should use FnService.List to list functions", func(t *testing.T) {
		cmd := List{
			Namespace: testNs,
		}

		mockFnHandler := mocks.NewFnHandler(t)
		mockFnHandler.On("List", testCtx, testNs).Return(openapi.FunctionListSuccess{Result: testResult}, nil)

		err := cmd.Run(testCtx, mockFnHandler, testLogger)
		require.NoError(t, err)
		mockFnHandler.AssertCalled(t, "List", testCtx, testNs)
		mockFnHandler.AssertNumberOfCalls(t, "List", 1)
		mockFnHandler.AssertExpectations(t)
	})
	t.Run("should correctly print result", func(t *testing.T) {
		cmd := List{
			Namespace: testNs,
		}

		mockFnHandler := mocks.NewFnHandler(t)
		mockFnHandler.On("List", testCtx, testNs).Return(openapi.FunctionListSuccess{Result: testResult}, nil)

		var outbuf bytes.Buffer
		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).Build()

		err := cmd.Run(testCtx, mockFnHandler, bufLogger)
		expected := fmt.Sprintf("%s\n%s\n", testResult[0], testResult[1])

		require.NoError(t, err)
		assert.Equal(t, expected, (&outbuf).String())
		mockFnHandler.AssertExpectations(t)
	})

	t.Run("should correctly print result when asked to count returned functions", func(t *testing.T) {
		cmd := List{
			Namespace: testNs,
			Count:     true,
		}

		mockFnHandler := mocks.NewFnHandler(t)
		mockFnHandler.On("List", testCtx, testNs).Return(openapi.FunctionListSuccess{Result: testResult}, nil)

		var outbuf bytes.Buffer
		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).Build()

		err := cmd.Run(testCtx, mockFnHandler, bufLogger)
		expected := fmt.Sprintf("%s\n%s\nCount: %d\n", testResult[0], testResult[1], len(testResult))

		require.NoError(t, err)
		assert.Equal(t, expected, (&outbuf).String())
		mockFnHandler.AssertExpectations(t)
	})

	t.Run("should return error if the list request is invalid", func(t *testing.T) {
		cmd := List{
			Namespace: testNs,
		}

		mockFnHandler := mocks.NewFnHandler(t)

		e := &openapi.GenericOpenAPIError{}
		mockFnHandler.On("List", testCtx, testNs).Return(openapi.FunctionListSuccess{}, e)

		err := cmd.Run(testCtx, mockFnHandler, testLogger)
		require.Error(t, err)
	})
}
