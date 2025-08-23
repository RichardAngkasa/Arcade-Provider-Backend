package wallet

import (
	"net/http"
	"provider/middleware"
	"provider/models"
	"provider/utils"

	"gorm.io/gorm"
)

func ClientTransactions(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
		clientID, err := middleware.MustClientID(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// QUERY
		var transactions []models.ClientWalletTransaction
		err = db.
			Where("client_id = ?", clientID).
			Order("created_at DESC").
			Find(&transactions).Error
		if err != nil {
			utils.JSONError(w, "failed to fetch client transactions", http.StatusInternalServerError)
			return
		}

		// RESPONSE
		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "client transactions called successfully",
			Data:    transactions,
		}, http.StatusOK)
	}
}
