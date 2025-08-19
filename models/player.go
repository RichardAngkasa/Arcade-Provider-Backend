package models

type Player struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	ClientID string `json:"client_id"`
}

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
	ClientWallet  ClientWallet `json:"client_wallet"`
	DepositAmount float64      `json:"amount"`
	PlayerID      int          `json:"player_id"`
	PlayerWallet  PlayerWallet `json:"player_Wallet"`
}

type PlayerTransactionsRequest struct {
	PlayerID int `json:"id"`
}

type PlayerWalletTransaction struct {
	ID       int     `json:"id"`
	ClientID int     `json:"client_id"`
	PlayerID int     `json:"player_id"`
	WalletID int     `json:"wallet_id"`
	Amount   float64 `json:"amount"`
	Type     string  `json:"type"`
}
