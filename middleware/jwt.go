package middleware

import (
	"database/sql"
	"net/http"
	"provider/utils"
)

func JwtAuth(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cookie *http.Cookie
		var err error

		cookie, err = r.Cookie("jwt_token_admin")
		role := "admin"
		if err != nil {
			cookie, err = r.Cookie("jwt_token_client")
			role = "client"
		}

		if err != nil {
			utils.JSONError(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseJWT(cookie.Value)
		if err != nil {
			utils.JSONError(w, "invalid token", http.StatusUnauthorized)
			return
		}

		if role == "admin" {
			if claims["role"] != "admin" {
				utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		if role == "client" {
			userID := int(claims["id"].(float64))

			var exists bool
			err = db.QueryRow(`
				SELECT EXISTS(SELECT 1 FROM clients 
				WHERE id = $1)
			`, userID).Scan(&exists)
			if err != nil || !exists {
				utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
