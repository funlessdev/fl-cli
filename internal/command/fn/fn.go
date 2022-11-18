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

import (
	"encoding/json"
	"errors"

	openapi "github.com/funlessdev/fl-client-sdk-go"
)

type Fn struct {
	Invoke Invoke `cmd:"" help:"run a function on an instance of the platform"`
	Create Create `cmd:"" help:"build and upload a function to the platform"`
	Delete Delete `cmd:"" help:"delete an existing function from the platform"`
	Build  Build  `cmd:"" help:"compile a function into a wasm binary"`
	// Upload Upload `cmd:"" help:"upload the wasm binary of a function to the platform"`
}

type FnError struct {
	Errors struct {
		Detail string `json:"detail"`
	} `json:"errors"`
}

func extractError(err error) error {
	var e FnError
	openApiError, castOk := err.(*openapi.GenericOpenAPIError)
	if !castOk {
		return err
	}
	if err := json.Unmarshal(openApiError.Body(), &e); err != nil {
		return err
	}
	return errors.New(e.Errors.Detail)
}

func (f *Fn) Help() string {
	return "More Description!"
}
