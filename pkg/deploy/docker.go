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
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type DockerShell interface {
	ComposeUp(composeFilePath string) error
	ComposeDown(composeFilePath string) error
}

type FLDockerShell struct{}

func (sh *FLDockerShell) ComposeUp(composeFilePath string) error {
	return runShellCmd("docker", "compose", "-f", composeFilePath, "up", "-d")
}

func (sh *FLDockerShell) ComposeDown(composeFilePath string) error {
	return runShellCmd("docker", "compose", "-f", composeFilePath, "down")
}

func runShellCmd(cmd string, args ...string) error {
	exe, params := parseCmd(cmd, args...)
	command := exec.Command(exe, params...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func parseCmd(cmd string, args ...string) (string, []string) {
	re := regexp.MustCompile(`[\r\t\n\f ]+`)
	a := strings.Split(re.ReplaceAllString(cmd, " "), " ")

	params := args
	if len(a) > 1 {
		params = append(a[1:], args...)
	}
	exe := a[0]

	return exe, params
}
