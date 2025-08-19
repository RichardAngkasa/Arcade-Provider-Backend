package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"provider/models"
	"provider/utils"
)

func AdminLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AdminLoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			utils.JSONError(w, "Invalid request", http.StatusBadRequest)
			return
		}

		adminPassword := os.Getenv("ADMIN_PASSWORD")
		if req.Password != adminPassword {
			utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateJWT(1, "admin")
		if err != nil {
			utils.JSONError(w, "Token generation failed", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt_token_admin",
			Value:    token,
			HttpOnly: true,
			Secure:   false,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		})

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "Admin login successfully",
			Data: models.AdminLoginResponse{
				Token: token,
			},
		}, http.StatusOK)
	}
}
