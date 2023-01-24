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

import (
	"bytes"
	"encoding/json"
	"errors"

	openapi "github.com/funlessdev/fl-client-sdk-go"
)

type Mod struct {
	Get    Get    `cmd:"" aliases:"g" help:"list functions and information of a module"`
	Delete Delete `cmd:"" aliases:"d,rm" help:"delete a module"`
	Update Update `cmd:"" aliases:"u,up" help:"update the name of a module"`
	Create Create `cmd:"" aliases:"c" help:"create a new module"`
	List   List   `cmd:"" aliases:"l,ls" help:"list all modules"`
}

// TODO: avoid repetition, move both ModError and FnError to separate file/package
type ModError struct {
	Errors struct {
		Detail string `json:"detail"`
	} `json:"errors"`
}

func extractError(err error) error {
	var e ModError
	openApiError, castOk := err.(*openapi.GenericOpenAPIError)
	if !castOk {
		return err
	}

	d := json.NewDecoder(bytes.NewReader(openApiError.Body()))
	d.DisallowUnknownFields()

	if err := d.Decode(&e); err != nil {
		return err
	}
	return errors.New(e.Errors.Detail)
}

func (f *Mod) Help() string {
	return "Manage modules (description TBD)"
}
