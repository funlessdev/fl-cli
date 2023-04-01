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
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/funlessdev/fl-cli/pkg/log"
)

type List struct {
	TemplateDir string `short:"d" type:"path" default:"." help:"The directory to read the templates from"`
}

func (f *List) Help() string {
	return `
DESCRIPTION

	List all available templates.
	The "--template-dir" can be used to specify a different path other than 
	the default one.

EXAMPLES

	$ fl template list --template-dir <your-templates-path>`
}

func (l *List) Run(ctx context.Context, logger log.FLogger) error {
	tpath := filepath.Join(l.TemplateDir, "template")
	var templates []string

	templateFolders, err := os.ReadDir(tpath)
	if os.IsNotExist(err) {
		logger.Info("No templates found! You can use 'fl template pull' to download some templates.\n")
		return nil
	}

	for _, file := range templateFolders {
		if file.IsDir() {
			templates = append(templates, file.Name())
		}
	}

	logger.Infof("Available templates:\n%s\n", formatTemplateList(templates))

	return nil
}

func formatTemplateList(availableTemplates []string) string {
	var result string
	sort.Slice(availableTemplates, func(i, j int) bool {
		return availableTemplates[i] < availableTemplates[j]
	})
	for _, template := range availableTemplates {
		result += fmt.Sprintf("- %s\n", template)
	}
	return result
}
