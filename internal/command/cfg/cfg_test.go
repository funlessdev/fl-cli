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

package cfg

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/funlessdev/fl-cli/pkg/client"
	"github.com/funlessdev/fl-cli/pkg/homedir"
	"github.com/funlessdev/fl-cli/pkg/log"
	"github.com/stretchr/testify/require"
)

func TestCfg(t *testing.T) {
	var outbuf bytes.Buffer
	testCtx := context.Background()
	testSetLogger, _ := log.NewLoggerBuilder().WithWriter(os.Stdout).Build()
	testGetLogger, _ := log.NewLoggerBuilder().WithWriter(&outbuf).DisableAnimation().Build()

	homedirPath, err := os.MkdirTemp("", "funless-test-cfg-")
	require.NoError(t, err)

	defer func() {
		homedir.GetHomeDir = os.UserHomeDir
		os.RemoveAll(homedirPath)
	}()

	config, err := client.NewConfig("config")
	require.NoError(t, err)

	setCmds := [4]CfgSet{
		{
			Key:   "host",
			Value: "test_host",
		},
		{
			Key:   "api_token",
			Value: "test_api_token",
		},
		{
			Key:   "admin_token",
			Value: "test_admin_token",
		},
		{
			Key:   "secret_key_base",
			Value: "test_secret_key_base",
		},
	}

	for _, c := range setCmds {
		err = c.Run(testCtx, testSetLogger, config)
		require.NoError(t, err)
	}

	getCmds := [4]CfgGet{
		{
			Key: "host",
		},
		{
			Key: "api_token",
		},
		{
			Key: "admin_token",
		},
		{
			Key: "secret_key_base",
		},
	}

	config, err = client.NewConfig("config")
	require.NoError(t, err)

	for _, c := range getCmds {
		err = c.Run(testCtx, testGetLogger, config)
		require.NoError(t, err)
	}

	expected :=
		"host=test_host\n" +
			"api_token=test_api_token\n" +
			"admin_token=test_admin_token\n" +
			"secret_key_base=test_secret_key_base\n"

	require.Equal(t, expected, outbuf.String())
}
