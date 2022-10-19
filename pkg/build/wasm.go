package build

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/funlessdev/fl-cli/pkg"
)

func BuildSource(ctx context.Context, c *client.Client, srcPath string, language string) error {
	image := pkg.FLRuntimes[language]
	absPath, err := filepath.Abs(srcPath)
	if err != nil {
		return err
	}
	outPath, _ := filepath.Abs("./out_wasm")
	err = os.MkdirAll(outPath, 0700)
	if err != nil {
		return err
	}

	containerConfig := &container.Config{
		Image:   image,
		Volumes: map[string]struct{}{},
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Source:   absPath,
				Target:   "/lib_fl/",
				ReadOnly: true,
				Type:     mount.TypeBind,
			},
			{
				Source: outPath,
				Target: "/out_wasm",
				Type:   mount.TypeBind,
			},
		},
	}
	return runContainer(ctx, c, hostConfig, containerConfig, pkg.FLRuntimeNames[language])
}

func GetBuilderImage(ctx context.Context, c *client.Client, language string) error {
	image, exists := pkg.FLRuntimes[language]

	if !exists {
		return errors.New("No corresponding builder found for the given language")
	}

	return pullImage(ctx, c, image)
}
