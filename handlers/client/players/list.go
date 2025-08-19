package players

import (
	"database/sql"
	"net/http"
	"provider/models"
	"provider/utils"
)

func ClientPlayers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.JSONError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		clientID, err := utils.GetIDFromToken(r, "jwt_token_client", "client")
		if err != nil {
			utils.JSONError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
			SELECT id, username, client_id
			FROM players
			WHERE client_id = $1
			ORDER BY created_at DESC
		`, clientID)
		if err != nil {
			utils.JSONError(w, "failed to fetch client players", http.StatusInternalServerError)
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
			Message: "client players called successfully",
			Data:    players,
		}, http.StatusOK)
	}
}
