package players

import (
	"database/sql"
	"net/http"
	"provider/models"
	"provider/utils"
)

func AdminPlayers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.JSONError(w, "methode not allowed", http.StatusMethodNotAllowed)
			return
		}

		_, err := utils.GetIDFromToken(r, "jwt_token_admin", "admin")
		if err != nil {
			utils.JSONError(w, "unauthorized", http.StatusUnauthorized)
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
