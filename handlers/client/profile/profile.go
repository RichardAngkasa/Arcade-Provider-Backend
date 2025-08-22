package clients

import (
	"database/sql"
	"net/http"
	"provider/middleware"
	"provider/models"
	utils "provider/utils"
)

type ClientProfileResponse struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	ApiKey   string  `json:"api_key"`
	Balance  float64 `json:"balance"`
}

func ClientProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID, err := middleware.MustClientID(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var client models.Client
		var wallet models.ClientWallet
		err = db.QueryRow(`
			SELECT c.id, c.username, c.email, c.api_key, w.balance
			FROM clients c
			LEFT JOIN client_wallets w ON c.id = w.client_id
			WHERE c.id = $1
		`, clientID).Scan(&client.ID, &client.Username, &client.Email, &client.ApiKey, &wallet.Balance)
		if err != nil {
			utils.JSONError(w, "failed to fetch client profile", http.StatusInternalServerError)
			return
		}

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "client profile fetched",
			Data: ClientProfileResponse{
				ID:       client.ID,
				Username: client.Username,
				Email:    client.Email,
				ApiKey:   client.ApiKey,
				Balance:  wallet.Balance,
			},
		}, http.StatusOK)

	}
}
