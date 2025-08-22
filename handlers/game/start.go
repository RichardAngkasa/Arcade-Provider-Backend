package game

import (
	"database/sql"
	"net/http"
	"provider/models"
	"provider/utils"
)

type SpinRequest struct {
	PlayerID  int     `json:"player_id"`
	BetAmount float64 `json:"bet_amount"`
	GameID    string  `json:"game_id"`
}

type SpinResponse struct {
	Symbols      map[string]string   `json:"symbols"`
	Type         string              `json:"type"`
	Amount       float64             `json:"amount"`
	PlayerWallet models.PlayerWallet `json:"player_wallet"`
}

func StartSpin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SpinRequest
		err := utils.BodyChecker(r, &req)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		clientKey := r.Header.Get("X-API-Key")
		if clientKey == "" {
			utils.JSONError(w, "api key missing", http.StatusBadRequest)
			return
		}
		clientID, err := utils.GetClientIdByApiKey(db, clientKey)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		err = utils.AmountLessThanZero(req.BetAmount)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = utils.PlayerMustExistUnderClient(db, clientID, req.PlayerID)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		playerWallet, err := utils.PlayerWallet(db, clientID, req.PlayerID)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clientWallet, err := utils.ClientWallet(db, clientID)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = utils.RequestAmountLessThanBalance(req.BetAmount, playerWallet.Balance)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			utils.JSONResponse(w, "failed to start game", http.StatusInternalServerError)
			return
		}
		var txErr error
		defer func() {
			if p := recover(); p != nil {
				tx.Rollback()
				panic(p)
			} else if txErr != nil {
				tx.Rollback()
			} else {
				err := tx.Commit()
				if err != nil {
					txErr = err
				}
			}
		}()

		var result models.GameResult
		if req.GameID == "pikachu" {
			r, err := RunPikachuGameLogic(req.BetAmount)
			if err != nil {
				utils.JSONError(w, "Spin Error", http.StatusInternalServerError)
				return
			}
			result = r
		}

		gameSessionID, txErr := utils.GameSessionLog(tx, req.PlayerID, clientID, req.BetAmount, result.Amount, req.GameID, result.Type)
		if txErr != nil {
			http.Error(w, txErr.Error(), http.StatusInternalServerError)
			return
		}

		txErr = utils.PlayerWithdraw(tx, clientID, req.PlayerID, req.BetAmount)
		if txErr != nil {
			utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
			return
		}

		txErr = utils.ClientDeposit(tx, clientID, req.BetAmount)
		if txErr != nil {
			utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
			return
		}

		if result.Type == "win" {
			txErr = utils.PlayerDeposit(tx, clientID, req.PlayerID, result.Amount)
			if txErr != nil {
				utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
				return
			}
			txErr = utils.PlayerLogTransaction(tx, playerWallet.ID, req.PlayerID, clientID, sql.NullInt64{Int64: int64(gameSessionID), Valid: true}, result.Amount, "bet_win_player")
			if txErr != nil {
				utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
				return
			}

			txErr = utils.ClientWithdraw(tx, clientID, result.Amount)
			if txErr != nil {
				utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
				return
			}
			txErr = utils.ClientLogTransaction(tx, clientWallet.ID, clientID, result.Amount, "bet_win_player")
			if txErr != nil {
				utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
				return
			}
		}

		if result.Type == "lose" {
			txErr = utils.PlayerLogTransaction(tx, playerWallet.ID, req.PlayerID, clientID, sql.NullInt64{Int64: int64(gameSessionID), Valid: true}, result.Amount, "bet_lose_player")
			if txErr != nil {
				utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
				return
			}

			txErr = utils.ClientLogTransaction(tx, clientWallet.ID, clientID, result.Amount, "bet_lose_player")
			if txErr != nil {
				utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
				return
			}
		}

		updatedPlayerWallet, txErr := utils.PlayerWallet(tx, clientID, req.PlayerID)
		if txErr != nil {
			utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
			return
		}

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "game spined successfully",
			Data: SpinResponse{
				Symbols:      result.Symbols,
				Type:         result.Type,
				Amount:       result.Amount,
				PlayerWallet: updatedPlayerWallet,
			},
		}, http.StatusOK)
	}
}
