package models

type GameResult struct {
	Symbols map[string]string `json:"symbols"`
	Type    string            `gorm:"type:text; not null"`
	Amount  float64           `json:"amount"`
}

type SpinRequest struct {
	PlayerID  int     `json:"player_id"`
	BetAmount float64 `json:"bet_amount"`
	GameID    string  `json:"game_id"`
}

type SpinResponse struct {
	Symbols      map[string]string `json:"symbols"`
	Type         string            `json:"type" gorm:"type:text; not null"`
	Amount       float64           `json:"amount"`
	PlayerWallet PlayerWallet      `json:"player_wallet"`
}

type GameSession struct {
	ID           int     `json:"id"`
	ClientID     int     `json:"client_id"`
	PlayerID     int     `json:"player_id"`
	GameID       string  `json:"game_id"`
	BetAmount    float64 `json:"bet_amount"`
	ResultAmount float64 `json:"result_amount"`
	Type         string  `gorm:"type:text; not null"`
}
