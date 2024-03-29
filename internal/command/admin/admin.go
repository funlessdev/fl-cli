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

package admin

import (
	deploy "github.com/funlessdev/fl-cli/internal/command/admin/deploy"
	user "github.com/funlessdev/fl-cli/internal/command/admin/user"
)

type Admin struct {
	Deploy deploy.Deploy `cmd:"" name:"deploy" aliases:"d" help:"Deploy FunLess on different setups"`
	User   user.User     `cmd:"" name:"user" aliases:"u" help:"Create/delete FunLess users"`
}
