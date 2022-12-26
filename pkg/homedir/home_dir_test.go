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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureConfigDir(t *testing.T) {
	homedirPath, err := os.MkdirTemp("", "funless-test-homedir-")
	require.NoError(t, err)

	GetHomeDir = func() (string, error) {
		return homedirPath, err
	}
	defer func() {
		GetHomeDir = os.UserHomeDir
		os.RemoveAll(homedirPath)
	}()

	t.Run("should create the config directory if it does not exist", func(t *testing.T) {
		path, err := EnsureConfigDir()
		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(homedirPath, ".fl"), path)
		assert.DirExists(t, path)
	})

	t.Run("should return the path to the config directory if it exists", func(t *testing.T) {
		path, err := EnsureConfigDir()
		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(homedirPath, ".fl"), path)
		assert.DirExists(t, path)
	})
}

func TestWriteToConfigDir(t *testing.T) {
	homedirPath, err := os.MkdirTemp("", "funless-test-homedir-")
	require.NoError(t, err)

	GetHomeDir = func() (string, error) {
		return homedirPath, err
	}
	defer func() {
		GetHomeDir = os.UserHomeDir
		os.RemoveAll(homedirPath)
	}()

	t.Run("should write the file to the config directory", func(t *testing.T) {
		path, err := WriteToConfigDir("test", []byte("test"), false)
		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(homedirPath, ".fl", "test"), path)
		assert.FileExists(t, path)
	})

	t.Run("should overwrite the file if overwrite is true", func(t *testing.T) {
		path, err := WriteToConfigDir("test", []byte("test"), true)
		assert.NoError(t, err)
		assert.FileExists(t, path)
	})

	t.Run("should return an error if the file already exists and overwrite is false", func(t *testing.T) {
		_, err := WriteToConfigDir("test", []byte("test"), false)
		assert.Error(t, err)
	})
}

func TestReadFromConfigDir(t *testing.T) {
	homedirPath, err := os.MkdirTemp("", "funless-test-homedir-")
	require.NoError(t, err)

	GetHomeDir = func() (string, error) {
		return homedirPath, err
	}
	defer func() {
		GetHomeDir = os.UserHomeDir
		os.RemoveAll(homedirPath)
	}()

	t.Run("should return an error if the file does not exist", func(t *testing.T) {
		_, _, err := ReadFromConfigDir("test")
		assert.Error(t, err)
	})

	t.Run("should return the file content if the file exists", func(t *testing.T) {
		_, err := WriteToConfigDir("test", []byte("test"), false)
		assert.NoError(t, err)
		content, path, err := ReadFromConfigDir("test")
		assert.NoError(t, err)
		assert.Equal(t, []byte("test"), content)
		assert.Equal(t, filepath.Join(homedirPath, ".fl", "test"), path)
	})

	t.Run("should return an error if the file is a directory", func(t *testing.T) {
		err := os.Mkdir(filepath.Join(homedirPath, ".fl", "test-dir"), 0755)
		assert.NoError(t, err)
		_, _, err = ReadFromConfigDir("test-dir")
		assert.Error(t, err)
	})
}

func TestCreateDirInConfigDir(t *testing.T) {
	homedirPath, err := os.MkdirTemp("", "funless-test-homedir-")
	require.NoError(t, err)

	GetHomeDir = func() (string, error) {
		return homedirPath, err
	}
	defer func() {
		GetHomeDir = os.UserHomeDir
		os.RemoveAll(homedirPath)
	}()

	t.Run("should create the directory in the config directory", func(t *testing.T) {
		path, err := CreateDirInConfigDir("test")
		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(homedirPath, ".fl", "test"), path)
		assert.DirExists(t, path)
	})

	t.Run("should return no error if the directory already exists", func(t *testing.T) {
		path, err := CreateDirInConfigDir("test")
		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(homedirPath, ".fl", "test"), path)
		assert.DirExists(t, path)
	})
}
