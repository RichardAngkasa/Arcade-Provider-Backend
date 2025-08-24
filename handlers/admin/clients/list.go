package clients

import (
	"net/http"
	"provider/middleware"
	"provider/models"
	"provider/utils"

	"gorm.io/gorm"
)

func AdminClients(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
		_, err := middleware.MustAdminID(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// QUERY
		var clients []models.Client
		err = db.
			Preload("Wallet").
			Order("created_at DESC").
			Find(&clients).Error
		if err != nil {
			utils.JSONError(w, "failed to fetch clients", http.StatusInternalServerError)
			return
		}

		// RESPONSE
		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "clients called successfully",
			Data:    clients,
		}, http.StatusOK)
	}
}
