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

package pkg

const CoreDevSecretKey = "dev-secret-key"

const CoreImg = "ghcr.io/funlessdev/fl-core:latest"
const WorkerImg = "ghcr.io/funlessdev/fl-worker:latest"
const PrometheusImg = "giusdp/fl-prometheus:latest"

const CoreContName = "fl-core"
const WorkerContName = "fl-worker"
const PrometheusContName = "prometheus"
const FLNet = "fl-net"

var FLBuilderImages = map[string]string{
	"js":   "ghcr.io/funlessdev/fl-js-builder:latest",
	"rust": "ghcr.io/funlessdev/fl-rust-builder:latest",
}

var FLBuilderNames = map[string]string{
	"js":   "fl-js-builder",
	"rust": "fl-rust-builder",
}
