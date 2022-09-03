package admin

import (
	"context"

	"github.com/funlessdev/fl-cli/pkg/log"
)

type join struct{}

func (d *join) Run(ctx context.Context, logger log.FLogger) error {
	logger.Info("admin join stub")
	return nil
}
