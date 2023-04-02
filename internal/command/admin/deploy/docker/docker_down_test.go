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

package admin_deploy_docker

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/funlessdev/fl-cli/pkg/homedir"
	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDockerDownRun(t *testing.T) {
	homedirPath, err := os.MkdirTemp("", "funless-test-homedir-")
	require.NoError(t, err)

	homedir.GetHomeDir = func() (string, error) {
		return homedirPath, nil
	}
	defer func() {
		homedir.GetHomeDir = os.UserHomeDir
		os.RemoveAll(homedirPath)
	}()

	down := Down{}
	ctx := context.TODO()

	mockDockerShell := mocks.NewDockerShell(t)
	out, logger := testLogger()

	t.Run("should return error when setup fails", func(t *testing.T) {
		homedir.GetHomeDir = func() (string, error) {
			return "", errors.New("some home error")
		}
		err := down.Run(ctx, mockDockerShell, logger)
		require.Error(t, err)

		homedir.GetHomeDir = func() (string, error) {
			return homedirPath, nil
		}
	})

	t.Run("should return error when compose down fails", func(t *testing.T) {
		path, err := downloadFile("docker-compose.yml", dockerComposeYmlUrl)
		require.NoError(t, err)

		out.Reset()
		mockDockerShell.On("ComposeDown", mock.Anything, path).Return(errors.New("some compose down error")).Once()

		err = down.Run(ctx, mockDockerShell, logger)
		require.Error(t, err)
	})

	t.Run("should return error when unable to find docker-compose.yml but local deployment is found", func(t *testing.T) {
		homedir.GetHomeDir = func() (string, error) {
			return "", os.ErrNotExist
		}

		defer func() {
			homedir.GetHomeDir = func() (string, error) {
				return homedirPath, nil
			}
		}()

		out.Reset()
		mockDockerShell.On("ComposeList", mock.Anything).Return([]string{"fl"}, nil).Once()

		err = down.Run(ctx, mockDockerShell, logger)
		require.Error(t, err)
		require.Equal(t, err.Error(), "unable to locate docker-compose.yml, but a local deployment was found. The file might have been moved or deleted.")

	})

	t.Run("should return error when unable to find docker-compose.yml and local deployment is not found", func(t *testing.T) {
		homedir.GetHomeDir = func() (string, error) {
			return "", os.ErrNotExist
		}

		defer func() {
			homedir.GetHomeDir = func() (string, error) {
				return homedirPath, nil
			}
		}()

		out.Reset()
		mockDockerShell.On("ComposeList", mock.Anything).Return([]string{"test"}, nil).Once()

		err = down.Run(ctx, mockDockerShell, logger)
		require.Error(t, err)
		require.Equal(t, err.Error(), "no local deployment found, nothing to remove. Use \"fl admin deploy docker up\" to create one.")

	})

	t.Run("should remove docker-compose.yml when succeds", func(t *testing.T) {
		path, err := downloadFile("docker-compose.yml", dockerComposeYmlUrl)
		require.NoError(t, err)

		out.Reset()
		mockDockerShell.On("ComposeDown", mock.Anything, path).Return(nil)

		err = down.Run(ctx, mockDockerShell, logger)
		require.NoError(t, err)
		require.Contains(t, out.String(), "\nAll clear!")

		require.NoFileExists(t, path)
	})
}
