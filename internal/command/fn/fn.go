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

type Fn struct {
	Invoke Invoke `cmd:"" aliases:"i" help:"Invoke a function"`
	Create Create `cmd:"" aliases:"c" help:"A combination of build and upload to create a function"`
	Delete Delete `cmd:"" aliases:"d" help:"Delete an existing function"`
	Build  Build  `cmd:"" aliases:"b" help:"Compile a function into a wasm binary"`
	Upload Upload `cmd:"" aliases:"up" help:"Create functions by uploading wasm binaries"`
	New    New    `cmd:"" aliases:"n" help:"Create a new function from a template"`

	Host string `short:"H" help:"API host/port of the platform (no protocol)"`
}
