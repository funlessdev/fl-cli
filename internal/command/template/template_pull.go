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
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/go-git/go-git/v5"
)

type Pull struct {
	Repository string `arg:"" default:"https://github.com/funlessdev/fl-templates.git" help:"the repository to pull the template folder from"`
	OutDir     string `short:"o" default:"." help:"the output directory where the template folder will be placed"`
	Force      bool   `short:"f" default:"false" help:"overwrite the template if it already exists"`
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
	logger.Infof("Retrieved %d templates from %s : %v", len(copyOpts.copiedTemplates), p.Repository, copyOpts.copiedTemplates)
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
		if err := copy(c.srcDir, c.destDir, language); err != nil {
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

func copy(srcDir string, destDir string, entry string) error {
	src := filepath.Join(srcDir, entry)
	dest := filepath.Join(destDir, entry)

	// Get properties of source
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	// If it's a directory, create it and copy recursively
	if info.IsDir() {
		// 1. Create the dest directory
		if err := os.MkdirAll(dest, info.Mode()); err != nil {
			return fmt.Errorf("error creating: %s - %s", dest, err.Error())
		}

		// 2. Read the source directory
		entries, err := os.ReadDir(src)
		if err != nil {
			// If we fail to read the source directory, remove the created directory
			os.RemoveAll(dest)
			return err
		}

		// 3. For each entry in the directory, recursively copy it
		for _, dirEntry := range entries {
			if err := copy(src, dest, dirEntry.Name()); err != nil {
				return err
			}
		}
		return nil
	}

	// If it's a file, copy it
	return copySingleFile(src, dest)
}

func copySingleFile(src, dest string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = ensureBaseDir(dest)
	if err != nil {
		return fmt.Errorf("error creating base directory: %s", err.Error())
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error creating dest file: %s", err.Error())
	}
	defer f.Close()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return fmt.Errorf("error setting dest file mode: %s", err.Error())
	}

	s, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening src file: %s", err.Error())
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	if err != nil {
		return fmt.Errorf("Error copying dest file: %s\n" + err.Error())
	}

	return nil
}

// ensureBaseDir creates the base directory of a given file path, if it does not exist.
func ensureBaseDir(fpath string) error {
	baseDir := path.Dir(fpath)
	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(baseDir, 0755)
}
