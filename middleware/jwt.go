package middleware

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"provider/utils"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey int

const (
	CtxAdminID ctxKey = iota
	CtxClientID
	CtxRole
)

func Validate(c *http.Cookie, allowedRoles ...string) (jwt.MapClaims, error) {
	if c == nil {
		return nil, errors.New("missing token")
	}
	claims, err := utils.ParseJWT(c.Value)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid role")
	}
	for _, r := range allowedRoles {
		if role == r {
			return claims, nil
		}
	}
	return claims, errors.New("unauthorized")
}

func JwtAuthAdmin(db *sql.DB, next http.Handler) http.Handler {
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

func JwtAuthClient(db *sql.DB, next http.Handler) http.Handler {
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

func MustClientID(r *http.Request) (int, error) {
	id, ok := r.Context().Value(CtxClientID).(int)
	if !ok {
		return 0, errors.New("unauthorized")
	}
	return id, nil
}

func MustAdminID(r *http.Request) (int, error) {
	id, ok := r.Context().Value(CtxAdminID).(int)
	if !ok {
		return 0, errors.New("unauthorized")
	}
	return id, nil
}
