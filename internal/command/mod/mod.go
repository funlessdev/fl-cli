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

package mod

type Mod struct {
	Get    Get    `cmd:"" aliases:"g" help:"List functions and information of a module"`
	Delete Delete `cmd:"" aliases:"d,rm" help:"Delete a module"`
	Update Update `cmd:"" aliases:"u,up" help:"Update the name of a module"`
	Create Create `cmd:"" aliases:"c" help:"Create a new module"`
	List   List   `cmd:"" aliases:"l,ls" help:"List all modules"`
}
