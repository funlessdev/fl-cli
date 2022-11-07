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
	testSource, _ := filepath.Abs("../../../test/fixtures/code.wasm")
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
	testResult :=
		`Building the given function using fl-runtimes...

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
	testSource, _ := filepath.Abs("../../../test/fixtures/test_code.txt")
	testDir, _ := filepath.Abs("../../../test/fixtures/test_dir/")
	testOutDir, _ := filepath.Abs("../../../test/fixtures")
	testCtx := context.Background()
	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).DisableAnimation().Build()

	mockBuilder := mocks.NewDockerBuilder(t)

	t.Run("should use FnService.Create to create functions", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			SourceDir: testDir,
			Language:  testLanguage,
			OutDir:    testOutDir,
		}

		mockInvoker := mocks.NewFnHandler(t)
		mockInvoker.On("Create", testCtx, testFn, testNs, mock.Anything).Return(openapi.FunctionCreationSuccess{Result: &testFn}, nil)

		mockBuilder := mocks.NewDockerBuilder(t)
		mockBuilder.On("Setup", testCtx, testLanguage, testOutDir).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", testCtx).Return(nil).Once()
		mockBuilder.On("BuildSource", testCtx, testDir).Return(nil).Once()

		err := cmd.Run(testCtx, mockBuilder, mockInvoker, testLogger)
		require.NoError(t, err)

		mockInvoker.AssertCalled(t, "Create", testCtx, testFn, testNs, mock.AnythingOfType("*os.File"))
		mockInvoker.AssertNumberOfCalls(t, "Create", 1)
		mockInvoker.AssertExpectations(t)
	})

	t.Run("should use DockerBuilder.BuildSource to build functions", func(t *testing.T) {
		cmd := Create{
			Name:      testFn,
			Namespace: testNs,
			SourceDir: testDir,
			Language:  testLanguage,
			OutDir:    testOutDir,
		}

		mockInvoker := mocks.NewFnHandler(t)
		mockInvoker.On("Create", testCtx, testFn, testNs, mock.Anything).Return(openapi.FunctionCreationSuccess{Result: &testFn}, nil)

		mockBuilder := mocks.NewDockerBuilder(t)
		mockBuilder.On("Setup", testCtx, testLanguage, testOutDir).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", testCtx).Return(nil).Once()
		mockBuilder.On("BuildSource", testCtx, testDir).Return(nil).Once()

		err := cmd.Run(testCtx, mockBuilder, mockInvoker, testLogger)
		require.NoError(t, err)

		mockBuilder.AssertCalled(t, "BuildSource", testCtx, testDir)
		mockBuilder.AssertNumberOfCalls(t, "BuildSource", 1)
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

		mockInvoker := mocks.NewFnHandler(t)
		mockInvoker.On("Create", testCtx, testFn, testNs, mock.Anything).Return(openapi.FunctionCreationSuccess{Result: &testFn}, nil)

		mockBuilder.On("Setup", testCtx, testLanguage, testOutDir).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", testCtx).Return(nil).Once()
		mockBuilder.On("BuildSource", testCtx, testDir).Return(nil).Once()

		var outbuf bytes.Buffer

		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).DisableAnimation().Build()

		err := cmd.Run(testCtx, mockBuilder, mockInvoker, bufLogger)

		require.NoError(t, err)
		assert.Equal(t, testResult, (&outbuf).String())
		mockInvoker.AssertExpectations(t)
		mockBuilder.AssertExpectations(t)
	})

	t.Run("should return error if asked to build a single source file", func(t *testing.T) {
		cmd := Create{
			Name:       testFn,
			Namespace:  testNs,
			SourceFile: testSource,
			Language:   testLanguage,
		}

		mockInvoker := mocks.NewFnHandler(t)

		err := cmd.Run(testCtx, mockBuilder, mockInvoker, testLogger)
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

		mockInvoker := mocks.NewFnHandler(t)

		mockBuilder.On("Setup", testCtx, testLanguage, testOutDir).Return(errors.New("some error")).Once()

		err := cmd.Run(testCtx, mockBuilder, mockInvoker, testLogger)
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

		mockInvoker := mocks.NewFnHandler(t)

		mockBuilder.On("Setup", testCtx, testLanguage, testOutDir).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", testCtx).Return(errors.New("some error")).Once()

		err := cmd.Run(testCtx, mockBuilder, mockInvoker, testLogger)
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

		mockInvoker := mocks.NewFnHandler(t)

		mockBuilder.On("Setup", testCtx, testLanguage, testOutDir).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", testCtx).Return(nil).Once()
		mockBuilder.On("BuildSource", testCtx, testDir).Return(errors.New("some error")).Once()

		err := cmd.Run(testCtx, mockBuilder, mockInvoker, testLogger)
		require.Error(t, err)
		mockBuilder.AssertExpectations(t)
	})
}
