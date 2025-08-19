package wallet

import (
	"database/sql"
	"net/http"
	"provider/models"
	"provider/utils"
)

func ClientTransactions(db *sql.DB) http.HandlerFunc {
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
			SELECT id, client_id, amount, type
			FROM client_wallet_transactions
			WHERE client_id = $1
			ORDER BY created_at DESC
		`, clientID)
		if err != nil {
			utils.JSONError(w, "failed to fetch client transactions", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var transactions []models.ClientWalletTransaction
		for rows.Next() {
			var tx models.ClientWalletTransaction
			err := rows.Scan(&tx.ID, &tx.ClientID, &tx.Amount, &tx.Type)
			if err != nil {
				utils.JSONError(w, "error parsing transactions", http.StatusInternalServerError)
				return
			}
			transactions = append(transactions, tx)
		}

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "client transactions called successfully",
			Data:    transactions,
		}, http.StatusOK)
	}
}
