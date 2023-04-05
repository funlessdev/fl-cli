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

	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestKubernetesUpRun(t *testing.T) {
	k8s := Up{}
	ctx := context.TODO()

	mockDeployer := mocks.NewKubernetesDeployer(t)
	_, logger := testLogger()

	t.Run("should return error when setting up Deployer fails", func(t *testing.T) {
		mockDeployer.On("WithConfig", mock.Anything).Return(errors.New("error")).Once()

		err := k8s.Run(ctx, mockDeployer, logger, client.Config{})
		require.Error(t, err)
		mockDeployer.AssertNumberOfCalls(t, "WithConfig", 1)
	})

	t.Run("should return error when creating Namespace fails", func(t *testing.T) {
		mockDeployer.On("WithConfig", mock.Anything).Return(nil)
		mockDeployer.On("CreateNamespace", mock.Anything).Return(errors.New("error")).Once()

		err := k8s.Run(ctx, mockDeployer, logger, client.Config{})
		require.Error(t, err)
		mockDeployer.AssertNumberOfCalls(t, "CreateNamespace", 1)
	})

	t.Run("should return error when creating ServiceAccount fails", func(t *testing.T) {
		mockDeployer.On("CreateNamespace", mock.Anything).Return(nil)
		mockDeployer.On("CreateSvcAccount", mock.Anything).Return(errors.New("error")).Once()

		err := k8s.Run(ctx, mockDeployer, logger, client.Config{})
		require.Error(t, err)
		mockDeployer.AssertNumberOfCalls(t, "CreateSvcAccount", 1)
	})

	t.Run("should return error when creating Role fails", func(t *testing.T) {
		mockDeployer.On("CreateSvcAccount", mock.Anything).Return(nil)
		mockDeployer.On("CreateRole", mock.Anything).Return(errors.New("error")).Once()

		err := k8s.Run(ctx, mockDeployer, logger, client.Config{})
		require.Error(t, err)
		mockDeployer.AssertNumberOfCalls(t, "CreateRole", 1)
	})

	t.Run("should return error when creating RoleBinding fails", func(t *testing.T) {
		mockDeployer.On("CreateRole", mock.Anything).Return(nil)
		mockDeployer.On("CreateRoleBinding", mock.Anything).Return(errors.New("error")).Once()

		err := k8s.Run(ctx, mockDeployer, logger, client.Config{})
		require.Error(t, err)
		mockDeployer.AssertNumberOfCalls(t, "CreateRoleBinding", 1)
	})

	t.Run("should return error when creating Prometheus ConfigMap fails", func(t *testing.T) {
		mockDeployer.On("CreateRoleBinding", mock.Anything).Return(nil)
		mockDeployer.On("CreatePrometheusConfigMap", mock.Anything).Return(errors.New("error")).Once()

		err := k8s.Run(ctx, mockDeployer, logger, client.Config{})
		require.Error(t, err)
		mockDeployer.AssertNumberOfCalls(t, "CreatePrometheusConfigMap", 1)
	})

	t.Run("should return error when deploying Prometheus fails", func(t *testing.T) {
		mockDeployer.On("CreatePrometheusConfigMap", mock.Anything).Return(nil)
		mockDeployer.On("DeployPrometheus", mock.Anything).Return(errors.New("error")).Once()

		err := k8s.Run(ctx, mockDeployer, logger, client.Config{})
		require.Error(t, err)
		mockDeployer.AssertNumberOfCalls(t, "DeployPrometheus", 1)
	})

	t.Run("should return error when deploying Prometheus Service fails", func(t *testing.T) {
		mockDeployer.On("DeployPrometheus", mock.Anything).Return(nil)
		mockDeployer.On("DeployPrometheusService", mock.Anything).Return(errors.New("error")).Once()

		err := k8s.Run(ctx, mockDeployer, logger, client.Config{})
		require.Error(t, err)
		mockDeployer.AssertNumberOfCalls(t, "DeployPrometheusService", 1)
	})

	t.Run("should return error when deploying Postgres fails", func(t *testing.T) {
		mockDeployer.On("DeployPrometheusService", mock.Anything).Return(nil)
		mockDeployer.On("DeployPostgres", mock.Anything).Return(errors.New("error")).Once()

		err := k8s.Run(ctx, mockDeployer, logger, client.Config{})
		require.Error(t, err)
		mockDeployer.AssertNumberOfCalls(t, "DeployPostgres", 1)
	})

	t.Run("should return error when deploying Postgres Service fails", func(t *testing.T) {
		mockDeployer.On("DeployPostgres", mock.Anything).Return(nil)
		mockDeployer.On("DeployPostgresService", mock.Anything).Return(errors.New("error")).Once()

		err := k8s.Run(ctx, mockDeployer, logger, client.Config{})
		require.Error(t, err)
		mockDeployer.AssertNumberOfCalls(t, "DeployPostgresService", 1)
	})

	t.Run("should return error when starting init-postgres fails", func(t *testing.T) {
		mockDeployer.On("DeployPostgresService", mock.Anything).Return(nil)
		mockDeployer.On("StartInitPostgres", mock.Anything).Return(errors.New("error")).Once()

		err := k8s.Run(ctx, mockDeployer, logger, client.Config{})
		require.Error(t, err)
		mockDeployer.AssertNumberOfCalls(t, "StartInitPostgres", 1)
	})

	t.Run("should return error when create Core Secrets fails", func(t *testing.T) {
		mockDeployer.On("StartInitPostgres", mock.Anything).Return(nil)
		mockDeployer.On("CreateCoreSecrets", mock.Anything).Return(errors.New("error")).Once()

		err := k8s.Run(ctx, mockDeployer, logger, client.Config{})
		require.Error(t, err)
		mockDeployer.AssertNumberOfCalls(t, "CreateCoreSecrets", 1)
	})

	t.Run("should return error when deploying Core fails", func(t *testing.T) {
		mockDeployer.On("CreateCoreSecrets", mock.Anything).Return(nil)
		mockDeployer.On("DeployCore", mock.Anything).Return(errors.New("error")).Once()

		err := k8s.Run(ctx, mockDeployer, logger, client.Config{})
		require.Error(t, err)
		mockDeployer.AssertNumberOfCalls(t, "DeployCore", 1)
	})

	t.Run("should return error when deploying Core Service fails", func(t *testing.T) {
		mockDeployer.On("DeployCore", mock.Anything).Return(nil)
		mockDeployer.On("DeployCoreService", mock.Anything).Return(errors.New("error")).Once()

		err := k8s.Run(ctx, mockDeployer, logger, client.Config{})
		require.Error(t, err)
		mockDeployer.AssertNumberOfCalls(t, "DeployCoreService", 1)
	})

	t.Run("should return error when deploying Workers fails", func(t *testing.T) {
		mockDeployer.On("DeployCoreService", mock.Anything).Return(nil)
		mockDeployer.On("DeployWorker", mock.Anything).Return(errors.New("error")).Once()

		err := k8s.Run(ctx, mockDeployer, logger, client.Config{})
		require.Error(t, err)
		mockDeployer.AssertNumberOfCalls(t, "DeployWorker", 1)
	})

	t.Run("should show error and then complete deployment when extracting Tokens fails", func(t *testing.T) {
		mockDeployer.On("DeployWorker", mock.Anything).Return(nil)
		mockDeployer.On("ExtractTokens", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("token error")).Once()
		outbuf, testLogger := testLogger()
		err := k8s.Run(ctx, mockDeployer, testLogger, client.Config{})

		expectedOutput := `Deploying FunLess on Kubernetes...

Setting things up...
done
Creating Namespace...
done
Creating ServiceAccount...
done
Creating Role...
done
Creating RoleBinding...
done
Creating Prometheus ConfigMap...
done
Deploying Prometheus...
done
Deploying Prometheus Service...
done
Deploying PostgreSQL...
done
Deploying PostgreSQL Service...
done
Starting init-postgres Job...
done
Creating Core Secrets...
done
Deploying Core...
done
Deploying Core Service...
done
Deploying Workers...
done
Extracting auth tokens...
failed

token error


Couldn't extract auth tokens from core pod. Completing deployment...
Deployment complete!
You can now start using FunLess! 🎉
`
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, outbuf.String())
		mockDeployer.AssertNumberOfCalls(t, "ExtractTokens", 1)
	})

	t.Run("successful prints when everything goes well", func(t *testing.T) {
		mockDeployer.On("ExtractTokens", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		outbuf, testLogger := testLogger()
		err := k8s.Run(ctx, mockDeployer, testLogger, client.Config{})

		expectedOutput := `Deploying FunLess on Kubernetes...

Setting things up...
done
Creating Namespace...
done
Creating ServiceAccount...
done
Creating Role...
done
Creating RoleBinding...
done
Creating Prometheus ConfigMap...
done
Deploying Prometheus...
done
Deploying Prometheus Service...
done
Deploying PostgreSQL...
done
Deploying PostgreSQL Service...
done
Starting init-postgres Job...
done
Creating Core Secrets...
done
Deploying Core...
done
Deploying Core Service...
done
Deploying Workers...
done
Extracting auth tokens...
done



Remember to add these tokens in ~/.fl/config as api_token and admin_token.

Deployment complete!
You can now start using FunLess! 🎉
`
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, outbuf.String())
	})
}
