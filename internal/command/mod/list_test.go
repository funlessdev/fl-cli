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
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/funlessdev/fl-cli/pkg"
	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/funlessdev/fl-cli/test/mocks"
	openapi "github.com/funlessdev/fl-client-sdk-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestModList(t *testing.T) {
	mod1 := "mod1"
	mod2 := "mod2"
	mod3 := "mod3"
	testMods := []string{mod1, mod2, mod3}
	testCtx := context.Background()
	testLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).Build()

	t.Run("should use ModService.List to list modules", func(t *testing.T) {
		cmd := List{}

		mockModHandler := mocks.NewModHandler(t)
		mockModHandler.On("List", mock.Anything).Return(pkg.ModuleNameList{Names: testMods}, nil)

		err := cmd.Run(testCtx, mockModHandler, testLogger, &Mod{})
		require.NoError(t, err)
		mockModHandler.AssertCalled(t, "List", mock.Anything)
		mockModHandler.AssertNumberOfCalls(t, "List", 1)
		mockModHandler.AssertExpectations(t)
	})

	t.Run("should correctly print result", func(t *testing.T) {
		cmd := List{}

		mockModHandler := mocks.NewModHandler(t)
		mockModHandler.On("List", mock.Anything).Return(pkg.ModuleNameList{Names: testMods}, nil)

		var outbuf bytes.Buffer
		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).Build()

		err := cmd.Run(testCtx, mockModHandler, bufLogger, &Mod{})

		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%s\n%s\n%s\n", testMods[0], testMods[1], testMods[2]), (&outbuf).String())
		mockModHandler.AssertExpectations(t)
	})

	t.Run("should correctly print result with count", func(t *testing.T) {
		cmd := List{
			Count: true,
		}

		mockModHandler := mocks.NewModHandler(t)
		mockModHandler.On("List", mock.Anything).Return(pkg.ModuleNameList{Names: testMods}, nil)

		var outbuf bytes.Buffer
		bufLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).Build()

		err := cmd.Run(testCtx, mockModHandler, bufLogger, &Mod{})

		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%s\n%s\n%s\nCount: %d\n", testMods[0], testMods[1], testMods[2], len(testMods)), (&outbuf).String())
		mockModHandler.AssertExpectations(t)
	})

	t.Run("should return error if invalid list request", func(t *testing.T) {
		cmd := List{}

		mockModHandler := mocks.NewModHandler(t)

		e := &openapi.GenericOpenAPIError{}
		mockModHandler.On("List", mock.Anything).Return(pkg.ModuleNameList{}, e)

		err := cmd.Run(testCtx, mockModHandler, testLogger, &Mod{})
		require.Error(t, err)
	})
}
