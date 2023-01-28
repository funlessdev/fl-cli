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

package client

import (
	"fmt"
	"regexp"
)

type InputValidatorHandler interface {
	ValidateName(name, entity string) error
}

type InputValidator struct{}

var _ InputValidatorHandler = &InputValidator{}

func (i *InputValidator) ValidateName(name, entity string) error {
	regex := regexp.MustCompile("^[_a-zA-Z0-9]+$")
	if regex.MatchString(name) {
		return nil
	} else {
		return fmt.Errorf("invalid %s name", entity)
	}
}
