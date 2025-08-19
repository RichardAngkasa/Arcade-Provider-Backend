package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

func BodyChecker(r *http.Request, req interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return errors.New("invalid request")
	}
	return nil
}

func AmountLessThanZero(ammount float64) error {
	if ammount <= 0 {
		return errors.New("amount must be greater than 0")
	} else {
		return nil
	}
}

func RequestAmountGreaterThanBalance(ammount, balance float64) error {
	if ammount > balance {
		return errors.New("insufficient balance")
	} else {
		return nil
	}
}

func RequestAmountLessThanBalance(ammount, balance float64) error {
	if ammount < balance {
		return nil
	} else {
		return errors.New("insufficient player balance")
	}
}

func ClientExistenceByID(db *sql.DB, client_id int) error {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM clients WHERE id = $1
		)
	`, client_id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("client not found")
	}
	return nil
}

func PlayerExistenceUnderClientByUsername(db *sql.DB, client_id int, username string) error {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM players 
			WHERE client_id = $1 AND username = $2
		)
	`, client_id, username).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("player already exist under this client")
	}
	return nil
}

func PlayerMustExistUnderClient(db *sql.DB, clientID, playerID int) error {
	var exists bool
	err := db.QueryRow(`
        SELECT EXISTS (
            SELECT 1 FROM players 
            WHERE client_id = $1 AND id  = $2
        )
    `, clientID, playerID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("player not found under this client")
	}
	return nil
}

func ClientUniqueness(db *sql.DB, username, email string) error {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM clients WHERE username=$1 OR email=$2
		)
	`, username, email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username or email already exists")
	}
	return nil
}
