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
	r.HandleFunc("/games", controllers.GameListHandler).Methods("GET")
	r.HandleFunc("/games", controllers.GameCreateHandler).Methods("POST")
	http.Handle("/", r)

	fmt.Println("go-man server running")
	log.Fatal(http.ListenAndServe(":8080", r))
}
