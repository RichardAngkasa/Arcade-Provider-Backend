package handlers

import (
	"database/sql"
	"net/http"
	client_auth "provider/handlers/client/auth"
	client_players "provider/handlers/client/players"
	client_profile "provider/handlers/client/profile"
	client_wallet "provider/handlers/client/wallet"

	player_auth "provider/handlers/player/auth"
	player_profile "provider/handlers/player/profile"
	player_wallet "provider/handlers/player/wallet"

	admin_auth "provider/handlers/admin/auth"
	admin_clients "provider/handlers/admin/clients"
	admin_players "provider/handlers/admin/players"
	admin_wallet "provider/handlers/admin/wallet"

	game "provider/handlers/game"
)

// !CLIENTS
// CLient Auth
func ClientRegister(db *sql.DB) http.HandlerFunc {
	return client_auth.ClientRegister(db)
}
func ClientLogin(db *sql.DB) http.HandlerFunc {
	return client_auth.ClientLogin(db)
}
func ClientLogout(db *sql.DB) http.HandlerFunc {
	return client_auth.ClientLogout()
}

// Client Profile
func ClientProfile(db *sql.DB) http.HandlerFunc {
	return client_profile.ClientProfile(db)
}

// Client Wallet
func ClientTransactions(db *sql.DB) http.HandlerFunc {
	return client_wallet.ClientTransactions(db)
}
func PlayerDeposit(db *sql.DB) http.HandlerFunc {
	return client_wallet.PlayerDeposit(db)
}
func PlayerWithdraw(db *sql.DB) http.HandlerFunc {
	return client_wallet.PlayerWithdraw(db)
}

// Client Players
func ClientPlayers(db *sql.DB) http.HandlerFunc {
	return client_players.ClientPlayers(db)
}

// !PLAYER
// Player Auth
func PlayerRegister(db *sql.DB) http.HandlerFunc {
	return player_auth.PlayerRegister(db)
}

// Player Profile
func PlayerProfile(db *sql.DB) http.HandlerFunc {
	return player_profile.PlayerProfile(db)
}

// Player Wallet
func PlayerTransactions(db *sql.DB) http.HandlerFunc {
	return player_wallet.PlayerTransactions(db)
}

// !GAME
func StartSpin(db *sql.DB) http.HandlerFunc {
	return game.StartSpin(db)
}

// !ADMIN
// Admin Auth
func AdminLogin(db *sql.DB) http.HandlerFunc {
	return admin_auth.AdminLogin(db)
}

// Admin banking
func ClientDeposit(db *sql.DB) http.HandlerFunc {
	return admin_wallet.ClientDeposit(db)
}
func ClientWithdraw(db *sql.DB) http.HandlerFunc {
	return admin_wallet.ClientWithdraw(db)
}
func AdminTransactions(db *sql.DB) http.HandlerFunc {
	return admin_wallet.AdminTransactions(db)
}

// Admin Clients
func AdminClients(db *sql.DB) http.HandlerFunc {
	return admin_clients.AdminClients(db)
}
func AdminClientProfile(db *sql.DB) http.HandlerFunc {
	return admin_clients.AdminClientProfile(db)
}
func AdminPlayers(db *sql.DB) http.HandlerFunc {
	return admin_players.AdminPlayers(db)
}
func AdminPlayerProfile(db *sql.DB) http.HandlerFunc {
	return admin_players.AdminPlayerProfile(db)
}
