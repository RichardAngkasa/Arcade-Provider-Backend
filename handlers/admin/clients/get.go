package clients

import (
	"database/sql"
	"net/http"
	"provider/middleware"
	"provider/models"
	"provider/utils"
)

func AdminClientProfile(db *sql.DB) http.HandlerFunc {
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

		var client models.Client
		err = db.QueryRow(`
			SELECT id, username, email, password, api_key
			FROM clients
			WHERE id = $1
		`, req.ID).Scan(&client.ID, &client.Username, &client.Email, &client.Password, &client.ApiKey)
		if err != nil {
			utils.JSONError(w, "failed to fetch client profile", http.StatusInternalServerError)
			return
		}

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "client profile fetched",
			Data:    client,
		}, http.StatusOK)
	}
}
