package models

type Client struct {
	ID       int          `json:"id"`
	Username string       `json:"username"`
	Email    string       `json:"email"`
	Password string       `json:"password"`
	ApiKey   string       `json:"api_key"`
	Wallet   ClientWallet `gorm:"foreignKey:ClientID"`
}

// WALLET
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
	ID             int     `json:"id"`
	ClientID       int     `json:"client_id"`
	Amount         float64 `json:"amount"`
	Type           string  `gorm:"type:text; not null"`
	ClientWalletID int     `json:"client_wallet_id"`
}

// AUTH
type ClientRegisterResponse struct {
	APIKey string `json:"api_key"`
}
type ClientLoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}
type ClientLoginResponse struct {
	APIKey string `json:"api_key"`
	Token  string `json:"token"`
}
