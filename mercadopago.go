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
	Address struct {
		Address string `json:"address"`
		City    string `json:"city"`
		State   string `json:"state"`
		ZipCode string `json:"zip_code"`
	} `json:"address"`
	AlternativePhone struct {
		AreaCode  string `json:"area_code"`
		Extension string `json:"extension"`
		Number    string `json:"number"`
	} `json:"alternative_phone"`
	BuyerReputation struct {
		CanceledTransactions int           `json:"canceled_transactions"`
		Tags                 []interface{} `json:"tags"`
		Transactions         struct {
			Canceled struct {
				Paid  int `json:"paid"`
				Total int `json:"total"`
			} `json:"canceled"`
			Completed   int `json:"completed"`
			NotYetRated struct {
				Paid  int `json:"paid"`
				Total int `json:"total"`
				Units int `json:"units"`
			} `json:"not_yet_rated"`
			Period  string `json:"period"`
			Total   int    `json:"total"`
			Unrated struct {
				Paid  int `json:"paid"`
				Total int `json:"total"`
			} `json:"unrated"`
		} `json:"transactions"`
	} `json:"buyer_reputation"`
	CountryID string `json:"country_id"`
	Credit    struct {
		Consumed      float64 `json:"consumed"`
		CreditLevelID string  `json:"credit_level_id"`
	} `json:"credit"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	ID             int    `json:"id"`
	Identification struct {
		Number string `json:"number"`
		Type   string `json:"type"`
	} `json:"identification"`
	LastName  string      `json:"last_name"`
	Logo      interface{} `json:"logo"`
	Nickname  string      `json:"nickname"`
	Permalink string      `json:"permalink"`
	Phone     struct {
		AreaCode  string `json:"area_code"`
		Extension string `json:"extension"`
		Number    string `json:"number"`
		Verified  bool   `json:"verified"`
	} `json:"phone"`
	Points           int       `json:"points"`
	RegistrationDate time.Time `json:"registration_date"`
	SellerExperience string    `json:"seller_experience"`
	SellerReputation struct {
		LevelID           string      `json:"level_id"`
		PowerSellerStatus interface{} `json:"power_seller_status"`
		Transactions      struct {
			Canceled  int    `json:"canceled"`
			Completed int    `json:"completed"`
			Period    string `json:"period"`
			Ratings   struct {
				Negative float64 `json:"negative"`
				Neutral  int     `json:"neutral"`
				Positive float64 `json:"positive"`
			} `json:"ratings"`
			Total int `json:"total"`
		} `json:"transactions"`
	} `json:"seller_reputation"`
	ShippingModes []string `json:"shipping_modes"`
	SiteID        string   `json:"site_id"`
	Status        struct {
		Billing struct {
			Allow bool          `json:"allow"`
			Codes []interface{} `json:"codes"`
		} `json:"billing"`
		Buy struct {
			Allow            bool          `json:"allow"`
			Codes            []interface{} `json:"codes"`
			ImmediatePayment struct {
				Reasons  []interface{} `json:"reasons"`
				Required bool          `json:"required"`
			} `json:"immediate_payment"`
		} `json:"buy"`
		ConfirmedEmail   bool `json:"confirmed_email"`
		ImmediatePayment bool `json:"immediate_payment"`
		List             struct {
			Allow            bool          `json:"allow"`
			Codes            []interface{} `json:"codes"`
			ImmediatePayment struct {
				Reasons  []interface{} `json:"reasons"`
				Required bool          `json:"required"`
			} `json:"immediate_payment"`
		} `json:"list"`
		Mercadoenvios          string `json:"mercadoenvios"`
		MercadopagoAccountType string `json:"mercadopago_account_type"`
		MercadopagoTcAccepted  bool   `json:"mercadopago_tc_accepted"`
		RequiredAction         string `json:"required_action"`
		Sell                   struct {
			Allow            bool          `json:"allow"`
			Codes            []interface{} `json:"codes"`
			ImmediatePayment struct {
				Reasons  []interface{} `json:"reasons"`
				Required bool          `json:"required"`
			} `json:"immediate_payment"`
		} `json:"sell"`
		SiteStatus string `json:"site_status"`
		UserType   string `json:"user_type"`
	} `json:"status"`
	Tags     []string `json:"tags"`
	UserType string   `json:"user_type"`
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
