package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"provider/router"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Missing DB_URL enviroment variable")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()

	router.RegisterCLientRoutes(r, db)
	router.RegisterPlayerRoutes(r, db)
	router.RegisterAdminRoutes(r, db)
	router.RegisterGameRoutes(r, db)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
