package admin

import (
	"context"

	"github.com/funlessdev/fl-cli/pkg/log"
)

type coreInit struct{} // init is reserved in Go

func (d *coreInit) Run(ctx context.Context, logger log.FLogger) error {
	logger.Info("admin init stub")
	return nil
}
