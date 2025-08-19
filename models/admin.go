package models

type AdminLoginRequest struct {
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}

type AdminWalletTransaction struct {
	ID       int     `json:"id"`
	ClientID int     `json:"client_id"`
	Amount   float64 `json:"amount"`
	Type     string  `json:"type"`
}
