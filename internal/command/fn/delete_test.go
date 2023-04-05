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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"

	openapi "github.com/funlessdev/fl-client-sdk-go"
)

func TestFnDelete(t *testing.T) {
	testFn := "test-fn"
	testMod := "test-mod"
	testCtx := context.Background()
	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).Build()

	t.Run("should use FnService.Delete to delete functions", func(t *testing.T) {
		cmd := Delete{
			Name:   testFn,
			Module: testMod,
		}

		mockFnHandler := mocks.NewFnHandler(t)
		mockFnHandler.On("Delete", mock.Anything, testFn, testMod).Return(nil)

		err := cmd.Run(testCtx, mockFnHandler, testLogger, &Fn{})
		require.NoError(t, err)
		mockFnHandler.AssertCalled(t, "Delete", mock.Anything, testFn, testMod)
		mockFnHandler.AssertNumberOfCalls(t, "Delete", 1)
		mockFnHandler.AssertExpectations(t)
	})
	t.Run("should correctly print result", func(t *testing.T) {
		cmd := Delete{
			Name:   testFn,
			Module: testMod,
		}

		mockFnHandler := mocks.NewFnHandler(t)
		mockFnHandler.On("Delete", mock.Anything, testFn, testMod).Return(nil)

		var outbuf bytes.Buffer
		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).Build()

		err := cmd.Run(testCtx, mockFnHandler, bufLogger, &Fn{})

		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("\nSuccessfully deleted function %s/%s.\n", testMod, testFn), (&outbuf).String())
		mockFnHandler.AssertExpectations(t)
	})

	t.Run("should return error if invalid delete request", func(t *testing.T) {
		cmd := Delete{
			Name:   testFn,
			Module: testMod,
		}

		mockFnHandler := mocks.NewFnHandler(t)

		e := &openapi.GenericOpenAPIError{}
		mockFnHandler.On("Delete", mock.Anything, testFn, testMod).Return(e)

		err := cmd.Run(testCtx, mockFnHandler, testLogger, &Fn{})
		require.Error(t, err)
	})
}
