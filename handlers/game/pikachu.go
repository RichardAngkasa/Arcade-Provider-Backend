package game

import (
	"errors"
	"math/rand"
	"time"

	"provider/models"
)

func RunPikachuGameLogic(betAmount float64) (models.GameResult, error) {
	if betAmount <= 0 {
		return models.GameResult{}, errors.New("invalid bet amount")
	}

	rand.Seed(time.Now().UnixNano())

	symbolPool := []string{"Pikachu", "Bulbasaur", "Charmander"}

	resultSymbols := map[string]string{
		"slot1": symbolPool[rand.Intn(len(symbolPool))],
		"slot2": symbolPool[rand.Intn(len(symbolPool))],
		"slot3": symbolPool[rand.Intn(len(symbolPool))],
	}

	isWin := resultSymbols["slot1"] == resultSymbols["slot2"] && resultSymbols["slot2"] == resultSymbols["slot3"]
	var payout float64
	var resultType string

	if isWin {
		payoutMultiplier := 2.0
		payout = betAmount * payoutMultiplier
		resultType = "win"
	} else {
		payout = 0
		resultType = "lose"
	}

	return models.GameResult{
		Symbols: resultSymbols,
		Type:    resultType,
		Amount:  payout,
	}, nil
}
