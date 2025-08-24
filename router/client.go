package router

import (
	"provider/handlers"
	"provider/middleware"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterCLientRoutes(r *mux.Router, db *gorm.DB) {
	r.HandleFunc("/api/client/register", handlers.ClientRegister(db)).Methods("POST")
	r.HandleFunc("/api/client/login", handlers.ClientLogin(db)).Methods("POST")
	r.HandleFunc("/api/client/logout", handlers.ClientLogout()).Methods("POST")

	r.Handle("/api/client/players", middleware.JwtAuthClient(handlers.ClientPlayers(db))).Methods("GET")
	r.Handle("/api/client/player/profile", middleware.JwtAuthClient(handlers.ClientPlayerProfile(db))).Methods("GET")
	r.Handle("/api/client/player/deposit", middleware.JwtAuthClient(handlers.ClientPlayerDeposit(db))).Methods("POST")
	r.Handle("/api/client/player/withdraw", middleware.JwtAuthClient(handlers.ClientPlayerWithdraw(db))).Methods("POST")

	r.Handle("/api/client/profile", middleware.JwtAuthClient(handlers.ClientProfile(db))).Methods("GET")
	r.Handle("/api/client/transactions", middleware.JwtAuthClient(handlers.ClientTransactions(db))).Methods("GET")
}
