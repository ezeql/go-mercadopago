package mercadopago

import (
	"encoding/json"
	"time"
)

type UserProfileResponse struct {
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

//https://api.mercadopago.com/users/me#options
func (mp *MP) Profile() (*UserProfileResponse, error) {
	r, err := restCall("GET", "/users/me", valuesWithToken(mp))
	if err != nil {
		return nil, err
	}

	u := &UserProfileResponse{}
	b := mustReadAll(r.Body)
	if err = json.Unmarshal(b, u); err != nil {
		return nil, err
	}
	return u, nil
}
