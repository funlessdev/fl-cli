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
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"

	openapi "github.com/funlessdev/fl-client-sdk-go"
)

func TestFnCreate(t *testing.T) {
	testResult := `Creating test-fn function...

Building function...🏗 ️
done
Uploading function... 📮
done

Successfully created function test-ns/test-fn.
`
	testFn := "test-fn"
	testNs := "test-ns"
	testLanguage := "js"
	testDir, _ := filepath.Abs("../../../test/fixtures/test_dir/")
	ctx := context.Background()
	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).DisableAnimation().Build()

	mockFnHandler := mocks.NewFnHandler(t)
	mockFnHandler.On("Create", ctx, testFn, testNs, mock.Anything).Return(openapi.FunctionCreationSuccess{Result: &testFn}, nil)

	mockBuilder := mocks.NewDockerBuilder(t)
	mockBuilder.On("Setup", mock.Anything, testLanguage, mock.Anything).Return(nil)
	mockBuilder.On("PullBuilderImage", ctx).Return(nil)
	mockBuilder.On("BuildSource", ctx, testDir).Return(nil)

	// monkey patch the openWasmFile function
	openWasmFile = func(path string) (*os.File, error) {
		return &os.File{}, nil
	}

	t.Run("success: should correctly print result when building from a directory", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			Source:    testDir,
			Language:  testLanguage,
		}

		var outbuf bytes.Buffer

		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).DisableAnimation().Build()

		err := cmd.Run(ctx, mockBuilder, mockFnHandler, bufLogger)

		require.NoError(t, err)
		assert.Equal(t, testResult, (&outbuf).String())
		mockFnHandler.AssertCalled(t, "Create", ctx, testFn, testNs, mock.AnythingOfType("*os.File"))
		mockFnHandler.AssertNumberOfCalls(t, "Create", 1)
		mockFnHandler.AssertExpectations(t)
		mockBuilder.AssertExpectations(t)
	})

	t.Run("should use DockerBuilder.BuildSource to build functions", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			Source:    testDir,
			Language:  testLanguage,
		}

		err := cmd.Run(ctx, mockBuilder, mockFnHandler, testLogger)
		require.NoError(t, err)

		mockBuilder.AssertCalled(t, "BuildSource", ctx, testDir)
		mockBuilder.AssertNumberOfCalls(t, "BuildSource", 2)
		mockBuilder.AssertExpectations(t)
	})

	t.Run("should return error if builder setup encounters errors", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			Source:    testDir,
			Language:  testLanguage,
		}

		mockBuilder := mocks.NewDockerBuilder(t)
		mockBuilder.On("Setup", mock.Anything, testLanguage, mock.Anything).Return(errors.New("some error")).Once()

		err := cmd.Run(ctx, mockBuilder, mockFnHandler, testLogger)
		require.Error(t, err)
		mockBuilder.AssertExpectations(t)
	})

	t.Run("should return error if builder image cannot be pulled", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			Source:    testDir,
			Language:  testLanguage,
		}

		mockBuilder := mocks.NewDockerBuilder(t)
		mockBuilder.On("Setup", mock.Anything, testLanguage, mock.Anything).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", ctx).Return(errors.New("some error")).Once()

		err := cmd.Run(ctx, mockBuilder, mockFnHandler, testLogger)
		require.Error(t, err)
		mockBuilder.AssertExpectations(t)
	})

	t.Run("should return error if builder image encounters errors", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			Source:    testDir,
			Language:  testLanguage,
		}

		mockBuilder := mocks.NewDockerBuilder(t)
		mockBuilder.On("Setup", mock.Anything, testLanguage, mock.Anything).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", ctx).Return(nil).Once()
		mockBuilder.On("BuildSource", ctx, testDir).Return(errors.New("some error")).Once()

		err := cmd.Run(ctx, mockBuilder, mockFnHandler, testLogger)
		require.Error(t, err)
		mockBuilder.AssertExpectations(t)
	})
}
