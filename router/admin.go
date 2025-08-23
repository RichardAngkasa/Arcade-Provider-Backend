package router

import (
	"provider/handlers"
	"provider/middleware"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterAdminRoutes(r *mux.Router, db *gorm.DB) {
	r.HandleFunc("/api/admin/login", handlers.AdminLogin()).Methods("POST")

	r.Handle("/api/admin/client/deposit", middleware.JwtAuthAdmin(handlers.AdminClientDeposit(db))).Methods("POST")
	r.Handle("/api/admin/client/withdraw", middleware.JwtAuthAdmin(handlers.AdminClientWithdraw(db))).Methods("POST")
	r.Handle("/api/admin/transactions", middleware.JwtAuthAdmin(handlers.AdminTransactions(db))).Methods("GET")

	r.Handle("/api/admin/clients", middleware.JwtAuthAdmin(handlers.AdminClients(db))).Methods("GET")
	r.Handle("/api/admin/client/profile", middleware.JwtAuthAdmin(handlers.AdminClientProfile(db))).Methods("GET")

	r.Handle("/api/admin/players", middleware.JwtAuthAdmin(handlers.AdminPlayers(db))).Methods("GET")
	r.Handle("/api/admin/player/profile", middleware.JwtAuthAdmin(handlers.AdminPlayerProfile(db))).Methods("GET")
}
