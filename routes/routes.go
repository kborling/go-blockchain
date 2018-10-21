package routes

import (
	"Blockchain/controllers"
	"github.com/gorilla/mux"
)

// Creates routes for getting and writing to the blockchain
// Write to blockchain via POST request using JSON ex. {"Data": "Test"}
func CreateRoutes(r *mux.Router, bc *controllers.BlockchainController) {
	r.HandleFunc("/", bc.HandleGetBlockchain).Methods("GET")
	r.HandleFunc("/", bc.HandleWriteBlock).Methods("POST")
}
