package auth

import (
	"net/http"
	"provider/models"
	"provider/utils"
	"strings"

	"gorm.io/gorm"
)

func ClientRegister(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
		var req models.Client
		err := utils.BodyChecker(r, &req)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.Username = strings.ToLower(req.Username)
		req.Email = strings.ToLower(req.Email)
		if req.Username == "" || req.Password == "" || req.Email == "" {
			utils.JSONError(w, "all fields required", http.StatusBadRequest)
			return
		}
		if len(req.Password) < 8 {
			utils.JSONError(w, "password must be at least 8 chars", http.StatusBadRequest)
			return
		}
		if !strings.Contains(req.Email, "@") {
			utils.JSONError(w, "email must contain @ symbol", http.StatusBadRequest)
			return
		}
		err = utils.ClientUniqueness(db, req.Username, req.Email)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// QUERY
		apiKey := utils.GeneratedAPIKey()
		hashedPassword, err := utils.HashedPassword(req.Password)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		client := models.Client{
			Username: req.Username,
			Email:    req.Email,
			Password: hashedPassword,
			ApiKey:   apiKey,
		}
		err = db.Create(&client).Error
		if err != nil {
			utils.JSONError(w, "failed create client", http.StatusInternalServerError)
			return
		}
		clientWallet := models.ClientWallet{
			ClientID: client.ID,
			Balance:  0,
		}
		err = db.
			Create(&clientWallet).Error
		if err != nil {
			utils.JSONError(w, "failed create client wallet", http.StatusInternalServerError)
			return
		}

		// RESPONSE
		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "client registered successfully",
			Data: models.ClientRegisterResponse{
				APIKey: apiKey,
			},
		}, http.StatusCreated)
	}
}
