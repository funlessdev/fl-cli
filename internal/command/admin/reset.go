package admin

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/funlessdev/fl-cli/pkg/admin"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type reset struct{}

func (r *reset) Run(ctx context.Context, logger log.FLogger) error {
	logger.Info("Removing local funless deployment...\n")

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return err
	}

	logger.StartSpinner("Removing Core container... ☠️")
	if err := logger.StopSpinner(admin.RemoveFLContainer(ctx, cli, "fl-core")); err != nil {
		return err
	}

	logger.StartSpinner("Removing Worker container... 🔪")
	if err := logger.StopSpinner(admin.RemoveFLContainer(ctx, cli, "fl-worker")); err != nil {
		return err
	}

	logger.StartSpinner("Removing the function containers... 🔫")
	if err := logger.StopSpinner(admin.RemoveFunctionContainers(ctx, cli)); err != nil {
		return err
	}

	logger.StartSpinner("Removing fl_net network... ✂️")
	if err := logger.StopSpinner(admin.RemoveFLNetwork(ctx, cli, "fl_net")); err != nil {
		return err
	}

	logger.Info("\nAll clear! 👍")

	return err
}
