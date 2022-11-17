package deploy

import (
	"context"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/docker"
	"github.com/mitchellh/go-homedir"
)

type DockerDeployer interface {
	Setup(coreImg, workerImg string) error
	CreateFLNetwork(ctx context.Context) error
	PullCoreImage(ctx context.Context) error
	PullWorkerImage(ctx context.Context) error
	PullPromImage(ctx context.Context) error
	StartCore(ctx context.Context) error
	StartWorker(ctx context.Context) error
	StartProm(ctx context.Context) error
}

type FLDockerDeployer struct {
	dockerClient *client.Client
	image        docker.ImageHandler
	container    docker.ContainerHandler
	network      docker.NetworkHandler

	logsPath            string
	flNetId             string
	flNetName           string
	coreImg             string
	coreContainerName   string
	workerImg           string
	workerContainerName string
	promContainerName   string
}

func NewDockerDeployer(img docker.ImageHandler, ctr docker.ContainerHandler, ntw docker.NetworkHandler) DockerDeployer {
	d := &FLDockerDeployer{
		image:     img,
		container: ctr,
		network:   ntw,
	}
	return d
}

func (d *FLDockerDeployer) Setup(coreImg, workerImg string) error {
	d.coreImg = coreImg
	d.workerImg = workerImg

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return err
	}
	d.dockerClient = cli

	h, err := homedir.Dir()
	if err != nil {
		return err
	}
	logsPath := filepath.Join(h, "funless-logs")
	if err := os.MkdirAll(logsPath, 0755); err != nil {
		return err
	}

	d.logsPath = logsPath
	return nil
}

func (d *FLDockerDeployer) CreateFLNetwork(ctx context.Context) error {
	// Network for Core + Worker
	exists, id, err := d.network.Exists(ctx, d.dockerClient, d.flNetName)
	if err != nil {
		return err
	}
	if exists {
		d.flNetId = id
		return nil
	}
	id, err = d.network.Create(ctx, d.dockerClient, d.flNetName)
	if err != nil {
		return err
	}
	d.flNetId = id
	return nil
}

func (d *FLDockerDeployer) PullCoreImage(ctx context.Context) error {
	return d.pull(ctx, d.coreImg)
}

func (d *FLDockerDeployer) PullWorkerImage(ctx context.Context) error {
	return d.pull(ctx, d.workerImg)
}

func (d *FLDockerDeployer) PullPromImage(ctx context.Context) error {
	return d.pull(ctx, pkg.PrometheusImg)
}

func (d *FLDockerDeployer) StartCore(ctx context.Context) error {
	containerConfig := coreContainerConfig(d.coreImg)
	hostConfig := coreHostConfig(d.logsPath)
	netConf := networkConfig(d.flNetName, d.flNetId)
	configs := configs(d.coreContainerName, containerConfig, hostConfig, netConf)
	return d.container.RunAsync(ctx, d.dockerClient, configs)
}

func (d *FLDockerDeployer) StartWorker(ctx context.Context) error {
	containerConfig := workerContainerConfig(d.workerImg)
	hostConf := workerHostConfig(d.logsPath)
	netConf := networkConfig(d.flNetName, d.flNetId)
	configs := configs(d.workerContainerName, containerConfig, hostConf, netConf)
	return d.container.RunAsync(ctx, d.dockerClient, configs)
}

func (d *FLDockerDeployer) StartProm(ctx context.Context) error {
	containerConfig := promContainerConfig()
	hostConf := promHostConfig()
	netConf := networkConfig(d.flNetName, d.flNetId)
	configs := configs(d.promContainerName, containerConfig, hostConf, netConf)
	return d.container.RunAsync(ctx, d.dockerClient, configs)
}

func (d *FLDockerDeployer) pull(ctx context.Context, img string) error {
	exists, err := d.image.Exists(ctx, d.dockerClient, img)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return d.image.Pull(ctx, d.dockerClient, img)

}

func coreContainerConfig(coreImg string) *container.Config {
	return &container.Config{
		Image: coreImg,
		ExposedPorts: nat.PortSet{
			"4000/tcp": struct{}{},
		},
		Env:     []string{"SECRET_KEY_BASE=" + pkg.CoreDevSecretKey},
		Volumes: map[string]struct{}{},
	}
}
func coreHostConfig(logsPath string) *container.HostConfig {
	return &container.HostConfig{
		PortBindings: nat.PortMap{
			"4000/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "4000",
				},
			},
		},
		Mounts: []mount.Mount{
			{
				Source: logsPath,
				Target: "/tmp/funless",
				Type:   mount.TypeBind,
			},
		},
	}
}
func workerContainerConfig(workerImg string) *container.Config {
	return &container.Config{
		Image: workerImg,
	}
}
func workerHostConfig(logsPath string) *container.HostConfig {
	return &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Source: logsPath,
				Target: "/tmp/funless",
				Type:   mount.TypeBind,
			},
		},
	}
}
func promContainerConfig() *container.Config {
	return &container.Config{
		Image: pkg.PrometheusImg,
	}
}
func promHostConfig() *container.HostConfig {
	return &container.HostConfig{
		PortBindings: nat.PortMap{
			"9090/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "9090",
				},
			},
		},
	}
}
func networkConfig(networkName, networkID string) *network.NetworkingConfig {
	endpoints := make(map[string]*network.EndpointSettings, 1)
	endpoints[networkName] = &network.EndpointSettings{
		NetworkID: networkID,
	}

	return &network.NetworkingConfig{
		EndpointsConfig: endpoints,
	}
}

func configs(name string, c *container.Config, h *container.HostConfig, n *network.NetworkingConfig) docker.ContainerConfigs {
	return docker.ContainerConfigs{
		ContName:   name,
		Container:  c,
		Host:       h,
		Networking: n,
	}
}
