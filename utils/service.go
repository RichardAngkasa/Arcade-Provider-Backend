package utils

import (
	"database/sql"
	"errors"
	"provider/models"

	"gorm.io/gorm"
)

// CLIENT
func ClientWallet(tx *gorm.DB, client_id int) (models.ClientWallet, error) {
	var wallet models.ClientWallet
	err := tx.
		Where("client_id = ?", client_id).
		First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return wallet, errors.New("client wallet not found")
	}
	return wallet, nil
}
func ClientAddBalance(tx *gorm.DB, client_id int, amount float64) error {
	wallet, err := ClientWallet(tx, client_id)
	if err != nil {
		return errors.New(err.Error())
	}
	wallet.Balance += amount
	err = tx.
		Save(&wallet).Error
	if err != nil {
		return errors.New("failed to update client wallet balance: " + err.Error())
	}
	return nil
}
func ClientDeductBalance(tx *gorm.DB, client_id int, amount float64) error {
	wallet, err := ClientWallet(tx, client_id)
	if err != nil {
		return errors.New(err.Error())
	}
	err = RequestAmountGreaterThanBalanceForbidden(amount, wallet.Balance)
	if err != nil {
		return errors.New(err.Error())
	}
	wallet.Balance -= amount
	err = tx.
		Save(&wallet).Error
	if err != nil {
		return errors.New("failed to update client wallet balance: " + err.Error())
	}
	return nil
}
func ClientLogTransaction(tx *gorm.DB, client_id int, ammount float64, transaction_type string) error {
	switch transaction_type {
	case
		"deposit",
		"withdraw",
		"bet_win_player",
		"bet_lose_player":
	default:
		return errors.New("invalid transaction type: " + transaction_type)
	}
	wallet, err := ClientWallet(tx, client_id)
	if err != nil {
		return errors.New(err.Error())
	}
	ClientLogTransaction := models.ClientWalletTransaction{
		ClientID:       client_id,
		Amount:         ammount,
		Type:           transaction_type,
		ClientWalletID: wallet.ID,
	}
	err = tx.
		Create(&ClientLogTransaction).Error
	if err != nil {
		return errors.New("failed to log client transaction" + err.Error())
	}
	return nil
}

// PLAYER
func PlayerWallet(tx *gorm.DB, client_id, player_id int) (models.PlayerWallet, error) {
	var wallet models.PlayerWallet
	err := tx.
		Where("client_id = ? AND player_id", client_id, player_id).
		First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return wallet, errors.New("player wallet not found")
	}
	return wallet, nil
}
func PlayerAddBalance(tx *gorm.DB, client_id, player_id int, amount float64) error {
	wallet, err := PlayerWallet(tx, client_id, player_id)
	if err != nil {
		return errors.New(err.Error())
	}
	wallet.Balance += amount
	if err != nil {
		return errors.New("failed to update player wallet balance: " + err.Error())
	}
	return nil
}
func PlayerDeductBalance(tx *gorm.DB, client_id, player_id int, amount float64) error {
	wallet, err := PlayerWallet(tx, client_id, player_id)
	if err != nil {
		return errors.New(err.Error())
	}
	err = RequestAmountGreaterThanBalanceForbidden(amount, wallet.Balance)
	if err != nil {
		return errors.New(err.Error())
	}
	wallet.Balance -= amount
	if err != nil {
		return errors.New("failed to update player wallet balance: " + err.Error())
	}
	return nil
}
func PlayerLogTransaction(tx *gorm.DB, client_id, player_id int, game_session_id sql.NullInt64, ammount float64, transaction_type string) error {
	switch transaction_type {
	case
		"deposit",
		"withdraw",
		"bet_win_player",
		"bet_lose_player":
	default:
		return errors.New("invalid transaction type: " + transaction_type)
	}
	wallet, err := PlayerWallet(tx, client_id, player_id)
	if err != nil {
		return errors.New(err.Error())
	}
	PlayerLogTransaction := models.PlayerWalletTransaction{
		ClientID:      client_id,
		PlayerID:      player_id,
		WalletID:      wallet.ID,
		GameSessionID: game_session_id,
		Amount:        ammount,
		Type:          transaction_type,
	}
	err = tx.
		Create(&PlayerLogTransaction).Error
	if err != nil {
		return errors.New("failed to log player transaction " + err.Error())
	}
	return nil
}

// ADMIN
func AdminLogTransaction(tx *gorm.DB, client_id int, ammount float64, transaction_type string) error {
	switch transaction_type {
	case
		"deposit",
		"withdraw":
	default:
		return errors.New("invalid transaction type: " + transaction_type)
	}
	AdminLogTransaction := models.AdminWalletTransaction{
		ClientID: client_id,
		Amount:   ammount,
		Type:     transaction_type,
	}
	err := tx.
		Create(&AdminLogTransaction).Error
	if err != nil {
		return errors.New("failed to log admin transaction " + err.Error())
	}
	return nil
}

// GAME
func GameLogSession(tx *gorm.DB, client_id, player_id int, betAmount, resultAmount float64, game_id, sessionType string) (int, error) {
	switch sessionType {
	case
		"win",
		"lose":
	default:
		return 0, errors.New("invalid transaction type: " + sessionType)
	}
	GameSession := models.GameSession{
		ClientID:     client_id,
		PlayerID:     player_id,
		GameID:       game_id,
		BetAmount:    betAmount,
		ResultAmount: resultAmount,
		Type:         sessionType,
	}
	err := tx.
		Create(&GameSession).Error
	if err != nil {
		return GameSession.ID, errors.New("failed to log game sessions")
	}
	return GameSession.ID, nil
}
