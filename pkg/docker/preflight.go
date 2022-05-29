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
	"fmt"
	"strings"

	"github.com/coreos/go-semver/semver"
)

const MinDockerVersion = "18.06.3-ce"

type preflightChecksPipeline struct {
	shell      shellExecutor
	dockerData string
	err        error
}

type checkStep func(pd *preflightChecksPipeline)

func (p *preflightChecksPipeline) step(f checkStep) {
	if p.err != nil {
		return
	}
	f(p)
}

// RunPreflightChecks performs preflight checks
// It ensures that docker is at least @MinDockerVersion.
// It returns an error if occured, nil otherwise
func RunPreflightChecks() error {

	// Preflight Checks pipeline
	pp := preflightChecksPipeline{shell: &baseShell{}}

	pp.step(extractDockerInfo)
	pp.step(ensureDockerVersion)

	return pp.err
}

func extractDockerInfo(p *preflightChecksPipeline) {
	p.dockerData, p.err = dockerInfo(p.shell)
}

func ensureDockerVersion(p *preflightChecksPipeline) {
	fmt.Printf("Check Docker version (min %s)\n", MinDockerVersion)
	// p.logger.StartSpinner("Check Docker version (min. " + MinDockerVersion + ")")
	version, err := dockerVersion(p.shell)
	if err != nil {
		p.err = err
		// p.logger.EndSpinner(false)
		return
	}
	vA := semver.New(MinDockerVersion)
	vB := semver.New(strings.TrimSpace(version))
	if vB.Compare(*vA) == -1 {
		p.err = fmt.Errorf("installed docker version %s is no longer supported", vB)
		// p.logger.EndSpinner(false)
		return
	}
	// p.logger.EndSpinner(true)
	fmt.Println("Docker supported")
}
