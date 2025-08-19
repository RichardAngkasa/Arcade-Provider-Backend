package handlers

import (
	"database/sql"
	"net/http"
	utils "provider/utils"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type ClientLoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type ClientLoginResponse struct {
	APIKey string `json:"api_key"`
	Token  string `json:"token"`
}

func ClientLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.JSONError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ClientLoginRequest
		err := utils.BodyChecker(r, &req)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		req.Identifier = strings.ToLower(req.Identifier)
		if req.Identifier == "" || req.Password == "" {
			utils.JSONError(w, "username or Email and password required", http.StatusBadRequest)
			return
		}

		var dbPassword, apiKey string
		var clientID int

		err = db.QueryRow(`
			SELECT password, api_key, id FROM clients
			WHERE username = $1 OR email = $1
		`, req.Identifier).Scan(&dbPassword, &apiKey, &clientID)

		if err != nil {
			if err == sql.ErrNoRows {
				utils.JSONError(w, "invalid credentials", http.StatusUnauthorized)
				return
			}
			utils.JSONError(w, "database error", http.StatusInternalServerError)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(req.Password))
		if err != nil {
			utils.JSONError(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateJWT(clientID, "client")
		if err != nil {
			utils.JSONError(w, "token generation failed", http.StatusInternalServerError)
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

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "client login successfully",
			Data: ClientLoginResponse{
				APIKey: apiKey,
				Token:  token,
			},
		}, http.StatusOK)

	}
}
