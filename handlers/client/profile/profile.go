package clients

import (
	"database/sql"
	"net/http"
	"provider/models"
	utils "provider/utils"
)

type ClientProfileResponse struct {
	ID       int     `json:"client_id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Balance  float64 `json:"balance"`
}

func ClientProfile(db *sql.DB) http.HandlerFunc {
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

		var client models.Client
		var wallet models.ClientWallet
		err = db.QueryRow(`
			SELECT c.id, c.username, c.email, w.balance
			FROM clients c
			LEFT JOIN client_wallets w ON c.id = w.client_id
			WHERE c.id = $1
		`, clientID).Scan(&client.ID, &client.Username, &client.Email, &wallet.Balance)
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
				Balance:  wallet.Balance,
			},
		}, http.StatusOK)

	}
}
