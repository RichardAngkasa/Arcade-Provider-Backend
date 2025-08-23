package clients

import (
	"net/http"
	"provider/middleware"
	"provider/models"
	utils "provider/utils"

	"gorm.io/gorm"
)

func ClientProfile(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
		clientID, err := middleware.MustClientID(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
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
			Data: models.Client{
				ID:       client.ID,
				Username: client.Username,
				Email:    client.Email,
				ApiKey:   client.ApiKey,
				Wallet:   client.Wallet,
			},
		}, http.StatusOK)
	}
}
