package router

import (
	"provider/handlers"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterPlayerRoutes(r *mux.Router, db *gorm.DB) {
	r.HandleFunc("/api/player/register", handlers.PlayerRegister(db)).Methods("POST")
	r.HandleFunc("/api/player/profile", handlers.PlayerProfile(db)).Methods("GET")

	r.Handle("/api/player/transactions", handlers.PlayerTransactions(db)).Methods("GET")
}
