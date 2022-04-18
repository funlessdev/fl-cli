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
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type ClientAPI interface {
	NewRequest(method string, url url.URL, body interface{}, includeNamespaceInUrl bool) (*http.Request, error)
}

type Client struct {
	client *http.Client
	Config Config
}

type Config struct {
	Host    string
	BaseURL *url.URL
}

func (c *Config) isHostMissing() bool {
	return len(c.Host) == 0

}

func (c *Config) prepareBaseURL() error {
	if c.BaseURL == nil {
		baseURL, err := makeBaseURL(c.Host)
		if err != nil {
			return fmt.Errorf("unable to create new client because the api host %s is invalid", c.Host)
		}
		c.BaseURL = baseURL
	}
	return nil
}

// NewClient creates a new funless client with the provided http client and fl configuration.
func NewClient(httpClient *http.Client, config Config) (*Client, error) {
	if config.isHostMissing() {
		return nil, errors.New("unable to create new client, missing API host")
	}

	err := config.prepareBaseURL()
	if err != nil {
		return nil, err
	}

	return &Client{client: httpClient, Config: config}, nil
}

func makeBaseURL(host string) (*url.URL, error) {
	url, err := url.Parse(host + "/api")

	if err != nil || len(url.Scheme) == 0 || len(url.Host) == 0 {
		url, err = url.Parse("https://" + host + "/api")
	}

	return url, err
}
