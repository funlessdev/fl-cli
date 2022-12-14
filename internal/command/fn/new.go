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
	"sort"

	"github.com/funlessdev/fl-cli/pkg/log"
)

type New struct {
	Name     string `arg:"" help:"the name of the function"`
	Language string `name:"lang" short:"l" required:"" xor:"list-lang" enum:"rust, js" help:"the language of the function"`
	List     bool   `aliases:"ls" xor:"list-lang" required:"" help:"list available templates in the current folder"`
}

const templateDirectory = "./template/"

func (n *New) Run(ctx context.Context, logger log.FLogger) error {
	if n.List {
		return listTemplates(logger)
	}

	logger.Info("Not implemented yet!")
	return nil
}

func listTemplates(logger log.FLogger) error {
	var templates []string

	templateFolders, err := os.ReadDir(templateDirectory)
	if os.IsNotExist(err) {
		logger.Info("No templates found! You can use 'fl template pull' to download some templates.")
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
