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
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/funlessdev/fl-cli/pkg"
)

type DockerShell interface {
	ComposeUp(ctx context.Context, composeFilePath string) error
	ComposeDown(ctx context.Context, composeFilePath string) error
	ComposeList(ctx context.Context) ([]string, error)
}

type FLDockerShell struct{}

func (sh *FLDockerShell) ComposeUp(ctx context.Context, composeFilePath string) error {
	return runShellCmd(ctx, os.Stdout, os.Stderr, "docker", "compose", "-f", composeFilePath, "up", "-d")
}

func (sh *FLDockerShell) ComposeDown(ctx context.Context, composeFilePath string) error {
	return runShellCmd(ctx, os.Stdout, os.Stderr, "docker", "compose", "-f", composeFilePath, "down")
}

func (sh *FLDockerShell) ComposeList(ctx context.Context) ([]string, error) {
	var buf bytes.Buffer
	err := runShellCmd(ctx, &buf, os.Stderr, "docker", "compose", "ls", "-q")
	lines := strings.Split(buf.String(), "\n")
	return lines, err
}

func runShellCmd(ctx context.Context, resultBuf io.Writer, errorBuf io.Writer, cmd string, args ...string) error {
	exe, params := parseCmd(ctx, cmd, args...)
	command := exec.Command(exe, params...)
	command.Stdout = resultBuf
	command.Stderr = errorBuf

	ctxEnv, ok := ctx.Value(pkg.FLContextKey("env")).(map[string]string)
	command.Env = os.Environ()

	if ok && ctxEnv != nil {
		for k := range ctxEnv {
			if ctxEnv[k] != "" {
				command.Env = append(command.Env, fmt.Sprintf("%s=%s", k, ctxEnv[k]))
			}
		}
	}

	return command.Run()
}

func parseCmd(ctx context.Context, cmd string, args ...string) (string, []string) {
	re := regexp.MustCompile(`[\r\t\n\f ]+`)
	a := strings.Split(re.ReplaceAllString(cmd, " "), " ")

	params := args
	if len(a) > 1 {
		params = append(a[1:], args...)
	}
	exe := a[0]

	return exe, params
}
