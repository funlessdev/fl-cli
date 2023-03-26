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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/funlessdev/fl-cli/pkg/homedir"
	openapi "github.com/funlessdev/fl-client-sdk-go"
)

const (
	DefaultModule = "_"
)

type Client struct {
	client    *http.Client
	Config    Config
	ApiClient *openapi.APIClient
}

type Config struct {
	Host          string
	BaseURL       *url.URL
	SecretKeyBase string // used when deploying the platform, unused when using API
	AdminToken    string
	APIToken      string
}

// NewConfig creates a new funless config, reading the information from the given configPath
func NewConfig(configPath string) (Config, error) {

	config, _, err := homedir.ReadFromConfigDir(configPath)

	outConfig := Config{
		Host:          "http://localhost:4000",
		SecretKeyBase: "",
		AdminToken:    "",
		APIToken:      "",
	}

	if err != nil {
		if os.IsNotExist(err) {
			return outConfig, nil
		} else {
			return Config{}, err
		}
	}
	configReader := bytes.NewReader(config)
	configScanner := bufio.NewScanner(configReader)
	configMap := make(map[string]string)
	for configScanner.Scan() {
		line := configScanner.Text()
		lineParts := strings.Split(line, "=")
		key, value := strings.TrimSpace(lineParts[0]), strings.TrimSpace(lineParts[1])
		configMap[key] = value
	}

	if err = configScanner.Err(); err != nil {
		return Config{}, err
	}

	if host, v := configMap["api_host"]; v {
		outConfig.Host = host
	}
	if secretKeyBase, v := configMap["secret_key_base"]; v {
		outConfig.SecretKeyBase = secretKeyBase
	}
	if adminToken, v := configMap["admin_token"]; v {
		outConfig.AdminToken = adminToken
	}
	if apiToken, v := configMap["api_token"]; v {
		outConfig.APIToken = apiToken
	}

	return outConfig, nil
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

	apiConfig := openapi.NewConfiguration()
	apiConfig.Servers[0].URL = config.Host
	apiClient := openapi.NewAPIClient(apiConfig)

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
