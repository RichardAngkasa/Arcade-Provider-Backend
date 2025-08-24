package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"provider/models"

	"gorm.io/gorm"
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

func RequestAmountGreaterThanBalanceForbidden(ammount, balance float64) error {
	if ammount > balance {
		return errors.New("insufficient balance")
	} else {
		return nil
	}
}

func ClientExistenceByID(db *gorm.DB, client_id int) error {
	var client models.Client
	err := db.
		Where("id", client_id).
		First(&client).Error
	if err != nil {
		return errors.New("client not found")
	}
	return err
}

func PlayerAlreadyExistUnderClientByUsername(db *gorm.DB, client_id int, username string) error {
	var player models.Player
	err := db.
		Where("client_id = ? AND username = ?", client_id, username).
		First(&player).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	} else if err != nil {
		return err
	}
	return errors.New("player already exist under this client")
}

func PlayerMustExistUnderClient(db *gorm.DB, clientID, playerID int) error {
	var player models.Player
	err := db.
		Where("id = ? AND client_id = ?", playerID, clientID).
		First(&player).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("player not found under this client")
	} else if err != nil {
		return err
	}
	return nil
}

func ClientUniqueness(db *gorm.DB, username, email string) error {
	var client models.Client
	err := db.
		Where("username = ? OR email = ?", username, email).
		First(&client).Error
	if err == nil {
		return errors.New("username or email already exists")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func GetClientIdByApiKey(db *gorm.DB, api_key string) (int, error) {
	var client models.Client
	err := db.
		Select("id").
		Where("api_key = ?", api_key).
		First(&client).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, errors.New("invalid client")
	} else if err != nil {
		return 0, err
	}
	return client.ID, nil
}

func GetClientIdByHeader(db *gorm.DB, r *http.Request) (int, error) {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		return 0, errors.New("api key missing")
	}
	clientID, err := GetClientIdByApiKey(db, apiKey)
	if err != nil {
		return 0, errors.New(err.Error())
	}
	return clientID, nil
}
