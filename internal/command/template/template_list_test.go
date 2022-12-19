package template

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/stretchr/testify/require"
)

const (
	FoundTemplates = "Available templates:"
	NoTemplates    = "No templates found! You can use 'fl template pull' to download some templates."
)

func TestListTemplates(t *testing.T) {
	ctx := context.TODO()
	var outbuf bytes.Buffer
	testLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).DisableAnimation().Build()

	t.Run("prints no templates found when no templates are available", func(t *testing.T) {
		listCmd := List{}

		err := listCmd.Run(ctx, testLogger)
		require.NoError(t, err)

		out := strings.Trim(outbuf.String(), "\n")
		require.Equal(t, NoTemplates, out)
	})

	t.Run("prints available templates", func(t *testing.T) {
		outbuf.Reset()

		tmpDir, err := os.MkdirTemp("", "funless-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		templatePullCmd := Pull{
			Repository: pkg.DefaultTemplateRepository,
			OutDir:     tmpDir,
		}

		templatePullCmd.Run(ctx, testLogger)

		listCmd := List{
			TemplateDir: tmpDir,
		}

		err = listCmd.Run(ctx, testLogger)
		require.NoError(t, err)

		require.Contains(t, outbuf.String(), FoundTemplates)
	})
}
