package fn

import (
	"context"

	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Delete struct {
	Name      string `arg:"" name:"name" help:"name of the function to delete"`
	Namespace string `name:"namespace" short:"n" default:"_" help:"namespace of the function to delete"`
}

func (f *Delete) Run(ctx context.Context, fnHandler client.FnHandler, logger log.FLogger) error {
	res, err := fnHandler.Delete(ctx, f.Name, f.Namespace)
	if err != nil {
		return extractError(err)
	}

	logger.Info(*res.Result)
	return nil
}
