package deploy

import "context"

type DockerDeployer interface {
	SetupClient(ctx context.Context) error
	SetupFLNetworks(ctx context.Context) error
	PullCoreImage(ctx context.Context) error
	PullWorkerImage(ctx context.Context) error
	StartCore(ctx context.Context) error
	StartWorker(ctx context.Context) error
}
