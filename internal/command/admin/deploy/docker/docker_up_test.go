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
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDockerUpRun(t *testing.T) {
	realGetFileInConfigDir := getFileInConfigDir
	getFileInConfigDir = func(string, string) (string, error) {
		return "", errors.New("get compose error")
	}
	defer func() {
		getFileInConfigDir = realGetFileInConfigDir
	}()

	up := Up{}
	ctx := context.TODO()

	mockDockerShell := mocks.NewDockerShell(t)
	_, logger := testLogger()

	t.Run("should return error when setup fails", func(t *testing.T) {
		err := up.Run(ctx, mockDockerShell, logger)
		assert.Error(t, err, "get compose error")
	})

	t.Run("should return error when compose up fails", func(t *testing.T) {
		getFileInConfigDir = func(string, string) (string, error) {
			return "", nil
		}

		mockDockerShell.On("ComposeUp", mock.Anything).Return(errors.New("compose up error")).Once()
		err := up.Run(ctx, mockDockerShell, logger)
		assert.Error(t, err, "compose up error")
	})

	t.Run("should complete successfully when compose up succeeds", func(t *testing.T) {
		mockDockerShell.On("ComposeUp", mock.Anything).Return(nil).Once()
		err := up.Run(ctx, mockDockerShell, logger)
		assert.NoError(t, err)
	})
}

func testLogger() (*bytes.Buffer, log.FLogger) {
	var outbuf bytes.Buffer
	testLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).DisableAnimation().Build()
	return &outbuf, testLogger
}
