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

	r.Handle("/api/client/players", middleware.JwtAuth(db, handlers.ClientPlayers(db))).Methods("GET")
	r.Handle("/api/client/player/deposit", middleware.JwtAuth(db, handlers.PlayerDeposit(db))).Methods("POST")
	r.Handle("/api/client/player/withdraw", middleware.JwtAuth(db, handlers.PlayerWithdraw(db))).Methods("POST")

	r.Handle("/api/client/profile", middleware.JwtAuth(db, handlers.ClientProfile(db))).Methods("GET")
	r.Handle("/api/client/transactions", middleware.JwtAuth(db, handlers.ClientTransactions(db))).Methods("GET")
}
