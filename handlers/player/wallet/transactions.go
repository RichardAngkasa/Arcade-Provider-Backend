package wallet

import (
	"net/http"
	"provider/models"
	"provider/utils"

	"gorm.io/gorm"
)

func PlayerTransactions(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
		clientID, err := utils.GetClientIdByHeader(db, r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		playerID, err := utils.GetIDfromQuery(r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
		}
		err = utils.PlayerMustExistUnderClient(db, clientID, playerID)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusForbidden)
			return
		}

		// QUERY
		var transactions []models.PlayerWalletTransactions
		err = db.
			Where("player_id = ? AND client_id = ?", playerID, clientID).
			Order("created_at DESC").
			Find(&transactions).Error
		if err != nil {
			utils.JSONError(w, "Failed to fetch player transactions", http.StatusInternalServerError)
			return
		}

		// RESPONSE
		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "player transactions called successfully",
			Data:    transactions,
		}, http.StatusOK)
	}
}
