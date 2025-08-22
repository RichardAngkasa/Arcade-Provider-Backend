package players

import (
	"database/sql"
	"net/http"
	"provider/middleware"
	"provider/models"
	"provider/utils"
)

func AdminPlayerProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := middleware.MustAdminID(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var req struct {
			ID int `json:"id"`
		}
		err = utils.BodyChecker(r, &req)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		var player models.Player
		err = db.QueryRow(`
			SELECT id, username, client_id
			FROM players
			WHERE id = $1 
		`, req.ID).Scan(&player.ID, &player.Username, &player.ClientID)
		if err != nil {
			utils.JSONError(w, "failed to fetch player profile", http.StatusInternalServerError)
			return
		}

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "profile fetched",
			Data:    player,
		}, http.StatusOK)
	}
}
