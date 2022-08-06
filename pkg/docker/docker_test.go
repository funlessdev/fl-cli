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

	"github.com/stretchr/testify/assert"
)

type fakeExec struct {
	out string
	err error
}

func (sh *fakeExec) runShellCmd(cmd string, args ...string) (string, error) {
	return sh.out, sh.err
}

func Test_dockerInfo(t *testing.T) {
	t.Run("should return cmd output when success", func(t *testing.T) {
		out, err := dockerInfo(&fakeExec{out: "ok", err: nil})
		assert.NoError(t, err)
		assert.Equal(t, "ok", out)
	})
	t.Run("should return 'docker is not running' error when error", func(t *testing.T) {
		_, err := dockerInfo(&fakeExec{out: "", err: assert.AnError})
		assert.ErrorContains(t, err, "docker is not running")
	})
}

func Test_dockerVersion(t *testing.T) {
	t.Run("should return cmd output when success", func(t *testing.T) {
		out, err := dockerVersion(&fakeExec{out: "ok", err: nil})
		assert.NoError(t, err)
		assert.Equal(t, "ok", out)
	})

	t.Run("should return cmd error when error", func(t *testing.T) {
		_, err := dockerInfo(&fakeExec{out: "", err: assert.AnError})
		assert.Error(t, err)
	})
}
