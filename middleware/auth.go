package middleware

import (
	"context"
	"net/http"
	"provider/utils"
)

func JwtAuthAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("jwt_token_admin")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		claims, err := Validate(c, "admin")
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		adminID := int(claims["id"].(float64))
		ctx := context.WithValue(r.Context(), CtxAdminID, adminID)
		ctx = context.WithValue(ctx, CtxRole, "admin")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func JwtAuthClient(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("jwt_token_client")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		claims, err := Validate(c, "client")
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		clientID := int(claims["id"].(float64))
		ctx := context.WithValue(r.Context(), CtxClientID, clientID)
		ctx = context.WithValue(ctx, CtxRole, "client")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
