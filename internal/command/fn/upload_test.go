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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFnUpload(t *testing.T) {
	testFn := "test-fn"
	testNs := "test-ns"
	testSource, _ := filepath.Abs("../../../test/fixtures/real.wasm")
	ctx := context.Background()
	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).DisableAnimation().Build()

	openWasmFile = readWasmFile

	t.Run("should return error if file is not valid", func(t *testing.T) {
		upload := Upload{
			Name:      testFn,
			Source:    "not-found.wasm",
			Namespace: testNs,
		}
		mockFnHandler := mocks.NewFnHandler(t)
		err := upload.Run(ctx, mockFnHandler, testLogger)
		require.Error(t, err)
	})

	t.Run("should return error when FnService.Create fails", func(t *testing.T) {
		cmd := Upload{
			Name:      testFn,
			Namespace: testNs,
			Source:    testSource,
		}

		mockFnHandler := mocks.NewFnHandler(t)
		mockFnHandler.On("Create", ctx, testFn, testNs, mock.Anything).Return(errors.New("error")).Once()

		err := cmd.Run(ctx, mockFnHandler, testLogger)
		require.Error(t, err)

		mockFnHandler.AssertCalled(t, "Create", ctx, testFn, testNs, mock.AnythingOfType("*os.File"))
		mockFnHandler.AssertNumberOfCalls(t, "Create", 1)
		mockFnHandler.AssertExpectations(t)
	})

	t.Run("should correctly print result when it works", func(t *testing.T) {
		testResult := `Reading wasm...
done
Uploading function...
done
Successfully uploaded function test-ns/test-fn 👌
`
		cmd := Upload{
			Name:      testFn,
			Namespace: testNs,
			Source:    testSource,
		}

		mockFnHandler := mocks.NewFnHandler(t)
		mockFnHandler.On("Create", ctx, testFn, testNs, mock.Anything).Return(nil)

		var outbuf bytes.Buffer

		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).DisableAnimation().Build()
		err := cmd.Run(ctx, mockFnHandler, bufLogger)

		require.NoError(t, err)
		assert.Equal(t, testResult, (&outbuf).String())
		mockFnHandler.AssertExpectations(t)
	})

}

func Test_obtainWasmFile(t *testing.T) {

	t.Run("should return error if no file is found", func(t *testing.T) {
		_, err := openWasmFile("not-exist.wasm")
		require.Error(t, err)
	})

	t.Run("should return error if file is not a wasm file", func(t *testing.T) {
		_, err := openWasmFile("not-wasm.file")
		require.Error(t, err)
	})

	t.Run("should return error if file does not contain wasm magic header", func(t *testing.T) {
		_, err := openWasmFile("../../../test/fixtures/fake.wasm")
		require.Error(t, err)
	})

	t.Run("should return file when given a valid wasm", func(t *testing.T) {
		file, err := openWasmFile("../../../test/fixtures/real.wasm")
		require.NoError(t, err)
		require.NotNil(t, file)
	})

}
