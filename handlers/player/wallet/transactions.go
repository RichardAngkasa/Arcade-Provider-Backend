package wallet

import (
	"database/sql"
	"net/http"
	"provider/models"
	"provider/utils"
)

func PlayerTransactions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ClientApiKey := r.Header.Get("X-API-Key")
		clientID, err := utils.GetClientIdByApiKey(db, ClientApiKey)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var req models.PlayerTransactionsRequest
		err = utils.BodyChecker(r, &req)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = utils.PlayerMustExistUnderClient(db, clientID, req.PlayerID)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusForbidden)
			return
		}

		rows, err := db.Query(`
			SELECT id, client_id, player_id, player_wallet_id, amount, type
			FROM player_wallet_transactions
			WHERE player_id = $1 AND client_id = $2
			ORDER BY created_at DESC
		`, req.PlayerID, clientID)
		if err != nil {
			utils.JSONError(w, "Failed to fetch player transactions", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var transactions []models.PlayerWalletTransaction
		for rows.Next() {
			var tx models.PlayerWalletTransaction
			err := rows.Scan(&tx.ID, &tx.ClientID, &tx.PlayerID, &tx.WalletID, &tx.Amount, &tx.Type)
			if err != nil {
				utils.JSONError(w, "Error parsing transactions", http.StatusInternalServerError)
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
