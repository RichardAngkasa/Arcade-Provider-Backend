package wallet

import (
	"database/sql"
	"net/http"
	"provider/middleware"
	"provider/models"
	"provider/utils"
)

func ClientWithdraw(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := middleware.MustAdminID(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var req models.ClientWalletRequest
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
		err = utils.ClientExistenceByID(db, req.ClientID)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		clientWallet, err := utils.ClientWallet(db, req.ClientID)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = utils.RequestAmountGreaterThanBalance(req.Amount, clientWallet.Balance)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			utils.JSONError(w, "failed to start transaction", http.StatusInternalServerError)
			return
		}

		var txErr error
		defer func() {
			p := recover()
			if p != nil {
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

		txErr = utils.ClientWithdraw(tx, req.ClientID, req.Amount)
		if txErr != nil {
			utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
			return
		}
		txErr = utils.ClientLogTransaction(tx, clientWallet.ID, req.ClientID, req.Amount, "withdraw")
		if txErr != nil {
			utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
			return
		}
		txErr = utils.AdminLogTransaction(tx, req.ClientID, req.Amount, "withdraw")
		if txErr != nil {
			utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
			return
		}

		updatedWallet, err := utils.ClientWallet(tx, req.ClientID)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "client withdraw successfully",
			Data:    updatedWallet,
		}, http.StatusOK)
	}
}
