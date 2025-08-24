package middleware

import (
	"context"
	"fmt"
	"net/http"
	"provider/utils"
)

func JwtAuthAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("jwt_token_admin")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusInternalServerError)
			return
		}
		claims, err := Validate(c, "admin")
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		adminID := int(claims["id"].(float64))
		redisKey := fmt.Sprintf("session:admin:%d", adminID)

		storedToken, err := utils.RedisClient.Get(utils.Ctx, redisKey).Result()
		if err != nil || storedToken != c.Value {
			utils.JSONError(w, "session expired or invalid", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), CtxAdminID, adminID)
		ctx = context.WithValue(ctx, CtxRole, "admin")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func JwtAuthClient(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("jwt_token_client")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusInternalServerError)
			return
		}
		claims, err := Validate(c, "client")
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		clientID := int(claims["id"].(float64))
		redisKey := fmt.Sprintf("session:client:%d", clientID)

		storedToken, err := utils.RedisClient.Get(utils.Ctx, redisKey).Result()
		if err != nil || storedToken != c.Value {
			utils.JSONError(w, "session expired or invalid", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), CtxClientID, clientID)
		ctx = context.WithValue(ctx, CtxRole, "client")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
