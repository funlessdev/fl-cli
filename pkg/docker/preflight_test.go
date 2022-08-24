// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package docker

import (
	"testing"

	"github.com/funlessdev/funless-cli/pkg/log"
	"github.com/stretchr/testify/assert"
)

type fakeShell struct {
	out string
	err error
}

func (sh *fakeShell) runShellCmd(cmd string, args ...string) (string, error) {
	return sh.out, sh.err
}

func Test_ensureDockerVersion(t *testing.T) {
	l, _ := log.NewLoggerBuilder().Build()
	p := preflightChecksPipeline{shell: &fakeShell{out: "19.03.5", err: nil}, logger: l}
	p.step(ensureDockerVersion)
	assert.NoError(t, p.err)

	p = preflightChecksPipeline{shell: &fakeShell{out: "10.03.5", err: nil}, logger: l}
	p.step(ensureDockerVersion)
	assert.ErrorContains(t, p.err, "installed docker version 10.3.5 is no longer supported")

	p = preflightChecksPipeline{shell: &fakeShell{out: MinDockerVersion, err: nil}, logger: l}
	p.step(ensureDockerVersion)
	assert.NoError(t, p.err)

	p = preflightChecksPipeline{shell: &fakeShell{out: "", err: assert.AnError}, logger: l}
	p.step(ensureDockerVersion)
	assert.Error(t, p.err)
}
