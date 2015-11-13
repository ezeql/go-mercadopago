// Copyright 2015 Ezequiel Moreno
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package mercadopago

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type MP struct {
	Sandbox      bool
	clientID     string
	clientSecret string
	accessData   Tokens
}

const baseURL = "https://api.mercadopago.com"

type Site string

const (
	Argentina Site = "MLA"
	Brazil    Site = "MLB"
	Mexico    Site = "MLM"
	Venezuela Site = "MLV"
	Colombia  Site = "MCO"
	Chile     Site = "MLC"
)

type Paging struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type Tokens struct {
	AccessToken string `json:"access_token"`
	FreshToken  string `json:"refresh_token"`
	LiveMode    bool   `json:"live_mode"`
}

func Build(clientID string, clientSecret string) *MP {
	mp := &MP{}
	mp.clientID = clientID
	mp.clientSecret = clientSecret
	mp.Sandbox = true
	return mp
}

//https://www.mercadopago.com.ar/developers/en/api-docs/custom-checkout/authentication/
func (mp *MP) AccessToken() (*Tokens, error) {
	data := url.Values{}
	data.Set("client_id", mp.clientID)
	data.Add("client_secret", mp.clientSecret)
	data.Add("grant_type", "client_credentials")

	r, err := restCall("POST", "/oauth/token", data)
	if err != nil {
		return nil, err
	}

	b := mustReadAll(r.Body)

	if err = json.Unmarshal(b, &mp.accessData); err != nil {
		return nil, err
	}

	return &mp.accessData, nil
}

func valuesWithToken(mp *MP) url.Values {
	v := url.Values{}
	v.Add("access_token", mp.accessData.AccessToken)
	return v
}

func restCall(method string, resource string, values url.Values) (*http.Response, error) {
	client := &http.Client{}
	apiURL := baseURL

	u, err := url.ParseRequestURI(apiURL)
	if err != nil {
		return nil, err
	}

	u.Path = resource
	urlStr := fmt.Sprintf("%v", u)
	r, _ := http.NewRequest(method, urlStr, bytes.NewBufferString(values.Encode()))
	r.Header.Add("Content-Length", strconv.Itoa(len(values.Encode())))
	r.Header.Add("content-type:", "application/x-www-form-urlencoded")

	return client.Do(r)
}

func mustReadAll(r io.Reader) []byte {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic("arrrgh")
	}
	return b
}
