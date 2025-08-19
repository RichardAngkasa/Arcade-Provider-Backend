package utils

import (
	"database/sql"
	"errors"
	"provider/models"
)

type Queryer interface {
	QueryRow(query string, args ...interface{}) *sql.Row
}

func ClientDeposit(tx *sql.Tx, client_id int, amount float64) error {
	result, err := tx.Exec(`
		UPDATE client_wallets 
		SET balance = balance + $1
		WHERE client_id = $2
	`, amount, client_id)
	if err != nil {
		return errors.New("failed to client deposit")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("could not check deposit result: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("deposit failed: client not found or no wallet got updated")
	}

	return nil
}

func ClientWithdraw(tx *sql.Tx, client_id int, amount float64) error {
	result, err := tx.Exec(`
		UPDATE client_wallets 
		SET balance = balance - $1
		WHERE client_id = $2
	`, amount, client_id)
	if err != nil {
		return errors.New("failed to client withdraw")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("could not check withdraw result: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("withdraw failed: client not found or no wallet got updated")
	}

	return nil
}

func PlayerDeposit(tx *sql.Tx, client_id, player_id int, amount float64) error {
	result, err := tx.Exec(`
		UPDATE player_wallets 
		SET balance = balance + $1
		WHERE player_id = $2 AND client_id = $3
	`, amount, player_id, client_id)
	if err != nil {
		return errors.New("failed to player deposit")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("could not check deposit result: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("deposit failed: player not found or no wallet got updated")
	}

	return nil
}

func PlayerWithdraw(tx *sql.Tx, client_id, player_id int, amount float64) error {
	result, err := tx.Exec(`
		UPDATE player_wallets 
		SET balance = balance - $1
		WHERE player_id = $2 AND client_id = $3
	`, amount, player_id, client_id)
	if err != nil {
		return errors.New("failed to player withdraw")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("could not check deposit result: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("deposit failed: player not found or no wallet got updated")
	}

	return nil
}

func ClientLogTransaction(tx *sql.Tx, client_wallet_id, client_id int, ammount float64, transaction_type string) error {
	switch transaction_type {
	case
		"deposit",
		"withdraw",
		"bet_win_player",
		"bet_lose_player":
	default:
		return errors.New("invalid transaction type: " + transaction_type)
	}

	_, err := tx.Exec(`
		INSERT INTO client_wallet_transactions (client_wallet_id, client_id, amount, type)
		VALUES ($1, $2, $3, $4)
	`, client_wallet_id, client_id, ammount, transaction_type)
	if err != nil {
		return errors.New("failed to log client transaction" + err.Error())
	}

	return nil
}

func PlayerLogTransaction(tx *sql.Tx, wallet_id, player_id, client_id int, game_session_id sql.NullInt64, ammount float64, transaction_type string) error {
	switch transaction_type {
	case
		"deposit",
		"withdraw",
		"bet_win_player",
		"bet_lose_player":
	default:
		return errors.New("invalid transaction type: " + transaction_type)
	}

	_, err := tx.Exec(`
		INSERT INTO player_wallet_transactions (player_wallet_id, player_id, client_id, game_session_id, amount, type)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, wallet_id, player_id, client_id, game_session_id, ammount, transaction_type)
	if err != nil {
		return errors.New("failed to log player transaction " + err.Error())
	}

	return nil
}

func AdminLogTransaction(tx *sql.Tx, client_id int, ammount float64, transaction_type string) error {
	switch transaction_type {
	case
		"deposit",
		"withdraw":
	default:
		return errors.New("invalid transaction type: " + transaction_type)
	}

	_, err := tx.Exec(`
		INSERT INTO admin_wallet_transactions (client_id, amount, type)
		VALUES ($1, $2, $3)
	`, client_id, ammount, transaction_type)
	if err != nil {
		return errors.New("failed to log admin transaction " + err.Error())
	}

	return nil
}

func ClientWallet(q Queryer, client_id int) (models.ClientWallet, error) {
	var wallet models.ClientWallet
	err := q.QueryRow(`
		SELECT id, client_id, balance 
		FROM client_wallets 
		WHERE client_id = $1
	`, client_id).Scan(&wallet.ID, &wallet.ClientID, &wallet.Balance)

	if err != nil {
		if err == sql.ErrNoRows {
			return wallet, errors.New("client wallet not found")
		}
		return wallet, errors.New("failed to fetch client wallet")
	}
	return wallet, nil
}

func PlayerWallet(q Queryer, client_id, player_id int) (models.PlayerWallet, error) {
	var wallet models.PlayerWallet
	err := q.QueryRow(`
		SELECT id, client_id, player_id, balance 
		FROM player_wallets 
		WHERE client_id = $1 AND player_id = $2
	`, client_id, player_id).Scan(&wallet.ID, &wallet.ClientID, &wallet.PlayerID, &wallet.Balance)

	if err != nil {
		if err == sql.ErrNoRows {
			return wallet, errors.New("player wallet not found")
		}
		return wallet, errors.New("failed to fetch player wallet")
	}
	return wallet, nil
}

func GameSessionLog(tx *sql.Tx, player_id, client_id int, betAmount, resultAmount float64, game_id, sessionType string) (int, error) {
	var gameSessionID int
	err := tx.QueryRow(`
		INSERT INTO game_sessions (player_id, game_id, bet_amount, type, result_amount, client_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		returning id
	`, player_id, game_id, betAmount, sessionType, resultAmount, client_id).Scan(&gameSessionID)
	if err != nil {
		return gameSessionID, errors.New("failed to log game sessions")
	}
	return gameSessionID, nil
}
