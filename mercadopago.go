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
	"time"
)

type MP struct {
	Sandbox      bool
	clientID     string
	clientSecret string
	accessData   Tokens
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

// https://www.mercadopago.com.ar/developers/en/api-docs/account/balance/
func (mp *MP) Balance(userID int) (*Balance, error) {
	v := valuesWithToken(mp)
	uri := fmt.Sprintf("/users/%v/mercadopago_account/balance", userID)

	r, err := restCall("GET", uri, v)
	if err != nil {
		return nil, err
	}

	balance := &Balance{}
	b := mustReadAll(r.Body)

	if err = json.Unmarshal(b, balance); err != nil {
		return nil, err
	}
	return balance, nil
}

//https://api.mercadopago.com/users/me#options
func (mp *MP) Profile() (*User, error) {
	r, err := restCall("GET", "/users/me", valuesWithToken(mp))
	if err != nil {
		return nil, err
	}

	u := &User{}
	b := mustReadAll(r.Body)
	if err = json.Unmarshal(b, u); err != nil {
		return nil, err
	}
	return u, nil
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

const baseURL = "https://api.mercadopago.com"

type Tokens struct {
	AccessToken string `json:"access_token"`
	FreshToken  string `json:"refresh_token"`
	LiveMode    bool   `json:"live_mode"`
}

type Site string

const (
	Argentina Site = "MLA"
	Brazil    Site = "MLB"
	Mexico    Site = "MLM"
	Venezuela Site = "MLV"
	Colombia  Site = "MCO"
	Chile     Site = "MLC"
)

type User struct {
	Address          struct{}    `json:"address"`
	AlternativePhone struct{}    `json:"alternative_phone"`
	BuyerReputation  struct{}    `json:"buyer_reputation"`
	CountryID        string      `json:"country_id"`
	Credit           struct{}    `json:"credit"`
	Email            string      `json:"email"`
	FirstName        string      `json:"first_name"`
	ID               int         `json:"id"`
	Identification   struct{}    `json:"identification"`
	LastName         string      `json:"last_name"`
	Logo             interface{} `json:"logo"`
	Nickname         string      `json:"nickname"`
	Permalink        string      `json:"permalink"`
	Phone            struct{}    `json:"phone"`
	Points           int         `json:"points"`
	RegistrationDate time.Time   `json:"registration_date"`
	SellerExperience string      `json:"seller_experience"`
	SellerReputation struct {
		LevelID           string      `json:"level_id"`
		PowerSellerStatus interface{} `json:"power_seller_status"`
		Transactions      struct{}    `json:"transactions"`
	} `json:"seller_reputation"`
	ShippingModes []interface{} `json:"shipping_modes"`
	SiteID        Site          `json:"site_id"`
	Status        struct{}      `json:"status"`
	Tags          []interface{} `json:"tags"`
	UserType      string        `json:"user_type"`
}

type Balance struct {
	AvailableBalance                  float64 `json:"available_balance"`
	AvailableBalanceByTransactionType []struct {
		Amount          float64 `json:"amount"`
		TransactionType string  `json:"transaction_type"`
	} `json:"available_balance_by_transaction_type"`
	CurrencyID                 string  `json:"currency_id"`
	TotalAmount                float64 `json:"total_amount"`
	UnavailableBalance         float64 `json:"unavailable_balance"`
	UnavailableBalanceByReason []struct {
		Amount float64 `json:"amount"`
		Reason string  `json:"reason"`
	} `json:"unavailable_balance_by_reason"`
	UserID int `json:"user_id"`
}
