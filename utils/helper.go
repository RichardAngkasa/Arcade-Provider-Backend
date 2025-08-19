package utils

import (
	"database/sql"
	"errors"
)

func GetClientIdByApiKey(db *sql.DB, api_key string) (int, error) {
	var clientID int
	err := db.QueryRow(`
		SELECT id FROM clients 
		WHERE api_key = $1
	`, api_key).Scan(&clientID)
	if err == sql.ErrNoRows {
		return 0, errors.New("invalid client")
	} else if err != nil {
		return 0, err
	}
	return clientID, nil
}
