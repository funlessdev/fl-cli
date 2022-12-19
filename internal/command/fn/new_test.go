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

package fn

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	ctx := context.TODO()

	var out bytes.Buffer
	testLogger, _ := log.NewLoggerBuilder().WithWriter(&out).DisableAnimation().Build()

	t.Run("error when function name already exists", func(t *testing.T) {
		testOutDir, err := os.MkdirTemp("", "funless-test-new-outdir-")
		require.NoError(t, err)
		defer os.RemoveAll(testOutDir)

		// Create a folder with the same name of the function
		err = os.MkdirAll(filepath.Join(testOutDir, "test"), 0755)
		require.NoError(t, err)

		newCmd := New{
			Name:   "test",
			OutDir: testOutDir,
		}

		err = newCmd.Run(ctx, testLogger)
		require.Error(t, err)

		require.Contains(t, err.Error(), "already exists")
	})

	t.Run("no error when template dir is not found", func(t *testing.T) {
		out.Reset()

		testOutDir, err := os.MkdirTemp("", "funless-test-new-outdir-")
		require.NoError(t, err)
		defer os.RemoveAll(testOutDir)

		newCmd := New{
			Name:        "test",
			Language:    "js",
			TemplateDir: filepath.Join(testOutDir, "not-exists"),
			OutDir:      testOutDir,
		}

		err = newCmd.Run(ctx, testLogger)
		require.NoError(t, err)

		require.Contains(t, out.String(), "Pulling default templates!")
	})

	t.Run("error when template is not found", func(t *testing.T) {
		out.Reset()

		testOutDir, err := os.MkdirTemp("", "funless-test-new-outdir-")
		require.NoError(t, err)
		defer os.RemoveAll(testOutDir)

		newCmd := New{
			Name:        "test",
			Language:    "not-exists",
			TemplateDir: filepath.Join(testOutDir, "not-exists"),
			OutDir:      testOutDir,
		}

		err = newCmd.Run(ctx, testLogger)
		require.Error(t, err)

		require.Contains(t, err.Error(), "no valid template for")
	})

	t.Run("creates a new function when no error", func(t *testing.T) {
		out.Reset()

		testOutDir, err := os.MkdirTemp("", "funless-test-new-outdir-")
		require.NoError(t, err)
		defer os.RemoveAll(testOutDir)

		newCmd := New{
			Name:        "test",
			Language:    "js",
			TemplateDir: filepath.Join(testOutDir, "not-exists"),
			OutDir:      testOutDir,
		}

		err = newCmd.Run(ctx, testLogger)
		require.NoError(t, err)

		require.Contains(t, out.String(), "Pulling default templates!")
		require.Contains(t, out.String(), "Function \"test\" created!")

		// Check if the function folder is created
		_, err = os.Stat(filepath.Join(testOutDir, "test"))
		require.NoError(t, err)
	})
}
