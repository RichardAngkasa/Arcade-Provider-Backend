package wallet

import (
	"database/sql"
	"net/http"
	"provider/middleware"
	"provider/models"
	"provider/utils"
)

func PlayerWithdraw(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		clientWallet, err := utils.ClientWallet(db, clientID)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		playerWallet, err := utils.PlayerWallet(db, clientID, req.PlayerID)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = utils.RequestAmountLessThanBalance(req.Amount, playerWallet.Balance)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			utils.JSONResponse(w, "failed to start transaction", http.StatusInternalServerError)
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

		txErr = utils.PlayerWithdraw(tx, clientID, req.PlayerID, req.Amount)
		if txErr != nil {
			utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
			return
		}
		txErr = utils.PlayerLogTransaction(tx, playerWallet.ID, req.PlayerID, clientID, sql.NullInt64{Valid: false}, req.Amount, "withdraw")
		if txErr != nil {
			utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
			return
		}

		txErr = utils.ClientDeposit(tx, clientID, req.Amount)
		if txErr != nil {
			utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
			return
		}
		txErr = utils.ClientLogTransaction(tx, clientWallet.ID, clientID, req.Amount, "deposit")
		if txErr != nil {
			utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
			return
		}

		updatedClientWallet, txErr := utils.ClientWallet(tx, clientID)
		if err != nil {
			utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
			return
		}
		updatedPlayerWallet, txErr := utils.PlayerWallet(tx, clientID, req.PlayerID)
		if err != nil {
			utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
			return
		}

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "player withdraw successfully",
			Data: models.PlayerWalletResponse{
				ClientWallet:  updatedClientWallet,
				DepositAmount: req.Amount,
				PlayerWallet:  updatedPlayerWallet,
			},
		}, http.StatusOK)

	}
}
