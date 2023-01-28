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
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestValidateFnModName(t *testing.T) {

	validator := &InputValidator{}
	validNames := []string{"_", "test_name", "test", "testname", "some____name", "__some__name__", "__init__", "_name_", "_name", "name_", "someName", "SOMENAME", "SomeName", "123somename", "123456"}
	invalidNames := []string{"", "test-name", "test name", "@@", "/somename", "\\somename", "[]name", "{}name", "##name", "èname", "àname", "ùname", "ç", "°", "§", "”", "some:name", "some.name", "some,name", "some;name"}
	entities := []string{"name", "mod", "package", "something"}

	t.Run("should return nil when given a valid name as input", func(t *testing.T) {
		var err error
		for _, n := range validNames {
			err = validator.ValidateName(n, "")
			require.NoError(t, err)
		}
	})

	t.Run("should return an error when given an invalid name as input", func(t *testing.T) {
		var err error
		for _, n := range invalidNames {
			err = validator.ValidateName(n, "")
			require.Error(t, err)
		}
	})

	t.Run("should return an error containing the entity name when given an invalid string as input", func(t *testing.T) {
		var err error
		for _, n := range entities {
			err = validator.ValidateName(invalidNames[0], n)
			require.Error(t, err)
			assert.Equal(t, err.Error(), fmt.Errorf("invalid %s name", n).Error())
		}
	})
}
