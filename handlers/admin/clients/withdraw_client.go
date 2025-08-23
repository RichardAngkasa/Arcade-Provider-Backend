package clients

import (
	"net/http"
	"provider/middleware"
	"provider/models"
	"provider/utils"

	"gorm.io/gorm"
)

func AdminClientWithdraw(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
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

		// QUERY
		var updatedClientWallet models.ClientWallet
		err = db.Transaction(func(tx *gorm.DB) error {
			err := utils.ClientDeductBalance(tx, req.ClientID, req.Amount)
			if err != nil {
				return err
			}
			err = utils.ClientLogTransaction(tx, req.ClientID, req.Amount, "withdraw")
			if err != nil {
				return err
			}
			err = utils.AdminLogTransaction(tx, req.ClientID, req.Amount, "deposit")
			if err != nil {
				return err
			}
			updatedClientWallet, err = utils.ClientWallet(tx, req.ClientID)
			if err != nil {
				return err
			}
			return nil
		})

		// RESPONSE
		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "client withdraw successfully",
			Data:    updatedClientWallet,
		}, http.StatusOK)
	}
}
