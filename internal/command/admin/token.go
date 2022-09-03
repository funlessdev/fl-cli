package admin

import (
	"context"

	"github.com/funlessdev/fl-cli/pkg/log"
)

type token struct{}

func (d *token) Run(ctx context.Context, logger log.FLogger) error {
	logger.Info("admin token stub")
	return nil
}
