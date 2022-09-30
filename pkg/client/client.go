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
	"errors"
	"fmt"
	"net/http"
	"net/url"

	swagger "github.com/funlessdev/fl-client-sdk-go"
)

const (
	DefaultNamespace = "_"
)

type Client struct {
	client    *http.Client
	Config    Config
	ApiClient *swagger.APIClient
}

type Config struct {
	Host      string
	Namespace string
	BaseURL   *url.URL
}

// NewClient creates a new funless client with the provided http client and configuration.
func NewClient(httpClient *http.Client, config Config) (*Client, error) {
	if len(config.Host) == 0 { // is host missing?
		return nil, errors.New("unable to create new client, missing API host")
	}

	if config.BaseURL == nil { // if BaseURL missing, create it
		u, err := buildBaseURL(config.Host)
		if err != nil {
			return nil, err
		}
		config.BaseURL = u
	}

	apiConfig := swagger.NewConfiguration()
	apiConfig.BasePath = config.Host
	apiClient := swagger.NewAPIClient(apiConfig)

	return &Client{client: httpClient, Config: config, ApiClient: apiClient}, nil
}

func buildBaseURL(host string) (*url.URL, error) {
	baseURL, err := url.Parse(host)

	if err != nil || len(baseURL.Scheme) == 0 || len(baseURL.Host) == 0 {
		baseURL, err = url.Parse("https://" + host)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to create new client because the api host %s is invalid", host)
	}

	return baseURL, nil
}
