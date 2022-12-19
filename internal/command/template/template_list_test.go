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

		err = templatePullCmd.Run(ctx, testLogger)
		require.NoError(t, err)

		listCmd := List{
			TemplateDir: tmpDir,
		}

		err = listCmd.Run(ctx, testLogger)
		require.NoError(t, err)

		require.Contains(t, outbuf.String(), FoundTemplates)
	})
}
