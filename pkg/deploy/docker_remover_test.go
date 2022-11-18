package deploy

import (
	"context"
	"errors"
	"testing"

	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDockerRemover(t *testing.T) {

	remover := NewDockerRemover("fl-net", "fl-core", "fl-worker", "fl-prom")

	mockDockerClient := mocks.NewDockerClient(t)
	remover.WithDockerClient(mockDockerClient)

	ctx := context.TODO()

	t.Run("RemoveFLNetwork should not call Remove if network does not exist", func(t *testing.T) {
		mockDockerClient.On("NetworkExists", ctx, mock.Anything).Return(false, "", nil).Once()

		err := remover.RemoveFLNetwork(ctx)
		assert.NoError(t, err)
		mockDockerClient.AssertNumberOfCalls(t, "NetworkExists", 1)
		mockDockerClient.AssertNumberOfCalls(t, "RemoveNetwork", 0)
	})

	t.Run("RemoveFLNetwork should return error when fails", func(t *testing.T) {
		mockDockerClient.On("NetworkExists", ctx, mock.Anything).Return(true, "", nil)
		mockDockerClient.On("RemoveNetwork", ctx, mock.Anything).Return(errors.New("error"))

		err := remover.RemoveFLNetwork(ctx)
		assert.Error(t, err)
		mockDockerClient.AssertNumberOfCalls(t, "RemoveNetwork", 1)
	})

	t.Run("RemoveXXXContainer should not call Remove if container does not exist", func(t *testing.T) {
		mockDockerClient.On("CtrExists", ctx, mock.Anything).Return(false, "", nil).Times(3)

		err := remover.RemoveCoreContainer(ctx)
		assert.NoError(t, err)

		err = remover.RemoveWorkerContainer(ctx)
		assert.NoError(t, err)

		err = remover.RemovePromContainer(ctx)
		assert.NoError(t, err)
		mockDockerClient.AssertNumberOfCalls(t, "CtrExists", 3)
		mockDockerClient.AssertNumberOfCalls(t, "RemoveCtr", 0)
	})

	t.Run("RemoveXXXContainer should return error when fails", func(t *testing.T) {
		mockDockerClient.On("CtrExists", ctx, mock.Anything).Return(true, "", nil).Times(3)
		mockDockerClient.On("RemoveCtr", ctx, mock.Anything).Return(errors.New("error")).Times(3)

		err := remover.RemoveCoreContainer(ctx)
		assert.Error(t, err)

		err = remover.RemoveWorkerContainer(ctx)
		assert.Error(t, err)

		err = remover.RemovePromContainer(ctx)
		assert.Error(t, err)
		mockDockerClient.AssertNumberOfCalls(t, "RemoveCtr", 3)
	})
}
