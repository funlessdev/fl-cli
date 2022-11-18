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

package deploy

import (
	"context"
	"errors"
	"testing"

	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDockerDeploy(t *testing.T) {

	mockDockerClient := mocks.NewDockerClient(t)

	deployer := NewDockerDeployer("test-net", "test-core", "test-worker", "test-prom")
	deployer.WithDockerClient(mockDockerClient)

	ctx := context.TODO()

	t.Run("PullXXXImage should return error if Pull fails", func(t *testing.T) {
		mockDockerClient.On("Pull", ctx, mock.Anything).Return(errors.New("test error"))
		mockDockerClient.On("ImageExists", ctx, mock.Anything).Return(false, nil).Times(3)

		err := deployer.PullCoreImage(ctx)
		assert.Error(t, err)

		err = deployer.PullWorkerImage(ctx)
		assert.Error(t, err)

		err = deployer.PullPromImage(ctx)
		assert.Error(t, err)

		mockDockerClient.AssertNumberOfCalls(t, "Pull", 3)
		mockDockerClient.AssertNumberOfCalls(t, "ImageExists", 3)
		mockDockerClient.AssertExpectations(t)
	})

	t.Run("PullXXXImage should not call Pull if image already Exists", func(t *testing.T) {
		mockDockerClient.On("ImageExists", ctx, mock.Anything).Return(true, nil).Times(3)

		err := deployer.PullCoreImage(ctx)
		assert.NoError(t, err)

		err = deployer.PullWorkerImage(ctx)
		assert.NoError(t, err)

		err = deployer.PullPromImage(ctx)
		assert.NoError(t, err)

		mockDockerClient.AssertNumberOfCalls(t, "ImageExists", 6)
		mockDockerClient.AssertNumberOfCalls(t, "Pull", 3)
	})

	t.Run("StartXXX methods should return error if RunAsync fails", func(t *testing.T) {
		mockDockerClient.On("RunAsync", ctx, mock.Anything).Return(errors.New("test error"))

		err := deployer.StartCore(ctx)
		assert.Error(t, err)

		err = deployer.StartWorker(ctx)
		assert.Error(t, err)

		err = deployer.StartProm(ctx)
		assert.Error(t, err)

		mockDockerClient.AssertNumberOfCalls(t, "RunAsync", 3)
		mockDockerClient.AssertExpectations(t)
	})

	t.Run("CreateFlNetwork should return error if Create fails", func(t *testing.T) {
		mockDockerClient.On("CreateNetwork", ctx, mock.Anything).Return("", errors.New("test error"))
		mockDockerClient.On("NetworkExists", ctx, mock.Anything).Return(false, "", nil).Once()

		err := deployer.CreateFLNetwork(ctx)
		assert.Error(t, err)

		mockDockerClient.AssertNumberOfCalls(t, "NetworkExists", 1)
		mockDockerClient.AssertNumberOfCalls(t, "CreateNetwork", 1)
		mockDockerClient.AssertExpectations(t)
	})

	t.Run("CreateFlNetwork should not call Create if network already Exists", func(t *testing.T) {
		mockDockerClient.On("NetworkExists", ctx, mock.Anything).Return(true, "id", nil)

		err := deployer.CreateFLNetwork(ctx)
		assert.NoError(t, err)

		mockDockerClient.AssertNumberOfCalls(t, "NetworkExists", 2)
		mockDockerClient.AssertNumberOfCalls(t, "CreateNetwork", 1)
	})
}
