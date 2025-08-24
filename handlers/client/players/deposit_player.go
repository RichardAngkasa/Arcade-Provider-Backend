package players

import (
	"database/sql"
	"net/http"
	"provider/middleware"
	"provider/models"
	"provider/utils"

	"gorm.io/gorm"
)

func ClientPlayerDeposit(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
		clientID, err := middleware.MustClientID(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var req models.PlayerWalletRequest
		err = utils.BodyChecker(r, &req)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = utils.AmountLessThanZero(req.Amount)
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
		var updatedClientWallet models.ClientWallet
		var updatedPlayerWallet models.PlayerWallet
		err = db.Transaction(func(tx *gorm.DB) error {
			// withdraw client wallet
			err := utils.ClientDeductBalance(tx, clientID, req.Amount)
			if err != nil {
				return err
			}
			err = utils.ClientLogTransaction(tx, clientID, req.Amount, "withdraw")
			if err != nil {
				return err
			}
			// deposit player wallet
			err = utils.PlayerAddBalance(tx, clientID, req.PlayerID, req.Amount)
			if err != nil {
				return err
			}
			err = utils.PlayerLogTransaction(tx, clientID, req.PlayerID, sql.NullInt64{Valid: false}, req.Amount, "deposit")
			if err != nil {
				return err
			}
			// update new wallet
			updatedClientWallet, err = utils.ClientWallet(tx, clientID)
			if err != nil {
				return err
			}
			updatedPlayerWallet, err = utils.PlayerWallet(tx, clientID, req.PlayerID)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// RESPONSE
		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "player deposit successfully",
			Data: models.PlayerWalletResponse{
				ClientWallet: updatedClientWallet,
				Amount:       req.Amount,
				PlayerWallet: updatedPlayerWallet,
			},
		}, http.StatusOK)
	}
}
