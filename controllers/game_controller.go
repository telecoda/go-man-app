package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/telecoda/go-man/models"
	"log"
	"net/http"
)

/* Thanks to the following source for an ASCII version of the game board
http://4coder.org/c-c-source-code/152/pacman/board.c.html

*/

func GameList(w http.ResponseWriter, r *http.Request) {
	addResponseHeaders(w)
	fmt.Fprint(w, Response{"success": true, "message": "Here are the current games", "method": r.Method})
}

func GameCreate(w http.ResponseWriter, r *http.Request) {

	log.Println("GameCreate started")
	addResponseHeaders(w)

	var board = newGameBoard()

	board.SaveGameBoard()

	log.Println("GameCreate finshed")
	returnBoardAsJson(w, board)
}

func GameById(w http.ResponseWriter, r *http.Request) {

	addResponseHeaders(w)

	vars := mux.Vars(r)
	gameId := vars["gameId"]

	board, err := models.LoadGameBoard(gameId)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	returnBoardAsJson(w, board)

}

func returnBoardAsJson(w http.ResponseWriter, board *models.GameBoard) {

	json.NewEncoder(w).Encode(&board)

}

func MovePlayerRight(w http.ResponseWriter, r *http.Request) {

	addResponseHeaders(w)

	// fetch latest board
	vars := mux.Vars(r)
	gameId := vars["gameId"]

	fmt.Println("Getting game board", gameId)

	var board, err = models.LoadGameBoard(gameId)

	if board == nil || err != nil {
		http.NotFound(w, r)
	}

	// move player right
	board.MainPlayer.Location.X++

	fmt.Println("Save game board", gameId)
	err = board.SaveGameBoard()

	if err != nil {
		fmt.Println(err)
	}

	returnBoardAsJson(w, board)
}

// received MainPlayer as JSON request
func UpdatePlayer(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Update player started")
	addResponseHeaders(w)

	jsonBody, err := getRequestBody(r)
	if err != nil {
		http.Error(w, "Failed to get request body", http.StatusInternalServerError)
		return
	}

	// unmarshall Player request
	mainPlayer, err := unmarshallPlayer(jsonBody)

	if err != nil {
		http.Error(w, "Failed to unmarshall player", http.StatusInternalServerError)
		return
	}

	// fetch current board
	vars := mux.Vars(r)
	gameId := vars["gameId"]

	fmt.Println("Getting game board", gameId)
	board, err := models.LoadGameBoard(gameId)

	if board == nil || err != nil {
		http.NotFound(w, r)
		return
	}

	// update board with player

	board.MainPlayer.Location = mainPlayer.Location

	// move player right
	//board.MainPlayer.Location.X++

	fmt.Println("Save game board", gameId)
	err = board.SaveGameBoard()

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	returnBoardAsJson(w, board)
}

func unmarshallPlayer(jsonBody []byte) (*models.Player, error) {

	var mainPlayer models.Player

	err := json.Unmarshal(jsonBody, &mainPlayer)

	return &mainPlayer, err

}
