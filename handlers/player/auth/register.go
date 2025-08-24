package handlers

import (
	"net/http"
	"provider/models"
	"provider/utils"
	"strings"

	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username string `json:"username"`
}

func PlayerRegister(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
		clientID, err := utils.GetClientIdByHeader(db, r)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		var req RegisterRequest
		err = utils.BodyChecker(r, &req)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.Username = strings.ToLower(req.Username)
		if req.Username == "" {
			utils.JSONError(w, "username required", http.StatusBadRequest)
			return
		}
		err = utils.PlayerAlreadyExistUnderClientByUsername(db, clientID, req.Username)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// QUERY
		player := models.Player{
			Username: req.Username,
			ClientID: clientID,
		}
		err = db.
			Create(&player).Error
		if err != nil {
			utils.JSONError(w, "insert new player failed", http.StatusInternalServerError)
			return
		}
		playerWallet := models.PlayerWallet{
			PlayerID: player.ID,
			ClientID: clientID,
			Balance:  0,
		}
		err = db.
			Create(&playerWallet).Error
		if err != nil {
			utils.JSONError(w, "failed creating new wallet for player", http.StatusInternalServerError)
			return
		}

		// RESPONSE
		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "client registered successfully",
		}, http.StatusCreated)

	}
}
