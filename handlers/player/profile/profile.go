package player

import (
	"net/http"
	"provider/models"
	"provider/utils"

	"gorm.io/gorm"
)

func PlayerProfile(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
		clientID, err := utils.GetClientIdByHeader(db, r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		playerID, err := utils.GetIDfromQuery(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
		}

		// QUERY
		var player models.Player
		err = db.
			Where("id = ? AND client_id = ?", playerID, clientID).
			Preload("Wallet").
			First(&player).Error
		if err != nil {
			utils.JSONError(w, "failed to fetch player profile", http.StatusInternalServerError)
			return
		}

		// RESPONSE
		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "player profile fetched",
			Data: models.Player{
				ID:       player.ID,
				Username: player.Username,
				ClientID: clientID,
				Wallet:   player.Wallet,
			},
		}, http.StatusOK)
	}
}
