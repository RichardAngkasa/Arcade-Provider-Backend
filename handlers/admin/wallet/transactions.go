package wallet

import (
	"net/http"
	"provider/middleware"
	"provider/models"
	"provider/utils"

	"gorm.io/gorm"
)

func AdminTransactions(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
		_, err := middleware.MustAdminID(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// QUERY
		var transactions []models.AdminWalletTransaction
		err = db.
			Order("created_at DESC").
			Find(&transactions).Error
		if err != nil {
			utils.JSONError(w, "failed to fetch admin transactions", http.StatusInternalServerError)
			return
		}

		// RESPONSE
		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "admin transactions called successfully",
			Data:    transactions,
		}, http.StatusOK)
	}
}
