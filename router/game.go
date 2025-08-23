package router

import (
	"provider/handlers"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterGameRoutes(r *mux.Router, db *gorm.DB) {
	r.HandleFunc("/api/game/start", handlers.StartSpin(db)).Methods("POST")
}
