package main

import (
	"fmt"
	"log"
	"net/http"
	"provider/config"
	"provider/router"
	"provider/utils"

	gorillaHandlers "github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	config.ConnectDatabase()
	utils.InitRedis()

	r := mux.NewRouter()

	router.RegisterCLientRoutes(r, config.DB)
	router.RegisterPlayerRoutes(r, config.DB)
	router.RegisterAdminRoutes(r, config.DB)
	router.RegisterGameRoutes(r, config.DB)

	cors := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"http://localhost:3000"}),
		gorillaHandlers.AllowedOrigins([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		gorillaHandlers.AllowCredentials(),
	)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", cors(r)))
}
