package router

import (
	"database/sql"
	"provider/handlers"
	"provider/middleware"

	"github.com/gorilla/mux"
)

func RegisterCLientRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/api/client/register", handlers.ClientRegister(db)).Methods("POST")
	r.HandleFunc("/api/client/login", handlers.ClientLogin(db)).Methods("POST")
	r.HandleFunc("/api/client/logout", handlers.ClientLogout(db)).Methods("POST")

	r.Handle("/api/client/players", middleware.JwtAuthClient(db, handlers.ClientPlayers(db))).Methods("GET")
	r.Handle("/api/client/player/deposit", middleware.JwtAuthClient(db, handlers.PlayerDeposit(db))).Methods("POST")
	r.Handle("/api/client/player/withdraw", middleware.JwtAuthClient(db, handlers.PlayerWithdraw(db))).Methods("POST")

	r.Handle("/api/client/profile", middleware.JwtAuthClient(db, handlers.ClientProfile(db))).Methods("GET")
	r.Handle("/api/client/transactions", middleware.JwtAuthClient(db, handlers.ClientTransactions(db))).Methods("GET")
}
