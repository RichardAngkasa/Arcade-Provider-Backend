package players

import (
	"net/http"
	"provider/middleware"
	"provider/models"
	"provider/utils"

	"gorm.io/gorm"
)

func AdminPlayerProfile(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
		_, err := middleware.MustAdminID(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}
		playerID, err := utils.GetIDfromQuery(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
		}

		// QUERY
		var player models.Player
		err = db.
			Preload("Wallet").
			First(&player, playerID).Error
		if err != nil {
			utils.JSONError(w, "failed to fetch player profile", http.StatusInternalServerError)
			return
		}

		// RESPONSE
		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "player profile fetched",
			Data:    player,
		}, http.StatusOK)
	}
}
