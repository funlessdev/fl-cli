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

package admin

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	dev := dev{}
	ctx := context.TODO()

	deployer := mocks.NewDockerDeployer(t)
	_, logger := testLogger()

	t.Run("should return error when setup client fails", func(t *testing.T) {
		deployer.On("WithImages", mock.Anything, mock.Anything).Return()
		deployer.On("WithDockerClient", mock.Anything, mock.Anything).Return()
		deployer.On("WithLogs", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
		err := dev.Run(ctx, deployer, logger)
		require.Error(t, err)
		deployer.AssertExpectations(t)
	})

	t.Run("should return error when docker network create fails", func(t *testing.T) {
		deployer.On("WithLogs", mock.Anything, mock.Anything).Return(nil)
		deployer.On("CreateFLNetwork", ctx).Return(errors.New("error")).Once()

		err := dev.Run(ctx, deployer, logger)
		require.Error(t, err)
	})

	t.Run("should return error when pulling core image fails", func(t *testing.T) {
		deployer.On("CreateFLNetwork", ctx).Return(nil)
		deployer.On("PullCoreImage", ctx).Return(errors.New("error")).Once()

		err := dev.Run(ctx, deployer, logger)
		require.Error(t, err)
	})

	t.Run("should return error when pulling worker image fails", func(t *testing.T) {
		deployer.On("PullCoreImage", ctx).Return(nil)
		deployer.On("PullWorkerImage", ctx).Return(errors.New("error")).Once()

		err := dev.Run(ctx, deployer, logger)
		require.Error(t, err)
	})

	t.Run("should return error when pulling prometheus image fails", func(t *testing.T) {
		deployer.On("PullWorkerImage", ctx).Return(nil)
		deployer.On("PullPromImage", ctx).Return(errors.New("error")).Once()

		err := dev.Run(ctx, deployer, logger)
		require.Error(t, err)
	})

	t.Run("should return error when starting core fails", func(t *testing.T) {
		deployer.On("PullPromImage", ctx).Return(nil)
		deployer.On("StartCore", ctx).Return(errors.New("error")).Once()

		err := dev.Run(ctx, deployer, logger)
		require.Error(t, err)
	})

	t.Run("should return error when starting worker fails", func(t *testing.T) {
		deployer.On("StartCore", ctx).Return(nil)
		deployer.On("StartWorker", ctx).Return(errors.New("error")).Once()

		err := dev.Run(ctx, deployer, logger)
		require.Error(t, err)
	})

	t.Run("should return error when starting prometheus fails", func(t *testing.T) {
		deployer.On("StartWorker", ctx).Return(nil)
		deployer.On("StartProm", ctx).Return(errors.New("error")).Once()

		err := dev.Run(ctx, deployer, logger)
		require.Error(t, err)
	})
	t.Run("successful prints when everything goes well", func(t *testing.T) {
		deployer.On("StartProm", ctx).Return(func(ctx context.Context) error {
			return nil
		})

		outbuf, testLogger := testLogger()

		err := dev.Run(ctx, deployer, testLogger)
		require.NoError(t, err)

		expectedOutput := `Deploying FunLess locally...

Setting things up...
done
pulling Core image () 🐋
done
pulling Worker image () 🐋
done
pulling Prometheus image 🐋
done
starting Core container 🎛️
done
starting Worker container 👷
done
starting Prometheus container 📊
done

Deployment complete!
You can now start using FunLess! 🎉
`
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, outbuf.String())

	})
}

func testLogger() (*bytes.Buffer, log.FLogger) {
	var outbuf bytes.Buffer
	testLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).DisableAnimation().Build()
	return &outbuf, testLogger
}
