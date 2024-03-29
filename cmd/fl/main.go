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

package main

import (
	"fmt"
	"os"

	"github.com/funlessdev/fl-cli/cmd/fl/app"
)

// CLIVersion holds the current version, to be set by the build with go build -ldflags "-X main.FLVersion=<version>"
var FLVersion = "vX.Y.Z.build"

func main() {
	if ctx, err := app.ParseCMD(FLVersion); err == nil {
		ctx.FatalIfErrorf(app.Run(ctx))
	} else {
		fmt.Println("Error parsing command line arguments:", err.Error())
		os.Exit(1)
	}
}
