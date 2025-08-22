package router

import (
	"database/sql"
	"provider/handlers"

	"github.com/gorilla/mux"
)

func RegisterPlayerRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/api/player/register", handlers.PlayerRegister(db)).Methods("POST")
	r.HandleFunc("/api/player/profile", handlers.PlayerProfile(db)).Methods("GET")

	r.Handle("/api/player/transactions", handlers.PlayerTransactions(db)).Methods("GET")
}
