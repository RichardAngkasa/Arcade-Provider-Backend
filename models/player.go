package models

import "database/sql"

type Player struct {
	ID       int          `json:"id"`
	Username string       `json:"username"`
	ClientID int          `json:"client_id"`
	Wallet   PlayerWallet `json:"wallet"`
}
type ProfileRequest struct {
	ID int `json:"id"`
}

// WALLET
type PlayerWallet struct {
	ID       int     `json:"id"`
	PlayerID int     `json:"player_id"`
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
type PlayerTransactionsRequest struct {
	PlayerID int `json:"id"`
}
type PlayerWalletTransaction struct {
	ID            int           `json:"id"`
	ClientID      int           `json:"client_id"`
	PlayerID      int           `json:"player_id"`
	WalletID      int           `json:"wallet_id"`
	GameSessionID sql.NullInt64 `json:"game_session_id"`
	Amount        float64       `json:"amount"`
	Type          string        `json:"type"`
}
