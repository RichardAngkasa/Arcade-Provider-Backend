package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"provider/models"
	"provider/utils"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func generatedAPIKey() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

type ClientRegisterResponse struct {
	APIKey string `json:"api_key"`
}

func ClientRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.Client
		err := utils.BodyChecker(r, &req)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		req.Username = strings.ToLower(req.Username)
		req.Email = strings.ToLower(req.Email)
		if req.Username == "" || req.Password == "" || req.Email == "" {
			utils.JSONError(w, "all fields required", http.StatusBadRequest)
			return
		}

		if len(req.Password) < 8 {
			utils.JSONError(w, "password must be at least 8 chars", http.StatusBadRequest)
			return
		}

		if !strings.Contains(req.Email, "@") {
			utils.JSONError(w, "email must contain @ symbol", http.StatusBadRequest)
			return
		}

		err = utils.ClientUniqueness(db, req.Username, req.Email)
		if err != nil {
			utils.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		apiKey := generatedAPIKey()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		if err != nil {
			utils.JSONError(w, "error hashing password", http.StatusInternalServerError)
			return
		}

		var clientID int
		err = db.QueryRow(`
			INSERT INTO clients (username, email, password, api_key)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, req.Username, req.Email, hashedPassword, apiKey).Scan(&clientID)
		if err != nil {
			utils.JSONError(w, "failed insert new client", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
			INSERT INTO client_wallets (client_id, balance)
			VALUES ($1, $2)
		`, clientID, 0)
		if err != nil {
			utils.JSONError(w, "failed create client wallet", http.StatusInternalServerError)
			return
		}

		utils.JSONResponse(w, utils.Response{
			Success: true,
			Message: "client registered successfully",
			Data: ClientRegisterResponse{
				APIKey: apiKey,
			},
		}, http.StatusCreated)
	}
}
