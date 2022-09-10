package admin

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/funlessdev/fl-cli/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRun_WhenSuccessful(t *testing.T) {
	dev := dev{}
	ctx := context.TODO()

	var outbuf bytes.Buffer
	testLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).DisableAnimation().Build()

	deployer := mocks.NewDockerDeployer(t)

	t.Run("print error when setup client fails", func(t *testing.T) {
		deployer.On("SetupClient", ctx).Return(func(ctx context.Context) error {
			return errors.New("error")
		}).Once()

		dev.Run(ctx, deployer, testLogger)

		expectedOutput := []string{
			"Deploying funless locally...\n",
			"\n",
			"Setting things up...\n",
			"failed\n",
			"",
		}

		assertOutput(t, expectedOutput, &outbuf)
	})

	t.Run("print error when docker networks setup fails", func(t *testing.T) {
		deployer.On("SetupClient", ctx).Return(func(ctx context.Context) error {
			return nil
		})
		deployer.On("SetupFLNetworks", ctx).Return(func(ctx context.Context) error {
			return errors.New("error")
		}).Once()

		dev.Run(ctx, deployer, testLogger)

		expectedOutput := []string{
			"Deploying funless locally...\n",
			"\n",
			"Setting things up...\n",
			"failed\n",
			"",
		}

		assertOutput(t, expectedOutput, &outbuf)
	})

	t.Run("print error when pulling core image fails", func(t *testing.T) {
		deployer.On("SetupFLNetworks", ctx).Return(func(ctx context.Context) error {
			return nil
		})
		deployer.On("PullCoreImage", ctx).Return(func(ctx context.Context) error {
			return errors.New("error")
		}).Once()

		dev.Run(ctx, deployer, testLogger)

		expectedOutput := []string{
			"Deploying funless locally...\n",
			"\n",
			"Setting things up...\n",
			"done\n",
			"pulling Core image (ghcr.io/funlessdev/fl-core:latest) 📦\n",
			"failed\n",
			"",
		}

		assertOutput(t, expectedOutput, &outbuf)
	})

	t.Run("print error when pulling worker image fails", func(t *testing.T) {
		deployer.On("PullCoreImage", ctx).Return(func(ctx context.Context) error {
			return nil
		})
		deployer.On("PullWorkerImage", ctx).Return(func(ctx context.Context) error {
			return errors.New("error")
		}).Once()

		dev.Run(ctx, deployer, testLogger)

		expectedOutput := []string{
			"Deploying funless locally...\n",
			"\n",
			"Setting things up...\n",
			"done\n",
			"pulling Core image (ghcr.io/funlessdev/fl-core:latest) 📦\n",
			"done\n",
			"pulling Worker image (ghcr.io/funlessdev/fl-worker:latest) 🗃\n",
			"failed\n",
			"",
		}

		assertOutput(t, expectedOutput, &outbuf)
	})

	t.Run("print error when starting core fails", func(t *testing.T) {
		deployer.On("PullWorkerImage", ctx).Return(func(ctx context.Context) error {
			return nil
		})
		deployer.On("StartCore", ctx).Return(func(ctx context.Context) error {
			return errors.New("error")
		}).Once()

		dev.Run(ctx, deployer, testLogger)

		expectedOutput := []string{
			"Deploying funless locally...\n",
			"\n",
			"Setting things up...\n",
			"done\n",
			"pulling Core image (ghcr.io/funlessdev/fl-core:latest) 📦\n",
			"done\n",
			"pulling Worker image (ghcr.io/funlessdev/fl-worker:latest) 🗃\n",
			"done\n",
			"starting Core container 🎛️\n",
			"failed\n",
			"",
		}

		assertOutput(t, expectedOutput, &outbuf)
	})

	t.Run("print error when starting worker fails", func(t *testing.T) {
		deployer.On("StartCore", ctx).Return(func(ctx context.Context) error {
			return nil
		})
		deployer.On("StartWorker", ctx).Return(func(ctx context.Context) error {
			return errors.New("error")
		}).Once()

		dev.Run(ctx, deployer, testLogger)

		expectedOutput := []string{
			"Deploying funless locally...\n",
			"\n",
			"Setting things up...\n",
			"done\n",
			"pulling Core image (ghcr.io/funlessdev/fl-core:latest) 📦\n",
			"done\n",
			"pulling Worker image (ghcr.io/funlessdev/fl-worker:latest) 🗃\n",
			"done\n",
			"starting Core container 🎛️\n",
			"done\n",
			"starting Worker container 👷\n",
			"failed\n",
			"",
		}

		assertOutput(t, expectedOutput, &outbuf)
	})

	t.Run("success prints when everything goes well", func(t *testing.T) {
		deployer.On("StartWorker", ctx).Return(func(ctx context.Context) error {
			return nil
		})

		dev.Run(ctx, deployer, testLogger)

		expectedOutput := []string{
			"Deploying funless locally...\n",
			"\n",
			"Setting things up...\n",
			"done\n",
			"pulling Core image (ghcr.io/funlessdev/fl-core:latest) 📦\n",
			"done\n",
			"pulling Worker image (ghcr.io/funlessdev/fl-worker:latest) 🗃\n",
			"done\n",
			"starting Core container 🎛️\n",
			"done\n",
			"starting Worker container 👷\n",
			"done\n",
			"\n",
			"Deployment complete!\n",
			"You can now start using Funless! 🎉\n",
			"",
		}

		assertOutput(t, expectedOutput, &outbuf)
	})

}

func assertOutput(t *testing.T, expected []string, outbuf *bytes.Buffer) {
	t.Helper()
	for _, expected := range expected {
		line, _ := outbuf.ReadString('\n')
		assert.Equal(t, expected, line)
	}
}
