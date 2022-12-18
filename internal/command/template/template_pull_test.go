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

package template

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/stretchr/testify/require"
)

func TestTemplatePull(t *testing.T) {
	ctx := context.TODO()
	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).DisableAnimation().Build()

	t.Run("template pull", func(t *testing.T) {
		testOutDir, err := os.MkdirTemp("", "funless-test-templates-outdir-")
		require.NoError(t, err)
		defer os.RemoveAll(testOutDir)

		templatePullCmd := Pull{
			Repository: pkg.DefaultTemplateRepository,
			Force:      false,
			OutDir:     testOutDir,
		}

		err = templatePullCmd.Run(ctx, testLogger)
		require.NoError(t, err)

		// Verify template directory exists
		_, err = os.Stat(filepath.Join(testOutDir, "template"))
		require.NoError(t, err)
	})

	t.Run("template pull with force", func(t *testing.T) {
		testOutDir, err := os.MkdirTemp("", "funless-test-templates-outdir-")
		require.NoError(t, err)
		defer os.RemoveAll(testOutDir)

		templatePullCmd := Pull{
			Repository: pkg.DefaultTemplateRepository,
			Force:      false,
			OutDir:     testOutDir,
		}

		// 1. first pull to have the template directory
		err = templatePullCmd.Run(ctx, testLogger)
		require.NoError(t, err)

		var out bytes.Buffer
		logger, _ := log.NewLoggerBuilder().WithWriter(&out).DisableAnimation().Build()

		// 2. second pull to check that it can't overwrite the existing template directory
		err = templatePullCmd.Run(ctx, logger)
		require.NoError(t, err)

		require.Contains(t, out.String(), "Skipped")

		// 3. third pull with force to check that it can overwrite the existing template directory
		out.Reset()
		templatePullCmd.Force = true
		err = templatePullCmd.Run(ctx, logger)
		require.NoError(t, err)
		require.NotContains(t, out.String(), "Skipped")

		// Verify template directory exists
		_, err = os.Stat(filepath.Join(testOutDir, "template"))
		require.NoError(t, err)
	})

	t.Run("template pull with invalid url", func(t *testing.T) {
		templatePullCmd := Pull{
			Repository: "invalid-url",
			Force:      false,
		}

		err := templatePullCmd.Run(ctx, testLogger)
		require.Error(t, err)
	})
}
