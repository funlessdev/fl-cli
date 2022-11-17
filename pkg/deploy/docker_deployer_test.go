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

	mockImgHandler := mocks.NewImageHandler(t)
	mockCtrHandler := mocks.NewContainerHandler(t)
	mockNtwHandler := mocks.NewNetworkHandler(t)

	deployer := NewDockerDeployer(mockImgHandler, mockCtrHandler, mockNtwHandler)

	ctx := context.TODO()

	t.Run("PullXXXImage should return error if Pull fails", func(t *testing.T) {
		mockImgHandler.On("Pull", ctx, mock.AnythingOfType("*client.Client"), mock.Anything).Return(errors.New("test error"))
		mockImgHandler.On("Exists", ctx, mock.AnythingOfType("*client.Client"), mock.Anything).Return(false, nil).Times(3)

		err := deployer.PullCoreImage(ctx)
		assert.Error(t, err)

		err = deployer.PullWorkerImage(ctx)
		assert.Error(t, err)

		err = deployer.PullCoreImage(ctx)
		assert.Error(t, err)

		mockImgHandler.AssertNumberOfCalls(t, "Pull", 3)
		mockImgHandler.AssertNumberOfCalls(t, "Exists", 3)
		mockImgHandler.AssertExpectations(t)
	})

	t.Run("PullXXXImage should not call Pull if image already Exists", func(t *testing.T) {
		mockImgHandler.On("Exists", ctx, mock.AnythingOfType("*client.Client"), mock.Anything).Return(true, nil).Times(3)

		err := deployer.PullCoreImage(ctx)
		assert.NoError(t, err)

		err = deployer.PullWorkerImage(ctx)
		assert.NoError(t, err)

		err = deployer.PullPromImage(ctx)
		assert.NoError(t, err)

		mockImgHandler.AssertNumberOfCalls(t, "Exists", 6)
		mockImgHandler.AssertNumberOfCalls(t, "Pull", 3)
	})

	t.Run("StartXXX methods should return error if RunAsync fails", func(t *testing.T) {
		mockCtrHandler.On("RunAsync", ctx, mock.AnythingOfType("*client.Client"), mock.Anything).Return(errors.New("test error"))

		err := deployer.StartCore(ctx)
		assert.Error(t, err)

		err = deployer.StartWorker(ctx)
		assert.Error(t, err)

		err = deployer.StartProm(ctx)
		assert.Error(t, err)

		mockCtrHandler.AssertNumberOfCalls(t, "RunAsync", 3)
		mockCtrHandler.AssertExpectations(t)
	})

	t.Run("CreateFlNetwork should return error if Create fails", func(t *testing.T) {
		mockNtwHandler.On("Create", ctx, mock.AnythingOfType("*client.Client"), mock.Anything).Return("", errors.New("test error"))
		mockNtwHandler.On("Exists", ctx, mock.AnythingOfType("*client.Client"), mock.Anything).Return(false, "", nil).Once()

		err := deployer.CreateFLNetwork(ctx)
		assert.Error(t, err)

		mockNtwHandler.AssertNumberOfCalls(t, "Exists", 1)
		mockNtwHandler.AssertNumberOfCalls(t, "Create", 1)
		mockNtwHandler.AssertExpectations(t)
	})

	t.Run("CreateFlNetwork should not call Create if network already Exists", func(t *testing.T) {
		mockNtwHandler.On("Exists", ctx, mock.AnythingOfType("*client.Client"), mock.Anything).Return(true, "id", nil)

		err := deployer.CreateFLNetwork(ctx)
		assert.NoError(t, err)

		mockNtwHandler.AssertNumberOfCalls(t, "Exists", 2)
		mockNtwHandler.AssertNumberOfCalls(t, "Create", 1)
	})
}
