package clients

import (
	"net/http"
	"provider/middleware"
	"provider/models"
	"provider/utils"

	"gorm.io/gorm"
)

func AdminClientProfile(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
		_, err := middleware.MustAdminID(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}
		clientID, err := utils.GetIDfromQuery(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
		}

		// QUERY
		var client models.Client
		err = db.
			Preload("Wallet").
			First(&client, clientID).Error
		if err != nil {
			utils.JSONError(w, "failed to fetch client profile", http.StatusInternalServerError)
			return
		}

		// RESPONSE
		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "client profile fetched",
			Data:    client,
		}, http.StatusOK)
	}
}
