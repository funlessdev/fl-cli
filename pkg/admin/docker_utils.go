package admin

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/funlessdev/funless-cli/pkg/log"
)

func flNetExists(ctx context.Context, client *client.Client) (bool, types.NetworkResource, error) {
	nets, err := client.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: "fl_net"}),
	})
	if err != nil {
		return false, types.NetworkResource{}, err
	}

	if len(nets) == 0 {
		return false, types.NetworkResource{}, nil
	}

	return true, nets[0], nil
}

func flNetCreate(ctx context.Context, client *client.Client, logger log.FLogger) (string, error) {
	res, err := client.NetworkCreate(ctx, "fl_net", types.NetworkCreate{})
	if err != nil {
		return "", err
	}
	if res.Warning != "" {
		logger.Infof("Warning creating fl_net network: %s", res.Warning)
	}
	return res.ID, nil
}
