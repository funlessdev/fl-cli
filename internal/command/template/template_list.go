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
	TemplateDir string `short:"d" type:"path" default:"./template" help:"the directory to read the templates from"`
}

func (l *List) Run(ctx context.Context, logger log.FLogger) error {
	tpath := filepath.Join(l.TemplateDir, "template")
	var templates []string

	templateFolders, err := os.ReadDir(tpath)
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
