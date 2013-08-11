package main

import (
	//"fmt"
	"github.com/gorilla/mux"
	"github.com/telecoda/go-man-app/controllers"
	"net/http"
)

var r *mux.Router

func init() {

	//fmt.Println("go-man server starting1")
	r = mux.NewRouter()

	r.HandleFunc("/", controllers.RootHandler).Methods("GET")
	//fmt.Println("go-man server starting2")
	// list games
	r.HandleFunc("/games", controllers.GameList).Methods("GET")
	// create new game
	r.HandleFunc("/games", controllers.GameCreate).Methods("POST")
	// get game by id
	r.HandleFunc("/games/{gameId}", controllers.GameById).Methods("GET")
	// add new player to game
	r.HandleFunc("/games/{gameId}/players", controllers.AddPlayer).Methods("POST")
	// update player
	r.HandleFunc("/games/{gameId}/players/{playerId}", controllers.UpdatePlayer).Methods("PUT")

	// options
	r.HandleFunc("/{path:.*}", controllers.OptionsHandler).Methods("OPTIONS")

	http.Handle("/", r)

}
