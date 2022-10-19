package build

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func imageExistsLocally(ctx context.Context, c *client.Client, image string) (bool, error) {
	_, _, err := c.ImageInspectWithRaw(ctx, image)
	notFound := client.IsErrNotFound(err)

	/* notFound being false means the error is something else; we still return false as we can't be sure the image actually exists */
	if err != nil && !notFound {
		return false, err
	}

	if notFound {
		return false, nil
	}

	return true, nil
}

func pullImage(ctx context.Context, c *client.Client, image string) error {
	exists, err := imageExistsLocally(ctx, c, image)
	if exists {
		return nil
	}
	if err != nil {
		return err
	}

	out, err := c.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()

	d := json.NewDecoder(out)

	type Event struct {
		Status         string `json:"status"`
		Error          string `json:"error"`
		Progress       string `json:"progress"`
		ProgressDetail struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"progressDetail"`
	}

	var event *Event
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

func runContainer(ctx context.Context, c *client.Client, hostConfig *container.HostConfig, containerConfig *container.Config, containerName string) error {
	resp, err := c.ContainerCreate(ctx, containerConfig, hostConfig, &network.NetworkingConfig{}, nil, containerName)

	if err != nil {
		return err
	}

	if err := c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	okC, errC := c.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	for {
		select {
		case <-okC:
			return nil
		case err = <-errC:
			return err
		}
	}
}
