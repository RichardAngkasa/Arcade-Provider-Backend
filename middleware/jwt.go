package middleware

import (
	"errors"
	"net/http"
	"provider/utils"

	"github.com/golang-jwt/jwt/v5"
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
