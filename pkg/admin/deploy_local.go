package admin

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/funlessdev/funless-cli/pkg/log"
)

func ObtainFLNet(ctx context.Context, client *client.Client, logger log.FLogger) (string, error) {
	exists, net, err := flNetExists(ctx, client)

	if err != nil {
		return "", err
	}
	if exists {
		return net.ID, nil
	}
	return flNetCreate(ctx, client, logger)
}
