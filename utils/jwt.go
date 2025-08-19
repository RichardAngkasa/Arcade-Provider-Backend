package utils

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("my-rahasia-key")

func GenerateJWT(id int, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":   id,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}

func GetIDFromToken(r *http.Request, cookieName string, expectedRole string) (int, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return 0, err
	}

	claims, err := ParseJWT(cookie.Value)
	if err != nil {
		return 0, err
	}

	role, ok := claims["role"].(string)
	if !ok || role != expectedRole {
		return 0, errors.New("invalid role in token")
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("invalid id in token")
	}

	return int(idFloat), nil
}
