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

package homedir

import (
	"errors"
	"os"
	"path/filepath"
)

var getHomeDir = os.UserHomeDir

// EnsureConfigDir return the path to the config directory.
// It creates it if needed.
func EnsureConfigDir() (string, error) {
	homedir, err := getHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(homedir, ".fl")
	// check if the directory exists
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// if not, create the directory
		if err := os.MkdirAll(path, 0755); err != nil {
			return "", err
		}
	}
	return path, nil
}

func WriteToConfigDir(filename string, data []byte, overwrite bool) error {
	homedir, err := EnsureConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(homedir, filename)
	if _, err := os.Stat(path); err == nil {
		if overwrite {
			os.Remove(path)
		} else {
			return errors.New("file already exists and overwrite is false")
		}
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return err
	}
	return nil
}

func ReadFromConfigDir(filename string) ([]byte, error) {
	homedir, err := EnsureConfigDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(homedir, filename)
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
