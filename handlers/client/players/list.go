package players

import (
	"net/http"
	"provider/middleware"
	"provider/models"
	"provider/utils"

	"gorm.io/gorm"
)

func ClientPlayers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
		clientID, err := middleware.MustClientID(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// QUERY
		var players []models.Player
		err = db.
			Where("client_id = ?", clientID).
			Order("created_at DESC").
			Find(&players).Error
		if err != nil {
			utils.JSONError(w, "failed to fetch client players", http.StatusInternalServerError)
			return
		}

		// RESPONSE
		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "client players called successfully",
			Data:    players,
		}, http.StatusOK)
	}
}
