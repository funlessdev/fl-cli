// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package docker

import (
	"errors"
	"os/exec"
	"regexp"
	"strings"
)

func dockerInfo(exec shellExecutor) (string, error) {
	out, err := exec.runShellCmd("docker info")
	if err != nil {
		return "", errors.New("docker is not running")
	}
	return out, nil
}

func dockerVersion(exec shellExecutor) (string, error) {
	return exec.runShellCmd("docker version --format {{.Server.Version}}")
}

type shellExecutor interface {
	runShellCmd(cmd string, args ...string) (string, error)
}

type baseShell struct{}

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

func (sh *baseShell) runShellCmd(cmd string, args ...string) (string, error) {
	exe, params := parseCmd(cmd, args...)
	out, err := exec.Command(exe, params...).CombinedOutput()
	return string(out), err
}
