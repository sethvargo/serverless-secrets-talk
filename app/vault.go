// Copyright 2019 Seth Vargo
// Copyright 2019 Google, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var vaultAddr = os.Getenv("VAULT_ADDR")

func vaultAccess(name string) (string, error) {
	if vaultAddr == "" {
		return "", fmt.Errorf("missing VAULT_ADDR")
	}

	vaultToken, err := vaultLogin()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodGet, vaultAddr+"/v1/"+name, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create access request: %s", err)
	}
	req.Header.Add("x-vault-token", vaultToken)

	client := insecureHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute access: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		var b bytes.Buffer
		if _, err := io.Copy(&b, resp.Body); err != nil {
			return "", err
		}
		return "", fmt.Errorf("bad response from access: %s", b.String())
	}

	s := struct {
		Data struct {
			Data struct {
				Value string `json:"value"`
			} `json:"data"`
		} `json:"data"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return "", err
	}

	return s.Data.Data.Value, nil
}

func vaultLogin() (string, error) {
	client := new(http.Client)

	// Get JWT token
	req, err := http.NewRequest(http.MethodGet, "http://metadata/computeMetadata/v1/instance/service-accounts/default/identity?audience=http://vault/myapp&format=full", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode > 299 {
		return "", fmt.Errorf("failed to get jwt: %s", body)
	}

	// Login to Vault with JWT token

	// It's a live demo, what could go wrong? In a real production setup, you
	// would want to validate the Vault server certificate.
	client = insecureHTTPClient()
	j := bytes.NewBufferString(`{"role":"myapp", "jwt":"` + string(body) + `"}`)

	req, err = http.NewRequest(http.MethodPost, vaultAddr+"/v1/auth/gcp-serverless/login", j)
	if err != nil {
		return "", fmt.Errorf("failed to make vault login request: %s", err)
	}

	resp, err = client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute vault login: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		var b bytes.Buffer
		if _, err := io.Copy(&b, resp.Body); err != nil {
			return "", err
		}
		return "", fmt.Errorf("bad response from vault login: %s", b.String())
	}

	s := struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return "", fmt.Errorf("failed to decode vault login response: %s", err)
	}

	return s.Auth.ClientToken, nil
}

func insecureHTTPClient() *http.Client {
	defTransport := http.DefaultTransport.(*http.Transport)
	return &http.Client{
		Transport: &http.Transport{
			Proxy:                 defTransport.Proxy,
			DialContext:           defTransport.DialContext,
			MaxIdleConns:          defTransport.MaxIdleConns,
			IdleConnTimeout:       defTransport.IdleConnTimeout,
			ExpectContinueTimeout: defTransport.ExpectContinueTimeout,
			TLSHandshakeTimeout:   defTransport.TLSHandshakeTimeout,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}
}
