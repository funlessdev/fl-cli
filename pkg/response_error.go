// Copyright 2023 Giuseppe De Palma, Matteo Trentin
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

import (
	"bytes"
	"encoding/json"
	"errors"

	openapi "github.com/funlessdev/fl-client-sdk-go"
)

type FLError struct {
	Errors struct {
		Detail string `json:"detail"`
	} `json:"errors"`
}

func ExtractError(err error) error {
	var e FLError
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
