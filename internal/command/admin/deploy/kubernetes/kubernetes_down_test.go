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

package admin_deploy_kubernetes

import (
	"context"
	"errors"
	"testing"

	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestKubernetesDownRun(t *testing.T) {
	k8sRm := Down{}
	ctx := context.TODO()

	mockRemover := mocks.NewKubernetesRemover(t)
	_, logger := testLogger()

	t.Run("should return error when setting up Remover fails", func(t *testing.T) {
		mockRemover.On("WithConfig", mock.Anything).Return(errors.New("error")).Once()

		err := k8sRm.Run(ctx, mockRemover, logger)
		require.Error(t, err)
		mockRemover.AssertNumberOfCalls(t, "WithConfig", 1)
	})

	t.Run("should return error when removing Namespace fails", func(t *testing.T) {
		mockRemover.On("WithConfig", mock.Anything).Return(nil)
		mockRemover.On("RemoveNamespace", mock.Anything).Return(errors.New("error")).Once()

		err := k8sRm.Run(ctx, mockRemover, logger)
		require.Error(t, err)
		mockRemover.AssertNumberOfCalls(t, "RemoveNamespace", 1)
	})

	t.Run("successful prints when everything goes well", func(t *testing.T) {
		mockRemover.On("RemoveNamespace", mock.Anything).Return(nil)

		outbuf, testLogger := testLogger()
		err := k8sRm.Run(ctx, mockRemover, testLogger)

		expectedOutput := `Removing Kubernetes FunLess deployment...

Setting things up...
done
Removing Namespace...
done

All clear!
`
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, outbuf.String())
	})

}
