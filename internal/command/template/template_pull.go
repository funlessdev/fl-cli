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
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/go-git/go-git/v5"
)

type Pull struct {
	Repository string `arg:"" default:"https://github.com/funlessdev/fl-templates.git" help:"The repository to pull the template folder from"`
	OutDir     string `short:"o" type:"existingdir" default:"." help:"The output directory where the template folder will be placed"`
	Force      bool   `short:"f" default:"false" help:"Overwrite the template if it already exists"`
}

func (f *Pull) Help() string {
	return `
DESCRIPTION

	Pull template folder from a repository, the default one is 
	https://github.com/funlessdev/fl-templates.git.
	An other repository can be used as argument to override the default one.
	The "--out-dir" can be used to specify a different path for the output 
	other than the default one. 
	The "--force" can be used to overwrite the template if it already exists.

EXAMPLES

	$ ls

		my_subfolder

	$ fl template pull <your-template-repository> --out-dir <your-template-output-dir> --force

	$ ls

		my_subfolder template

---

	For default template repository:

	$ fl template pull

	$ ls ./template/

		js rust
`

}

type copyOpts struct {
	// Map with the templates that are ok to move (e.g. "rust" => true)
	okToCopyTemplate map[string]bool

	notCopiedTemplates []string
	copiedTemplates    []string

	srcDir  string
	destDir string

	force bool

	err error
}

func (p *Pull) Run(ctx context.Context, logger log.FLogger) error {
	// Create a tmp tmpDir for the repository

	_ = logger.StartSpinner("Cloning templates repository...")
	tmpDir, err := os.MkdirTemp("", "funless-templates-")
	if err != nil {
		return logger.StopSpinner(err)
	}
	defer os.RemoveAll(tmpDir)

	// Clone the repository in the tmp dir
	err = logger.StopSpinner(cloneRepo(p.Repository, tmpDir))
	if err != nil {
		return err
	}

	copyOpts := &copyOpts{
		okToCopyTemplate:   make(map[string]bool),
		notCopiedTemplates: make([]string, 0),
		copiedTemplates:    make([]string, 0),
		srcDir:             filepath.Join(tmpDir, "template"),
		destDir:            filepath.Join(p.OutDir, "template"),
		force:              p.Force,
	}

	_ = logger.StartSpinner("Preparing templates...")
	// Move the templates from the repository to the template folder
	err = logger.StopSpinner(copyTemplates(copyOpts))
	if err != nil {
		return err
	}

	if len(copyOpts.notCopiedTemplates) > 0 {
		logger.Infof("Skipped %d template(s) (already present): %v\n", len(copyOpts.notCopiedTemplates), copyOpts.notCopiedTemplates)
	}

	if len(copyOpts.copiedTemplates) == 0 {
		logger.Info("No new templates retrieved.")
	} else {
		logger.Infof("Retrieved %d templates from %s : %v\n", len(copyOpts.copiedTemplates), p.Repository, copyOpts.copiedTemplates)
	}

	return nil
}

func cloneRepo(url, dir string) error {
	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:   url,
		Depth: 1,
	})
	if err != nil {
		return err
	}

	return nil
}

func copyTemplates(c *copyOpts) error {
	templates, err := os.ReadDir(c.srcDir)
	if err != nil {
		return errors.New("cannot find templates. Folder 'template' missing?")
	}

	for _, dir := range templates {
		copySingleTemplate(dir, c)
	}
	return c.err
}

func copySingleTemplate(dir fs.DirEntry, c *copyOpts) {
	// Skip non-directories
	if !dir.IsDir() {
		return
	}
	language := dir.Name()

	// if we don't know we can copy the template, check if it already exists
	if _, found := c.okToCopyTemplate[language]; !found {
		c.okToCopyTemplate[language] = canCopyTemplate(c.destDir, language) || c.force
	}

	if c.okToCopyTemplate[language] {
		c.copiedTemplates = append(c.copiedTemplates, language)

		// Now actually copy the template
		src := filepath.Join(c.srcDir, language)
		dest := filepath.Join(c.destDir, language)
		if err := pkg.Copy(src, dest); err != nil {
			c.err = err
		}
	} else {
		c.notCopiedTemplates = append(c.notCopiedTemplates, language)
	}
}

// canCopyTemplate checks if the folder of a particular language exists.
// If it does, it returns false (no copy needed), otherwise true
func canCopyTemplate(destDir string, language string) bool {
	dir := filepath.Join(destDir, language)
	if _, err := os.Stat(dir); err == nil {
		return false
	}
	return true
}
