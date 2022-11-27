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
	"context"
	"errors"
	"testing"

	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDockerDownRun(t *testing.T) {
	down := docker_down{}
	ctx := context.TODO()

	mockRemover := mocks.NewDockerRemover(t)
	_, logger := testLogger()

	t.Run("should return error when removing Core fails", func(t *testing.T) {
		mockRemover.On("WithDockerClient", mock.Anything).Return()
		mockRemover.On("RemoveCoreContainer", mock.Anything).Return(errors.New("error")).Once()

		err := down.Run(ctx, mockRemover, logger)
		require.Error(t, err)
		mockRemover.AssertNumberOfCalls(t, "RemoveCoreContainer", 1)
	})

	t.Run("should return error when removing Worker fails", func(t *testing.T) {
		mockRemover.On("RemoveCoreContainer", mock.Anything).Return(nil)
		mockRemover.On("RemoveWorkerContainer", mock.Anything).Return(errors.New("error")).Once()

		err := down.Run(ctx, mockRemover, logger)
		require.Error(t, err)
		mockRemover.AssertNumberOfCalls(t, "RemoveWorkerContainer", 1)
	})

	t.Run("should return error when removing Prometheus fails", func(t *testing.T) {
		mockRemover.On("RemoveWorkerContainer", mock.Anything).Return(nil)
		mockRemover.On("RemovePromContainer", mock.Anything).Return(errors.New("error")).Once()

		err := down.Run(ctx, mockRemover, logger)
		require.Error(t, err)
		mockRemover.AssertNumberOfCalls(t, "RemovePromContainer", 1)
	})

	t.Run("should return error when removing FL network fails", func(t *testing.T) {
		mockRemover.On("RemovePromContainer", mock.Anything).Return(nil)
		mockRemover.On("RemoveFLNetwork", mock.Anything).Return(errors.New("error")).Once()

		err := down.Run(ctx, mockRemover, logger)
		require.Error(t, err)
		mockRemover.AssertNumberOfCalls(t, "RemoveFLNetwork", 1)
	})

	t.Run("successful prints when everything goes well", func(t *testing.T) {
		mockRemover.On("RemoveFLNetwork", mock.Anything).Return(nil)

		outbuf, testLogger := testLogger()
		err := down.Run(ctx, mockRemover, testLogger)

		expectedOutput := `Removing local FunLess deployment...

Removing Core container... ☠️
done
Removing Worker container... 🔪
done
Removing Prometheus container... ⚰️
done
Removing fl network... ✂️
done

All clear! 👍
`
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, outbuf.String())
	})

}
