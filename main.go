package main

import (
	"Blockchain/controllers"
	"Blockchain/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	// Create new Mux Router
	r := mux.NewRouter()

	// Create new Blockchain Controller
	blockchainController := controllers.NewBlockchainController()

	// Create routes
	routes.CreateRoutes(r, blockchainController)

	// Listen and Serve on port 8080
	log.Fatalln(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r)))
}
