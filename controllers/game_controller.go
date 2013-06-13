package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/telecoda/go-man/models"
	"net/http"
)

func GameList(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, Response{"success": true, "message": "Here are the current games", "method": r.Method})
}

func GameCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var board = models.NewGameBoard()

	board.SaveGameBoard()

	bJson, err := json.Marshal(board)

	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Fprint(w, string(bJson))
	}
}

func GameById(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Println("Getting game board", id)
	var board, err = models.LoadGameBoard(id)

	if board == nil || err != nil {
		http.NotFound(w, r)
	}

	bJson, err := json.Marshal(board)

	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Fprint(w, string(bJson))
	}
}
