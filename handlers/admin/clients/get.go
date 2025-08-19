package clients

import (
	"database/sql"
	"net/http"
	"provider/models"
	"provider/utils"
)

func AdminClientProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.JSONError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		_, err := utils.GetIDFromToken(r, "jwt_token_admin", "admin")
		if err != nil {
			utils.JSONError(w, "unauthorize", http.StatusUnauthorized)
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
			SELECT id, username, email
			FROM clients
			WHERE id = $1
		`, req.ID).Scan(&client.ID, &client.Username, &client.Email)
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
