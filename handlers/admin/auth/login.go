package auth

import (
	"net/http"
	"os"
	"provider/models"
	"provider/utils"
)

func AdminLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// AUTH
		var req models.AdminLoginRequest
		err := utils.BodyChecker(r, &req)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		adminPassword := os.Getenv("ADMIN_PASSWORD")
		if req.Password != adminPassword {
			utils.JSONError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// OPERATE
		token, err := utils.GenerateJWT(1, "admin")
		if err != nil {
			utils.JSONError(w, "token generation failed", http.StatusInternalServerError)
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

		// RESPONSE
		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "admin login successfully",
			Data: models.AdminLoginResponse{
				Token: token,
			},
		}, http.StatusOK)
	}
}
