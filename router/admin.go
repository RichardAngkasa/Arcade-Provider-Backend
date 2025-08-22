package router

import (
	"database/sql"
	"provider/handlers"
	"provider/middleware"

	"github.com/gorilla/mux"
)

func RegisterAdminRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/api/admin/login", handlers.AdminLogin(db)).Methods("POST")

	r.Handle("/api/admin/client/deposit", middleware.JwtAuthAdmin(db, handlers.ClientDeposit(db))).Methods("POST")
	r.Handle("/api/admin/client/withdraw", middleware.JwtAuthAdmin(db, handlers.ClientWithdraw(db))).Methods("POST")
	r.Handle("/api/admin/transactions", middleware.JwtAuthAdmin(db, handlers.AdminTransactions(db))).Methods("GET")

	r.Handle("/api/admin/clients", middleware.JwtAuthAdmin(db, handlers.AdminClients(db))).Methods("GET")
	r.Handle("/api/admin/client/profile", middleware.JwtAuthAdmin(db, handlers.AdminClientProfile(db))).Methods("GET")

	r.Handle("/api/admin/players", middleware.JwtAuthAdmin(db, handlers.AdminPlayers(db))).Methods("GET")
	r.Handle("/api/admin/player/profile", middleware.JwtAuthAdmin(db, handlers.AdminPlayerProfile(db))).Methods("GET")
}
