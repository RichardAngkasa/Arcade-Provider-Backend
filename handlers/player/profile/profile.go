package player

import (
	"database/sql"
	"net/http"
	"provider/models"
	"provider/utils"
)

type ProfileRequest struct {
	ID int `json:"id"`
}

type PlayerProfileResponse struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
	ClientID int     `json:"client_id"`
}

func PlayerProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.JSONError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ProfileRequest
		ClientApiKey := r.Header.Get("X-API-Key")

		err := utils.BodyChecker(r, &req)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if ClientApiKey == "" {
			utils.JSONError(w, "api key missing", http.StatusBadRequest)
			return
		}

		clientID, err := utils.GetClientIdByApiKey(db, ClientApiKey)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var player models.Player
		err = db.QueryRow(`
			SELECT id, username 
			FROM players 
			WHERE id = $1 AND client_id = $2
		`, req.ID, clientID).Scan(&player.ID, &player.Username)
		if err != nil {
			utils.JSONError(w, "failed to fetch player profile", http.StatusInternalServerError)
			return
		}

		wallet, err := utils.PlayerWallet(db, clientID, req.ID)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "player profile fetched",
			Data: PlayerProfileResponse{
				ID:       player.ID,
				Username: player.Username,
				Balance:  wallet.Balance,
				ClientID: clientID,
			},
		}, http.StatusOK)
	}
}
