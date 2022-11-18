package build

import (
	"context"
	"errors"
	"testing"

	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBuildSource(t *testing.T) {
	wasmBuilder := NewWasmBuilder()
	mockDockerClient := mocks.NewDockerClient(t)
	wasmBuilder.Setup(mockDockerClient, "js", "test")

	ctx := context.TODO()

	t.Run("PullBuilderImage should return error if Pull fails", func(t *testing.T) {
		mockDockerClient.On("Pull", ctx, mock.Anything).Return(errors.New("test error"))
		mockDockerClient.On("ImageExists", ctx, mock.Anything).Return(false, nil).Once()

		err := wasmBuilder.PullBuilderImage(ctx)
		assert.Error(t, err)

		mockDockerClient.AssertNumberOfCalls(t, "Pull", 1)
		mockDockerClient.AssertNumberOfCalls(t, "ImageExists", 1)
		mockDockerClient.AssertExpectations(t)
	})

	t.Run("PullBuilderImage should not call Pull if image already Exists", func(t *testing.T) {
		mockDockerClient.On("ImageExists", ctx, mock.Anything).Return(true, nil).Once()

		err := wasmBuilder.PullBuilderImage(ctx)
		assert.NoError(t, err)

		mockDockerClient.AssertNumberOfCalls(t, "ImageExists", 2)
		mockDockerClient.AssertNumberOfCalls(t, "Pull", 1)
	})

	t.Run("BuildSource should return error if fails", func(t *testing.T) {
		mockDockerClient.On("RunAndWait", ctx, mock.Anything).Return(errors.New("test error"))

		err := wasmBuilder.BuildSource(ctx, "test")
		assert.Error(t, err)

		mockDockerClient.AssertNumberOfCalls(t, "RunAndWait", 1)
		mockDockerClient.AssertExpectations(t)
	})
}
