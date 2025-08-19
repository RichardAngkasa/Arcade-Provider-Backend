package clients

import (
	"database/sql"
	"net/http"
	"provider/models"
	"provider/utils"
)

func AdminClients(db *sql.DB) http.HandlerFunc {
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
			SELECT id, username, email
			FROM clients
		`)
		if err != nil {
			utils.JSONError(w, "failed to fetch clients", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var clients []models.Client
		for rows.Next() {
			var c models.Client
			err := rows.Scan(&c.ID, &c.Username, &c.Email)
			if err != nil {
				utils.JSONError(w, "error parsing clients", http.StatusInternalServerError)
				return
			}
			clients = append(clients, c)
		}

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "clients called successfully",
			Data:    clients,
		}, http.StatusOK)
	}
}
