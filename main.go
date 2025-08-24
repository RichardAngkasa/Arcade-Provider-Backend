package main

import (
	"fmt"
	"log"
	"net/http"
	"provider/config"
	"provider/router"
	"provider/utils"

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

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
