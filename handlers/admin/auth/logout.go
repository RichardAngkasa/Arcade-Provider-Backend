package auth

import (
	"fmt"
	"net/http"
	"provider/middleware"
	"provider/utils"
)

func AdminLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt_token_admin")
		if err != nil {
			utils.JSONError(w, "no cookie", http.StatusUnauthorized)
			return
		}
		claims, err := middleware.Validate(cookie, "client")
		if err != nil {
			utils.JSONError(w, "invalid token", http.StatusUnauthorized)
			return
		}
		clientID := int(claims["id"].(float64))
		redisKey := fmt.Sprintf("session:admin:%d", clientID)
		utils.RedisClient.Del(utils.Ctx, redisKey)

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt_token_admin",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "logged out",
		}, http.StatusOK)
	}
}
