package wallet

import (
	"database/sql"
	"net/http"
	"provider/models"
	"provider/utils"
)

func ClientDeposit(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.JSONError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		_, err := utils.GetIDFromToken(r, "jwt_token_admin", "admin")
		if err != nil {
			utils.JSONError(w, "unauthorize", http.StatusUnauthorized)
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

		txErr = utils.ClientDeposit(tx, req.ClientID, req.Amount)
		if txErr != nil {
			utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
			return
		}
		txErr = utils.ClientLogTransaction(tx, clientWallet.ID, req.ClientID, req.Amount, "deposit")
		if txErr != nil {
			utils.JSONError(w, txErr.Error(), http.StatusInternalServerError)
			return
		}
		txErr = utils.AdminLogTransaction(tx, req.ClientID, req.Amount, "deposit")
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
			Message: "client deposit successfully",
			Data:    updatedWallet,
		}, http.StatusOK)
	}
}
