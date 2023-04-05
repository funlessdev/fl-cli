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

package deploy

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseCmd(t *testing.T) {
	t.Run("should split cmd from params in cmd string", func(t *testing.T) {
		testCtx := context.Background()
		exe, _ := parseCmd(testCtx, "docker info")
		assert.Equal(t, "docker", exe)
	})
	t.Run("should append arg in cmd string at the start of params array", func(t *testing.T) {
		testCtx := context.Background()
		_, params := parseCmd(testCtx, "docker info", "hello")
		assert.Equal(t, []string{"info", "hello"}, params)
	})
}

func Test_runShellCmd(t *testing.T) {

	t.Run("should return the command output in output buffer", func(t *testing.T) {
		var outBuf bytes.Buffer
		testCtx := context.Background()

		err := runShellCmd(testCtx, &outBuf, os.Stderr, "echo", "hello")
		assert.Nil(t, err)
		assert.Equal(t, "hello\n", outBuf.String())
	})

	t.Run("should return in error buffer all the contents of stderr", func(t *testing.T) {
		var errBuf bytes.Buffer
		testCtx := context.Background()

		errSecond := runShellCmd(testCtx, os.Stdout, &errBuf, "/bin/sh", "-c", "echo hello err 1>&2")
		assert.Nil(t, errSecond)
		assert.Equal(t, "hello err\n", errBuf.String())

	})

	t.Run("should return an error in case of command fail", func(t *testing.T) {
		testCtx := context.Background()

		err := runShellCmd(testCtx, os.Stdout, os.Stderr, "exit", "1")
		assert.NotNil(t, err)
	})

	t.Run("should return the content of StdOut and StdErr despite the command fail", func(t *testing.T) {
		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		testCtx := context.Background()

		errFirst := runShellCmd(testCtx, &outBuf, &errBuf, "echo", "hello out")
		assert.Nil(t, errFirst)

		errSecond := runShellCmd(testCtx, os.Stdout, &errBuf, "/bin/sh", "-c", "echo hello err 1>&2")
		assert.Nil(t, errSecond)

		errThird := runShellCmd(testCtx, &outBuf, &errBuf, "exit", "1")
		assert.NotNil(t, errThird)

		assert.Equal(t, "hello out\n", outBuf.String())
		assert.Equal(t, "hello err\n", errBuf.String())

	})

}
