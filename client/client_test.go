// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package client

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("should create a client when config is valid", func(t *testing.T) {
		host := "test-host.com"
		expectedBaseUrl := "https://" + host + "/api"

		config := Config{Host: host}
		client, err := NewClient(http.DefaultClient, config)

		require.NoError(t, err)
		assert.Equal(t, host, client.Config.Host)
		assert.Equal(t, expectedBaseUrl, client.Config.BaseURL.String())
	})

	t.Run("should fail when missing api host in config", func(t *testing.T) {
		config := Config{}
		client, err := NewClient(http.DefaultClient, config)

		require.Error(t, err)
		require.Contains(t, err.Error(), "unable to create new client, missing API host")
		require.Nil(t, client)
	})

	t.Run("should fail to create base URL when api host is invalid", func(t *testing.T) {
		host := "ht://some bad .url/_20_%+off_60000_"

		config := Config{Host: host}
		client, err := NewClient(http.DefaultClient, config)

		require.Error(t, err)
		require.Contains(t, err.Error(), "unable to create new client because the api host "+host+" is invalid")
		require.Nil(t, client)
	})

	t.Run("should create client when host and base url are of different values", func(t *testing.T) {
		host := "test-host.com"
		baseUrl, _ := url.Parse("https://another-url.com")

		config := Config{Host: host, BaseURL: baseUrl}
		client, err := NewClient(http.DefaultClient, config)

		require.NoError(t, err)
		require.Equal(t, host, client.Config.Host)
		require.Equal(t, baseUrl, client.Config.BaseURL)
	})
}

func Test_buildRequestURL(t *testing.T) {
	client, _ := NewClient(nil, Config{Host: "test-host.com"})
	uStr := "fn/hello"

	t.Run("should create {baseURL}/{DefaultNamespace}/{urlStr} when namespace is missing", func(t *testing.T) {
		u, err := client.buildRequestURL(uStr)
		require.NoError(t, err)
		require.Equal(t, client.Config.BaseURL.String()+"/"+DefaultNamespace+"/"+uStr, u.String())
	})

	t.Run("should create url with custom namespace when not missing", func(t *testing.T) {
		client.Config.Namespace = "some-ns"
		u, err := client.buildRequestURL(uStr)
		require.NoError(t, err)
		require.Equal(t, client.Config.BaseURL.String()+"/some-ns/"+uStr, u.String())
	})

	t.Run("should fail with an invalid endpoint", func(t *testing.T) {
		uStr = "_20_%+off_6"
		_, err := client.buildRequestURL(uStr)
		require.Error(t, err)
		require.Equal(t, err.Error(), "invalid endpoint given "+uStr)
	})
}
