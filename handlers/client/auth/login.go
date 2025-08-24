package auth

import (
	"errors"
	"fmt"
	"net/http"
	"provider/models"
	utils "provider/utils"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ClientLogin(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
		var req models.ClientLoginRequest
		err := utils.BodyChecker(r, &req)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.Identifier = strings.ToLower(req.Identifier)
		if req.Identifier == "" || req.Password == "" {
			utils.JSONError(w, "username or email and password required", http.StatusBadRequest)
			return
		}

		// QUERYING
		var client models.Client
		err = db.
			Where("username = ? OR email = ?", req.Identifier, req.Identifier).
			First(&client).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.JSONError(w, "invalid credentials", http.StatusUnauthorized)
				return
			}
			utils.JSONError(w, "database error", http.StatusInternalServerError)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(req.Password))
		if err != nil {
			utils.JSONError(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		token, err := utils.GenerateJWT(client.ID, "client")
		if err != nil {
			utils.JSONError(w, "token generation failed", http.StatusInternalServerError)
			return
		}
		err = utils.RedisClient.Set(utils.Ctx,
			"session:client:"+fmt.Sprint(client.ID),
			token,
			24*time.Hour,
		).Err()
		if err != nil {
			utils.JSONError(w, "redis", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt_token_client",
			Value:    token,
			HttpOnly: true,
			Secure:   false,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		})

		// RESPONSE
		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "client login successfully",
			Data: models.ClientLoginResponse{
				APIKey: client.ApiKey,
				Token:  token,
			},
		}, http.StatusOK)

	}
}
