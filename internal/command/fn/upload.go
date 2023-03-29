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
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/log"
)

type Upload struct {
	Name   string `arg:"" help:"The name of the function"`
	Source string `arg:"" type:"existingfile" help:"Path of the wasm binary"`
	Module string `short:"m" default:"_" help:"The module of the function"`
}

func (c *Upload) Help() string {
	return `
DESCRIPRION

	It uploads a function with the specified name from the specified 
	path of the wasm binary.
	The "--module" flag can be used to choose a module other than 
	the default one. 

EXAMPLES
	
	$ fl fn upload <your-function-name> <your-wasm-path> --module=<your-module-name>
`
}

func (u *Upload) Run(ctx context.Context, fnHandler client.FnHandler, logger log.FLogger) error {
	_ = logger.StartSpinner("Reading wasm...")
	code, err := openWasmFile(u.Source)
	if err != nil {
		return logger.StopSpinner(err)
	}
	_ = logger.StopSpinner(nil)

	_ = logger.StartSpinner("Uploading function...")
	err = fnHandler.Create(ctx, u.Name, u.Module, code)
	if err != nil {
		return logger.StopSpinner(err)
	}
	_ = logger.StopSpinner(nil)

	logger.Info(fmt.Sprintf("Successfully uploaded function %s/%s 👌\n", u.Module, u.Name))
	return nil
}

// I swap the function with a fake one in create_test.go
var openWasmFile = func(path string) (*os.File, error) {
	return readWasmFile(path)
}

func readWasmFile(path string) (*os.File, error) {
	if !strings.HasSuffix(path, ".wasm") {
		return nil, errors.New("can only create function with a .wasm file")
	}

	wasmPath := filepath.Clean(path)

	code, err := os.Open(wasmPath)
	if err != nil {
		return nil, err
	}
	stat, err := code.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, errors.New("the specified file is a directory")
	}

	wasmMagicHeader := []byte{0x00, 0x61, 0x73, 0x6d}
	var buf = make([]byte, 4)
	if _, err := code.Read(buf); err != nil {
		return nil, err
	}

	if !bytes.Equal(buf, wasmMagicHeader) {
		return nil, errors.New("the file is not a valid wasm binary")
	}

	// Reset the file pointer to the beginning of the file
	_, err = code.Seek(0, io.SeekStart)
	return code, err
}
