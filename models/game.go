package models

type GameResult struct {
	Symbols map[string]string `json:"symbols"`
	Type    string            `json:"type"`
	Amount  float64           `json:"amount"`
}
