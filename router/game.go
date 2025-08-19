package router

import (
	"database/sql"
	"provider/handlers"

	"github.com/gorilla/mux"
)

func RegisterGameRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/api/game/start", handlers.StartSpin(db)).Methods("POST")
}
