package handlers

import (
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

	"gorm.io/gorm"
)

// ADMIN
func AdminLogin() http.HandlerFunc {
	return admin_auth.AdminLogin()
}
func AdminClientDeposit(db *gorm.DB) http.HandlerFunc {
	return admin_clients.AdminClientDeposit(db)
}
func AdminClientWithdraw(db *gorm.DB) http.HandlerFunc {
	return admin_clients.AdminClientWithdraw(db)
}
func AdminTransactions(db *gorm.DB) http.HandlerFunc {
	return admin_wallet.AdminTransactions(db)
}
func AdminClients(db *gorm.DB) http.HandlerFunc {
	return admin_clients.AdminClients(db)
}
func AdminClientProfile(db *gorm.DB) http.HandlerFunc {
	return admin_clients.AdminClientProfile(db)
}
func AdminPlayers(db *gorm.DB) http.HandlerFunc {
	return admin_players.AdminPlayers(db)
}
func AdminPlayerProfile(db *gorm.DB) http.HandlerFunc {
	return admin_players.AdminPlayerProfile(db)
}

// CLIENTS
func ClientRegister(db *gorm.DB) http.HandlerFunc {
	return client_auth.ClientRegister(db)
}
func ClientLogin(db *gorm.DB) http.HandlerFunc {
	return client_auth.ClientLogin(db)
}
func ClientLogout() http.HandlerFunc {
	return client_auth.ClientLogout()
}
func ClientProfile(db *gorm.DB) http.HandlerFunc {
	return client_profile.ClientProfile(db)
}
func ClientTransactions(db *gorm.DB) http.HandlerFunc {
	return client_wallet.ClientTransactions(db)
}
func ClientPlayerDeposit(db *gorm.DB) http.HandlerFunc {
	return client_players.ClientPlayerDeposit(db)
}
func ClientPlayerWithdraw(db *gorm.DB) http.HandlerFunc {
	return client_players.ClientPlayerWithdraw(db)
}
func ClientPlayers(db *gorm.DB) http.HandlerFunc {
	return client_players.ClientPlayers(db)
}

// PLAYER
func PlayerRegister(db *gorm.DB) http.HandlerFunc {
	return player_auth.PlayerRegister(db)
}
func PlayerProfile(db *gorm.DB) http.HandlerFunc {
	return player_profile.PlayerProfile(db)
}
func PlayerTransactions(db *gorm.DB) http.HandlerFunc {
	return player_wallet.PlayerTransactions(db)
}

// !GAME
func StartSpin(db *gorm.DB) http.HandlerFunc {
	return game.StartSpin(db)
}
