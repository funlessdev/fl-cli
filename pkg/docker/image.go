// Copyright 2022 Giuseppe De Palma, Matteo Trentin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type FLImageHandler struct{}

func NewFLImageHandler() ImageHandler {
	return &FLImageHandler{}
}

func (dp *FLImageHandler) Exists(ctx context.Context, dockerClient *client.Client, image string) (bool, error) {
	_, _, err := dockerClient.ImageInspectWithRaw(ctx, image)
	notFound := client.IsErrNotFound(err)

	// notFound being false means the error is something else
	// we still return false as we can't be sure the image actually exists
	if err != nil && !notFound {
		return false, err
	}

	if notFound {
		return false, nil
	}

	return true, nil
}

func (dp *FLImageHandler) Pull(ctx context.Context, dockerClient *client.Client, image string) error {
	out, err := dockerClient.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()

	d := json.NewDecoder(out)

	var event *dockerEvent
	for {
		if err := d.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if event.Error != "" {
			return fmt.Errorf("pulling image: %s", event.Error)
		}
	}
	return nil
}

// struct for decoding docker events, used in PullImage to check if an error occurred during pulling
type dockerEvent struct {
	Status         string `json:"status"`
	Error          string `json:"error"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}
