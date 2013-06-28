package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/telecoda/go-man/controllers"
	"log"
	"net/http"
)

func main() {

	fmt.Println("go-man server starting")

	r := mux.NewRouter()
	r.HandleFunc("/", controllers.RootHandler).Methods("GET")
	// list games
	r.HandleFunc("/games", controllers.GameList).Methods("GET")
	// create new game
	r.HandleFunc("/games", controllers.GameCreate).Methods("POST")
	// get game by id
	r.HandleFunc("/games/{gameId}", controllers.GameById).Methods("GET")
	// update MainPlayer
	r.HandleFunc("/games/{gameId}", controllers.UpdatePlayer).Methods("PUT")

	// options
	r.HandleFunc("/{path:.*}", controllers.OptionsHandler).Methods("OPTIONS")

	http.Handle("/", r)

	fmt.Println("go-man server running")
	log.Fatal(http.ListenAndServe(":8080", r))
}
