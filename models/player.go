package models

import "database/sql"

type Player struct {
	ID       int          `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string       `json:"username"`
	ClientID int          `json:"client_id"`
	Wallet   PlayerWallet `gorm:"foreignKey:PlayerID" json:"wallet"`
}

// WALLET
type PlayerWallet struct {
	ID       int     `gorm:"primaryKey;autoIncrement" json:"id"`
	PlayerID int     `gorm:"not null" json:"player_id"`
	ClientID int     `json:"client_id"`
	Balance  float64 `json:"balance"`
}
type PlayerWalletRequest struct {
	Amount   float64 `json:"amount"`
	PlayerID int     `json:"id"`
}
type PlayerWalletResponse struct {
	ClientWallet ClientWallet `json:"client_wallet"`
	Amount       float64      `json:"amount"`
	PlayerWallet PlayerWallet `json:"player_wallet"`
}
type PlayerWalletTransactions struct {
	ID            int           `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID      int           `json:"client_id"`
	PlayerID      int           `json:"player_id"`
	WalletID      int           `gorm:"column:player_wallet_id" json:"player_wallet_id"`
	GameSessionID sql.NullInt64 `json:"game_session_id"`
	Amount        float64       `json:"amount"`
	Type          string        `gorm:"type:text; not null"`
}
