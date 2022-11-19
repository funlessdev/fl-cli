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

func TestFnCreateNoBuild(t *testing.T) {
	testResult := "test-fn"
	testFn := "test-fn"
	testNs := "test-ns"
	testSource, _ := filepath.Abs("../../../test/fixtures/real.wasm")
	testCtx := context.Background()
	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).DisableAnimation().Build()

	mockBuilder := mocks.NewDockerBuilder(t)

	t.Run("should use FnService.Create to create functions", func(t *testing.T) {
		cmd := Create{
			Name:       testFn,
			Namespace:  testNs,
			SourceFile: testSource,
			NoBuild:    true,
		}

		mockInvoker := mocks.NewFnHandler(t)
		mockInvoker.On("Create", testCtx, testFn, testNs, mock.Anything).Return(openapi.FunctionCreationSuccess{Result: &testResult}, nil)

		err := cmd.Run(testCtx, mockBuilder, mockInvoker, testLogger)
		require.NoError(t, err)
		mockInvoker.AssertCalled(t, "Create", testCtx, testFn, testNs, mock.AnythingOfType("*os.File"))
		mockInvoker.AssertNumberOfCalls(t, "Create", 1)
		mockInvoker.AssertExpectations(t)
	})

	t.Run("should correctly print result with single file", func(t *testing.T) {
		cmd := Create{
			Name:       testFn,
			Namespace:  testNs,
			SourceFile: testSource,
			NoBuild:    true,
		}

		mockInvoker := mocks.NewFnHandler(t)
		mockInvoker.On("Create", testCtx, testFn, testNs, mock.Anything).Return(openapi.FunctionCreationSuccess{Result: &testResult}, nil)

		var outbuf bytes.Buffer

		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).Build()

		err := cmd.Run(testCtx, mockBuilder, mockInvoker, bufLogger)

		require.NoError(t, err)
		assert.Equal(t, testResult+"\n", (&outbuf).String())
		mockInvoker.AssertExpectations(t)
	})

	t.Run("should return error if given an invalid create request", func(t *testing.T) {
		cmd := Create{
			Name:       testFn,
			Namespace:  testNs,
			SourceFile: testSource,
			NoBuild:    true,
		}

		mockInvoker := mocks.NewFnHandler(t)
		e := &openapi.GenericOpenAPIError{}
		mockInvoker.On("Create", testCtx, testFn, testNs, mock.Anything).Return(openapi.FunctionCreationSuccess{}, e)

		err := cmd.Run(testCtx, mockBuilder, mockInvoker, testLogger)
		require.Error(t, err)
	})

	t.Run("should return error if source file does not exist", func(t *testing.T) {
		cmd := Create{
			Name:       testFn,
			Namespace:  testNs,
			SourceFile: "no_file",
			NoBuild:    true,
		}

		mockInvoker := mocks.NewFnHandler(t)

		err := cmd.Run(testCtx, mockBuilder, mockInvoker, testLogger)
		require.Error(t, err)
	})
}

func TestFnCreateBuild(t *testing.T) {
	testResult := `Building the given function using fl-runtimes...

Setting up...
done
Pulling builder image for js 📦
done
Building source using builder image 🛠️
done
test-fn
`
	testFn := "test-fn"
	testNs := "test-ns"
	testLanguage := "js"
	testSource, _ := filepath.Abs("../../../test/fixtures/real.wasm")
	testDir, _ := filepath.Abs("../../../test/fixtures/test_dir/")
	testOutDir, _ := filepath.Abs("../../../test/fixtures")
	ctx := context.Background()
	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).DisableAnimation().Build()

	mockFnHandler := mocks.NewFnHandler(t)
	mockBuilder := mocks.NewDockerBuilder(t)
	mockBuilder.On("Setup", mock.Anything, testLanguage, testOutDir).Return(nil)
	mockBuilder.On("PullBuilderImage", ctx).Return(nil)
	mockBuilder.On("BuildSource", ctx, testDir).Return(nil)
	mockBuilder.On("GetWasmFile", testFn).Return(nil, nil)

	t.Run("should use FnService.Create to create functions", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			SourceDir: testDir,
			Language:  testLanguage,
			OutDir:    testOutDir,
		}

		mockFnHandler.On("Create", ctx, testFn, testNs, mock.Anything).Return(openapi.FunctionCreationSuccess{Result: &testFn}, nil)

		err := cmd.Run(ctx, mockBuilder, mockFnHandler, testLogger)
		require.NoError(t, err)

		mockFnHandler.AssertCalled(t, "Create", ctx, testFn, testNs, mock.AnythingOfType("*os.File"))
		mockFnHandler.AssertNumberOfCalls(t, "Create", 1)
		mockFnHandler.AssertExpectations(t)
	})

	t.Run("should use DockerBuilder.BuildSource to build functions", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			SourceDir: testDir,
			Language:  testLanguage,
			OutDir:    testOutDir,
		}

		mockFnHandler := mocks.NewFnHandler(t)
		mockFnHandler.On("Create", ctx, testFn, testNs, mock.Anything).Return(openapi.FunctionCreationSuccess{Result: &testFn}, nil)

		err := cmd.Run(ctx, mockBuilder, mockFnHandler, testLogger)
		require.NoError(t, err)

		mockBuilder.AssertCalled(t, "BuildSource", ctx, testDir)
		mockBuilder.AssertNumberOfCalls(t, "BuildSource", 2)
		mockBuilder.AssertExpectations(t)
	})

	t.Run("should correctly print result when building from a directory", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			SourceDir: testDir,
			Language:  testLanguage,
			OutDir:    testOutDir,
		}

		mockFnHandler := mocks.NewFnHandler(t)
		mockFnHandler.On("Create", ctx, testFn, testNs, mock.Anything).Return(openapi.FunctionCreationSuccess{Result: &testFn}, nil)

		var outbuf bytes.Buffer

		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).DisableAnimation().Build()

		err := cmd.Run(ctx, mockBuilder, mockFnHandler, bufLogger)

		require.NoError(t, err)
		assert.Equal(t, testResult, (&outbuf).String())
		mockFnHandler.AssertExpectations(t)
		mockBuilder.AssertExpectations(t)
	})

	t.Run("should return error if asked to build a single source file", func(t *testing.T) {
		cmd := Create{
			Name:       testFn,
			Namespace:  testNs,
			SourceFile: testSource,
			Language:   testLanguage,
		}

		mockFnHandler := mocks.NewFnHandler(t)

		mockBuilder := mocks.NewDockerBuilder(t)
		err := cmd.Run(ctx, mockBuilder, mockFnHandler, testLogger)
		require.Error(t, err)
	})

	t.Run("should return error if builder setup encounters errors", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			SourceDir: testDir,
			Language:  testLanguage,
			OutDir:    testOutDir,
		}

		mockFnHandler := mocks.NewFnHandler(t)

		mockBuilder := mocks.NewDockerBuilder(t)
		mockBuilder.On("Setup", mock.Anything, testLanguage, testOutDir).Return(errors.New("some error")).Once()

		err := cmd.Run(ctx, mockBuilder, mockFnHandler, testLogger)
		require.Error(t, err)
		mockBuilder.AssertExpectations(t)
	})

	t.Run("should return error if builder image cannot be pulled", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			SourceDir: testDir,
			Language:  testLanguage,
			OutDir:    testOutDir,
		}

		mockFnHandler := mocks.NewFnHandler(t)

		mockBuilder := mocks.NewDockerBuilder(t)
		mockBuilder.On("Setup", mock.Anything, testLanguage, testOutDir).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", ctx).Return(errors.New("some error")).Once()

		err := cmd.Run(ctx, mockBuilder, mockFnHandler, testLogger)
		require.Error(t, err)
		mockBuilder.AssertExpectations(t)
	})

	t.Run("should return error if builder image encounters errors", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			SourceDir: testDir,
			Language:  testLanguage,
			OutDir:    testOutDir,
		}

		mockFnHandler := mocks.NewFnHandler(t)

		mockBuilder := mocks.NewDockerBuilder(t)
		mockBuilder.On("Setup", mock.Anything, testLanguage, testOutDir).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", ctx).Return(nil).Once()
		mockBuilder.On("BuildSource", ctx, testDir).Return(errors.New("some error")).Once()

		err := cmd.Run(ctx, mockBuilder, mockFnHandler, testLogger)
		require.Error(t, err)
		mockBuilder.AssertExpectations(t)
	})
}
