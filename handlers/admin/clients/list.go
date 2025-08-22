package clients

import (
	"database/sql"
	"net/http"
	"provider/middleware"
	"provider/models"
	"provider/utils"
)

func AdminClients(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := middleware.MustAdminID(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
			SELECT id, username, email, password, api_key
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
			err := rows.Scan(&c.ID, &c.Username, &c.Email, &c.Password, &c.ApiKey)
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
