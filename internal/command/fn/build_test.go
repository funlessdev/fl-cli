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
	"github.com/stretchr/testify/require"
)

func TestFnBuild(t *testing.T) {

	testFn := "test-fn"
	testLanguage := "js"
	// testSource, _ := filepath.Abs("../../../test/fixtures/test_code.txt")
	testDir, _ := filepath.Abs("../../../test/fixtures/test_dir/")
	testOutDir, _ := filepath.Abs("../../../test/fixtures")
	testCtx := context.Background()
	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).DisableAnimation().Build()

	mockBuilder := mocks.NewDockerBuilder(t)

	t.Run("should use DockerBuilder to build functions", func(t *testing.T) {
		cmd := Build{
			Name:        testFn,
			Source:      testDir,
			Destination: testOutDir,
			Language:    testLanguage,
		}

		mockBuilder := mocks.NewDockerBuilder(t)
		mockBuilder.On("Setup", testCtx, testLanguage, testOutDir).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", testCtx).Return(nil).Once()
		mockBuilder.On("BuildSource", testCtx, testDir).Return(nil).Once()

		err := cmd.Run(testCtx, mockBuilder, testLogger)
		require.NoError(t, err)

		mockBuilder.AssertCalled(t, "BuildSource", testCtx, testDir)
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

		mockBuilder.On("Setup", testCtx, testLanguage, testOutDir).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", testCtx).Return(nil).Once()
		mockBuilder.On("BuildSource", testCtx, testDir).Return(nil).Once()

		var outbuf bytes.Buffer

		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).DisableAnimation().Build()

		err := cmd.Run(testCtx, mockBuilder, bufLogger)

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

		mockBuilder.On("Setup", testCtx, testLanguage, testOutDir).Return(errors.New("some error")).Once()

		err := cmd.Run(testCtx, mockBuilder, testLogger)
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

		mockBuilder.On("Setup", testCtx, testLanguage, testOutDir).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", testCtx).Return(errors.New("some error")).Once()

		err := cmd.Run(testCtx, mockBuilder, testLogger)
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

		mockBuilder.On("Setup", testCtx, testLanguage, testOutDir).Return(nil).Once()
		mockBuilder.On("PullBuilderImage", testCtx).Return(nil).Once()
		mockBuilder.On("BuildSource", testCtx, testDir).Return(errors.New("some error")).Once()

		err := cmd.Run(testCtx, mockBuilder, testLogger)
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
