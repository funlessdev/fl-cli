package admin

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/admin"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type dev struct{}

func (d *dev) Run(ctx context.Context, logger log.FLogger) error {
	logger.Info("Deploying funless locally...\n")

	logger.StartSpinner("Setting things up...")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return err
	}
	deployer := admin.NewLocalDeployer(ctx, cli, "fl_net")

	if err := logger.StopSpinner(deployer.Apply(admin.SetupFLNetwork)); err != nil {
		return err
	}

	logger.StartSpinner(fmt.Sprintf("pulling Core image (%s) 📦", pkg.FLCore))
	if err := logger.StopSpinner(deployer.Apply(admin.PullCoreImage)); err != nil {
		return err
	}

	logger.StartSpinner(fmt.Sprintf("pulling Worker image (%s) 🗃", pkg.FLWorker))
	if err := logger.StopSpinner(deployer.Apply(admin.PullWorkerImage)); err != nil {
		return err
	}

	logger.StartSpinner("starting Core container 🎛️")
	if err := logger.StopSpinner(deployer.Apply(admin.StartCoreContainer)); err != nil {
		return err
	}

	logger.StartSpinner("starting Worker container 👷")
	if err := logger.StopSpinner(deployer.Apply(admin.StartWorkerContainer)); err != nil {
		return err
	}

	logger.Info("\nDeployment complete!")
	logger.Info("You can now start using Funless! 🎉")

	return err
}
