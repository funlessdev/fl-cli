// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	DefaultNamespace = "_"
)

type Client struct {
	client *http.Client
	Config Config
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

	return &Client{client: httpClient, Config: config}, nil
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

func (c *Client) buildRequestURL(endPoint string) (*url.URL, error) {
	ns := DefaultNamespace
	if len(c.Config.Namespace) != 0 { // If namespace not missing
		ns = c.Config.Namespace
	}
	ep := fmt.Sprintf("%s/%s/%s", c.Config.BaseURL.String(), ns, endPoint)

	u, err := url.Parse(ep)
	if err != nil {
		// todo Debug "url.Parse(%s) error: %s\n", urlStr, err
		return nil, fmt.Errorf("invalid endpoint given %s", endPoint)
	}
	return u, nil
}

func (c *Client) CreateGet(urlStr string) (*http.Request, error) {
	u, err := c.buildRequestURL(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		// Debug(DbgError, "http.NewRequest(%v, %s, buf) error: %s\n", method, u.String(), err)
		// errStr := wski18n.T("Error initializing request: {{.err}}", map[string]interface{}{"err": err})
		// werr := MakeWskError(errors.New(errStr), EXIT_CODE_ERR_GENERAL, DISPLAY_MSG, NO_DISPLAY_USAGE)
		return nil, errors.New("error initializing request")
	}

	return req, nil
}

func (c *Client) Send(request *http.Request) (*http.Response, error) {
	// Issue the request to the funless server endpoint
	res, err := c.client.Do(request)
	if err != nil {
		// Debug(DbgError, "HTTP Do() [req %s] error: %s\n", req.URL.String(), err)
		// werr := MakeWskError(err, EXIT_CODE_ERR_NETWORK, DISPLAY_MSG, NO_DISPLAY_USAGE)
		return nil, fmt.Errorf("error sending request %s", request.URL.String())
	}

	return res, nil
}

// func (c *Client) CreatePostRequest(u url.URL, body interface{}) (*http.Request, error) {
// 	var buf io.ReadWriter
// 	if body != nil {
// 		buf = new(bytes.Buffer)
// 		encoder := json.NewEncoder(buf)
// 		encoder.SetEscapeHTML(false)
// 		err := encoder.Encode(body)

// 		if err != nil {
// 			// Debug(DbgError, "json.Encode(%#v) error: %s\n", body, err)
// 			// errStr := wski18n.T("Error encoding request body: {{.err}}", map[string]interface{}{"err": err})
// 			// werr := MakeWskError(errors.New(errStr), EXIT_CODE_ERR_GENERAL, DISPLAY_MSG, NO_DISPLAY_USAGE)
// 			return nil, errors.New("error encoding request body")
// 		}
// 	}

// 	req, err := http.NewRequest(http.MethodPost, u.String(), buf)
// 	if err != nil {
// 		// Debug(DbgError, "http.NewRequest(%v, %s, buf) error: %s\n", method, u.String(), err)
// 		// errStr := wski18n.T("Error initializing request: {{.err}}", map[string]interface{}{"err": err})
// 		// werr := MakeWskError(errors.New(errStr), EXIT_CODE_ERR_GENERAL, DISPLAY_MSG, NO_DISPLAY_USAGE)
// 		return nil, errors.New("error initializing request")
// 	}

// 	req.Header.Add("Content-Type", "application/json")

// 	return req, nil
// }
