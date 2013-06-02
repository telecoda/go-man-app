package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/telecoda/go-man/models"
	"net/http"
)

func GameListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, Response{"success": true, "message": "Here are the current games", "method": r.Method})
}

func GameCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var board = models.NewGameBoard()

	bJson, err := json.Marshal(board)

	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Fprint(w, string(bJson))
	}
}
