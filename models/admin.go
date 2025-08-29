package models

type AdminLoginRequest struct {
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}

type AdminWalletTransactions struct {
	ID       int     `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID int     `json:"client_id"`
	Amount   float64 `json:"amount"`
	Type     string  `gorm:"type:text; not null"`
}
