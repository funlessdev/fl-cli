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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFnBuild(t *testing.T) {

	testFn := "test-fn"
	testLanguage := "js"
	testDir, _ := filepath.Abs("../../../test/fixtures/test_dir/")
	testOutDir, _ := filepath.Abs("../../../test/fixtures")

	ctx := context.TODO()

	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).DisableAnimation().Build()

	mockBuilder := mocks.NewDockerBuilder(t)
	mockBuilder.On("RenameCodeWasm", mock.Anything).Return(nil)

	t.Run("should use DockerBuilder to build functions", func(t *testing.T) {
		cmd := Build{
			Name:        testFn,
			Source:      testDir,
			Destination: testOutDir,
			Language:    testLanguage,
		}

		mockBuilder.On("Setup", mock.Anything, testLanguage, testOutDir).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", ctx).Return(nil).Once()
		mockBuilder.On("BuildSource", ctx, testDir).Return(nil).Once()

		err := cmd.Run(ctx, mockBuilder, testLogger)
		require.NoError(t, err)

		mockBuilder.AssertNumberOfCalls(t, "BuildSource", 1)
		mockBuilder.AssertExpectations(t)
	})

	t.Run("should correctly print result when building from a directory", func(t *testing.T) {
		cmd := Build{
			Name:        testFn,
			Source:      testDir,
			Destination: testOutDir,
			Language:    testLanguage,
		}

		output := genTestOutput(testFn, pkg.FLBuilderImages[testLanguage], testOutDir)

		mockBuilder.On("Setup", mock.Anything, testLanguage, testOutDir).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", ctx).Return(nil).Once()
		mockBuilder.On("BuildSource", ctx, testDir).Return(nil).Once()

		var outbuf bytes.Buffer
		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).DisableAnimation().Build()

		err := cmd.Run(ctx, mockBuilder, bufLogger)

		require.NoError(t, err)
		require.Equal(t, output, (&outbuf).String())
		mockBuilder.AssertExpectations(t)
	})

	t.Run("should return error if builder setup encounters errors", func(t *testing.T) {
		cmd := Build{
			Name:        testFn,
			Source:      testDir,
			Destination: testOutDir,
			Language:    testLanguage,
		}

		mockBuilder.On("Setup", mock.Anything, testLanguage, testOutDir).Return(errors.New("some error")).Once()

		err := cmd.Run(ctx, mockBuilder, testLogger)
		require.Error(t, err)
		mockBuilder.AssertExpectations(t)
	})

	t.Run("should return error if builder image cannot be pulled", func(t *testing.T) {
		cmd := Build{
			Name:        testFn,
			Source:      testDir,
			Destination: testOutDir,
			Language:    testLanguage,
		}

		mockBuilder.On("Setup", mock.Anything, testLanguage, testOutDir).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", ctx).Return(errors.New("some error")).Once()

		err := cmd.Run(ctx, mockBuilder, testLogger)
		require.Error(t, err)
		mockBuilder.AssertExpectations(t)
	})

	t.Run("should return error if builder image encounters errors", func(t *testing.T) {
		cmd := Build{
			Name:        testFn,
			Source:      testDir,
			Language:    testLanguage,
			Destination: testOutDir,
		}

		mockBuilder.On("Setup", mock.Anything, testLanguage, testOutDir).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", ctx).Return(nil).Once()
		mockBuilder.On("BuildSource", ctx, testDir).Return(errors.New("some error")).Once()

		err := cmd.Run(ctx, mockBuilder, testLogger)
		require.Error(t, err)
		mockBuilder.AssertExpectations(t)
	})
}

func genTestOutput(name, image, dest string) string {
	return fmt.Sprintf(`Building %s into a wasm binary...

Setting up...
done
Pulling Javascript builder image (%s) 📦
done
Building source... 🛠️
done

Successfully built function at %s/%s.wasm 🥳🥳
`, name, image, dest, name)
}
