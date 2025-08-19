package handlers

import (
	"database/sql"
	"net/http"
	"provider/utils"
	"strings"
)

type RegisterRequest struct {
	Username string `json:"username"`
}

func PlayerRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.JSONError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ClientApiKey := r.Header.Get("X-API-Key")
		if ClientApiKey == "" {
			utils.JSONError(w, "api key missing", http.StatusBadRequest)
			return
		}

		var req RegisterRequest
		err := utils.BodyChecker(r, &req)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		req.Username = strings.ToLower(req.Username)
		if req.Username == "" {
			utils.JSONError(w, "username required", http.StatusBadRequest)
			return
		}

		clientID, err := utils.GetClientIdByApiKey(db, ClientApiKey)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		err = utils.PlayerExistenceUnderClientByUsername(db, clientID, req.Username)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		var playerID int
		err = db.QueryRow(`
			INSERT INTO players (username, client_id) 
			VALUES ($1, $2) RETURNING id
		`, req.Username, clientID).Scan(&playerID)
		if err != nil {
			utils.JSONError(w, "insert new user failed", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
			INSERT INTO player_wallets (player_id, balance, client_id) 
			VALUES ($1, $2, $3)
		`, playerID, 0, clientID)
		if err != nil {
			utils.JSONError(w, "creating new wallet for player failed", http.StatusInternalServerError)
			return
		}

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "client registered successfully",
		}, http.StatusCreated)

	}
}
