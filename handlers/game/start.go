package game

import (
	"database/sql"
	"net/http"
	"provider/models"
	"provider/utils"

	"gorm.io/gorm"
)

func StartSpin(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
		clientID, err := utils.GetClientIdByHeader(db, r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		var req models.SpinRequest
		err = utils.BodyChecker(r, &req)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
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

		// QUERY
		var updatedPlayerWallet models.PlayerWallet
		var result models.GameResult
		err = db.Transaction(func(tx *gorm.DB) error {
			result, err = RunPikachuGameLogic(req.BetAmount)
			if err != nil {
				return err
			}
			gameSessionID, err := utils.GameLogSession(tx, clientID, req.PlayerID, req.BetAmount, result.Amount, req.GameID, result.Type)
			if err != nil {
				return err
			}
			// player pay to play
			err = utils.PlayerDeductBalance(tx, clientID, req.PlayerID, req.BetAmount)
			if err != nil {
				return err
			}
			err = utils.ClientAddBalance(tx, clientID, req.BetAmount)
			if err != nil {
				return err
			}
			if result.Type == "win" {
				// add balance to player
				err = utils.PlayerAddBalance(tx, clientID, req.PlayerID, result.Amount)
				if err != nil {
					return err
				}
				err = utils.PlayerLogTransaction(tx, clientID, req.PlayerID, sql.NullInt64{Int64: int64(gameSessionID), Valid: true}, result.Amount, "bet_win_player")
				if err != nil {
					return err
				}
				// deduct balance from cleint
				err := utils.ClientDeductBalance(tx, clientID, result.Amount)
				if err != nil {
					return err
				}
				err = utils.ClientLogTransaction(tx, clientID, result.Amount, "bet_win_player")
				if err != nil {
					return err
				}
			}
			if result.Type == "lose" {
				// log player lost
				err = utils.PlayerLogTransaction(tx, clientID, req.PlayerID, sql.NullInt64{Int64: int64(gameSessionID), Valid: true}, result.Amount, "bet_lose_player")
				if err != nil {
					return err
				}
				err = utils.ClientLogTransaction(tx, clientID, result.Amount, "bet_lose_player")
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "game spined successfully",
			Data: models.SpinResponse{
				Symbols:      result.Symbols,
				Type:         result.Type,
				Amount:       result.Amount,
				PlayerWallet: updatedPlayerWallet,
			},
		}, http.StatusOK)
	}
}
