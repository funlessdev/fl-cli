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

package pkg

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

func Copy(src string, dest string) error {
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
			newSrc := filepath.Join(src, dirEntry.Name())
			newDest := filepath.Join(dest, dirEntry.Name())
			if err := Copy(newSrc, newDest); err != nil {
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
