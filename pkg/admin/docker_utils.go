// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
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
