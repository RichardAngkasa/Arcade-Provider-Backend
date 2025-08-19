package models

type Client struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	ApiKey   string `json:"api_key"`
}

type ClientWallet struct {
	ID       int     `json:"id"`
	ClientID int     `json:"client_id"`
	Balance  float64 `json:"balance"`
}

type ClientWalletRequest struct {
	ClientID int     `json:"client_id"`
	Amount   float64 `json:"amount"`
}

type ClientWalletTransaction struct {
	ID       int     `json:"id"`
	ClientID int     `json:"client_id"`
	Amount   float64 `json:"amount"`
	Type     string  `json:"type"`
}
