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
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type New struct {
	Name        string `arg:"" help:"the name of the function"`
	Language    string `name:"lang" short:"l" required:"" xor:"list-lang" enum:"rust, js" help:"the language of the function"`
	TemplateDir string `short:"t" type:"path" default:"./template" help:"the directory where the template are located"`
}

func (n *New) Run(ctx context.Context, logger log.FLogger) error {
	// check that Language is an existing template
	if !isValidTemplate(n.TemplateDir, n.Language) {
		return fmt.Errorf("no valid template for \"%s\" found", n.Language)
	}

	// copy template to current directory
	src := filepath.Join(n.TemplateDir, n.Language)
	dst := filepath.Join(".", n.Name)
	pkg.Copy(src, dst)

	logger.Info("Implementing...")
	return nil
}

func isValidTemplate(tDir string, lang string) bool {
	path := filepath.Join(tDir, lang)
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
