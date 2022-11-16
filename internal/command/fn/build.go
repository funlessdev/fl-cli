package fn

import (
	"context"
	"fmt"

	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/build"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Build struct {
	Name        string `arg:"" help:"the name of the function"`
	Source      string `short:"s" required:"" xor:"dir-file,dir-build" type:"existingdir" help:"path of the source directory"`
	Destination string `short:"d" type:"path" help:"path where the compiled wasm file will be saved"`
	Language    string `short:"l" enum:"rust,js" required:"" help:"programming language of the function"`
}

func (b *Build) Run(ctx context.Context, builder build.DockerBuilder, logger log.FLogger) error {
	logger.Info(fmt.Sprintf("Building %s into a wasm binary...\n", b.Name))

	_ = logger.StartSpinner("Setting up...")
	if build_err := logger.StopSpinner(builder.Setup(ctx, b.Language, b.Destination)); build_err != nil {
		return build_err
	}

	_ = logger.StartSpinner(fmt.Sprintf("Pulling %s builder image (%s) 📦", langNames[b.Language], pkg.FLBuilderImages[b.Language]))
	if build_err := logger.StopSpinner(builder.PullBuilderImage(ctx)); build_err != nil {
		return build_err
	}
	_ = logger.StartSpinner("Building source... 🛠️")
	if build_err := logger.StopSpinner(builder.BuildSource(ctx, b.Source)); build_err != nil {
		return build_err
	}

	logger.Info(fmt.Sprintf("\nSuccessfully built function at %s/%s.wasm 🥳🥳", b.Destination, b.Name))
	return nil
}

var langNames = map[string]string{
	"js":   "Javascript",
	"rust": "Rust",
}
