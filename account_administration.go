package mercadopago

import (
	"encoding/json"
	"fmt"
)

type BalanceResponse struct {
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

// type BalanceHistoryResponse struct {
// }

type MovementsResponse struct {
	Paging  Paging `json:"paging"`
	Results []struct {
		Amount          float64       `json:"amount"`
		BalancedAmount  float64       `json:"balanced_amount"`
		ClientID        int           `json:"client_id"`
		CurrencyID      string        `json:"currency_id"`
		DateCreated     string        `json:"date_created"`
		DateReleased    string        `json:"date_released"`
		Detail          string        `json:"detail"`
		FinancialEntity string        `json:"financial_entity"`
		ID              int           `json:"id"`
		Label           []interface{} `json:"label"`
		OriginalMoveID  interface{}   `json:"original_move_id"`
		ReferenceID     int           `json:"reference_id"`
		SiteID          Site          `json:"site_id"`
		Status          string        `json:"status"`
		Type            string        `json:"type"`
		UserID          int           `json:"user_id"`
	} `json:"results"`
}

//Balance returns users balance for the specified userID
func (mp *MP) Balance(userID int) (*BalanceResponse, error) {
	v := valuesWithToken(mp)
	uri := fmt.Sprintf("/users/%v/mercadopago_account/balance", userID)

	r, err := restCall("GET", uri, v)
	if err != nil {
		return nil, err
	}

	balance := &BalanceResponse{}
	b := mustReadAll(r.Body)

	if err = json.Unmarshal(b, balance); err != nil {
		return nil, err
	}
	return balance, nil
}

// https://www.mercadopago.com.ar/developers/en/api-docs/account/balance/
// func (mp *MP) BalanceHistory(userID int) (*Balance, error) {
// }

//Movements allows you to obtain information about single or multiple movements by applying different filters.
func (mp *MP) Movements() (*MovementsResponse, error) {
	v := valuesWithToken(mp)
	//TODO: FILTERS
	uri := "/mercadopago_account/movements/search"

	r, err := restCall("GET", uri, v)
	if err != nil {
		return nil, err
	}

	movements := &MovementsResponse{}
	b := mustReadAll(r.Body)

	if err = json.Unmarshal(b, movements); err != nil {
		return nil, err
	}
	return movements, nil
}
