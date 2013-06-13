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
	r.HandleFunc("/", controllers.RootHandler)
	r.HandleFunc("/games", controllers.GameList).Methods("GET")
	r.HandleFunc("/games", controllers.GameCreate).Methods("POST")
	r.HandleFunc("/games/{id}", controllers.GameById).Methods("GET")
	http.Handle("/", r)

	fmt.Println("go-man server running")
	log.Fatal(http.ListenAndServe(":8080", r))
}
