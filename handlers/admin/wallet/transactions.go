package wallet

import (
	"database/sql"
	"net/http"
	"provider/middleware"
	"provider/models"
	"provider/utils"
)

func AdminTransactions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := middleware.MustAdminID(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
			SELECT id, client_id, amount, type
			FROM admin_wallet_transactions
			ORDER BY created_at DESC
		`)
		if err != nil {
			utils.JSONError(w, "failed to fetch admin transactions", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var transactions []models.AdminWalletTransaction
		for rows.Next() {
			var tx models.AdminWalletTransaction
			err := rows.Scan(&tx.ID, &tx.ClientID, &tx.Amount, &tx.Type)
			if err != nil {
				utils.JSONError(w, "error parsing transactions", http.StatusInternalServerError)
				return
			}
			transactions = append(transactions, tx)
		}

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "admin transactions called successfully",
			Data:    transactions,
		}, http.StatusOK)
	}
}
