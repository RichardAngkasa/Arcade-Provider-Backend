package players

import (
	"database/sql"
	"net/http"
	"provider/middleware"
	"provider/models"
	"provider/utils"
)

func AdminPlayers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := middleware.MustAdminID(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
			SELECT id, username, client_id
			FROM players
		`)
		if err != nil {
			utils.JSONError(w, "failed to fetch players", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var players []models.Player
		for rows.Next() {
			var p models.Player
			err := rows.Scan(&p.ID, &p.Username, &p.ClientID)
			if err != nil {
				utils.JSONError(w, "error parsing players", http.StatusInternalServerError)
				return
			}
			players = append(players, p)
		}

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "players called successfully",
			Data:    players,
		}, http.StatusOK)

	}
}
