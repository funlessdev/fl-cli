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

	"github.com/funlessdev/fl-cli/internal/command/template"
	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type New struct {
	Name        string `arg:"" help:"the name of the function"`
	Language    string `name:"lang" short:"l" required:"" enum:"rust, js" help:"the language of the function"`
	TemplateDir string `short:"t" type:"path" default:"." help:"the directory where the template are located"`
	OutDir      string `short:"o" type:"path" default:"." help:"the directory where the function will be created"`
}

func (n *New) Run(ctx context.Context, logger log.FLogger) error {
	srcLanguageTemplate := filepath.Join(n.TemplateDir, "template", n.Language)
	destFunc := filepath.Join(n.OutDir, n.Name)

	// Check that function is not already present
	if folderExists(destFunc) {
		return fmt.Errorf("function \"%s\" already exists", n.Name)
	}

	// if template folder not found, pull default templates
	if !folderExists(filepath.Join(n.TemplateDir, "template")) {
		logger.Infof("Folder \"template\" not found in %s. Pulling default templates!\n", n.TemplateDir)
		pullCmd := template.Pull{
			Repository: pkg.DefaultTemplateRepository,
			OutDir:     n.TemplateDir,
		}
		if err := pullCmd.Run(ctx, logger); err != nil {
			return err
		}
	}

	// if language template is still not available, return error
	if !folderExists(srcLanguageTemplate) {
		return fmt.Errorf("no valid template for \"%s\" found", n.Language)
	}

	// copy template to current directory with the name of the function
	if err := pkg.Copy(srcLanguageTemplate, destFunc); err != nil {
		return err
	}

	logger.Infof("Function \"%s\" created!", n.Name)

	return nil
}

func folderExists(template string) bool {
	if _, err := os.Stat(template); err == nil {
		return true
	}
	return false
}
