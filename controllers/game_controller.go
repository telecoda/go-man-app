package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/telecoda/go-man-app/models"
	"log"
	"net/http"
)

func GameList(w http.ResponseWriter, r *http.Request) {
	addResponseHeaders(w)

	stateFilter := r.URL.Query().Get("state")
	fmt.Println("Query parameters:", stateFilter)

	// get all games
	boards, err := models.ReadAllGameBoards(stateFilter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	returnBoardsSummaryAsJson(w, boards)
}

func GameCreate(w http.ResponseWriter, r *http.Request) {

	log.Println("GameCreate started")
	addResponseHeaders(w)

	jsonBody, err := getRequestBody(r)
	if err != nil {
		http.Error(w, "Failed to get request body", http.StatusBadRequest)
		return
	}

	// unmarshall create game request
	newGame, err := unmarshallGameBoardSummary(jsonBody)

	if err != nil {
		http.Error(w, "Failed to unmarshall game"+err.Error(), http.StatusBadRequest)
		return
	}

	var board *models.GameBoard
	board, err = models.NewGameBoard(newGame.Name, newGame.MaxGoMenAllowed, newGame.MaxGoGhostsAllowed, newGame.WaitForPlayersSeconds)
	if err != nil {
		http.Error(w, "Failed to create a new game"+err.Error(), http.StatusBadRequest)
		return
	}

	board.CreateGameBoard()

	log.Println("GameCreate finshed")
	returnBoardAsJson(w, board)
}

func GameById(w http.ResponseWriter, r *http.Request) {

	addResponseHeaders(w)

	vars := mux.Vars(r)
	gameId := vars["gameId"]

	board, err := models.LoadGameBoard(gameId)

	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}

	returnBoardAsJson(w, board)

}

func returnBoardsSummaryAsJson(w http.ResponseWriter, boardsSummary *[]models.GameBoardSummary) {

	json.NewEncoder(w).Encode(&boardsSummary)

}

func returnBoardAsJson(w http.ResponseWriter, board *models.GameBoard) {

	json.NewEncoder(w).Encode(&board)

}

func returnPlayerAsJson(w http.ResponseWriter, player *models.Player) {

	json.NewEncoder(w).Encode(&player)

}

// received MainPlayer as JSON request
func AddPlayer(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Add player started")
	addResponseHeaders(w)

	jsonBody, err := getRequestBody(r)
	if err != nil {
		http.Error(w, "Failed to get request body", http.StatusBadRequest)
		return
	}

	// unmarshall Player request
	player, err := unmarshallPlayer(jsonBody)

	if err != nil {
		http.Error(w, "Failed to unmarshall player"+err.Error(), http.StatusBadRequest)
		return
	}

	// fetch current board
	vars := mux.Vars(r)
	gameId := vars["gameId"]

	fmt.Println("Getting game board", gameId)
	board, err := models.LoadGameBoard(gameId)

	if board == nil {
		http.NotFound(w, r)
		return
	}

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	player, err = board.AddPlayer(player)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Save game board", gameId)
	err = board.SaveGameBoard()

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	returnPlayerAsJson(w, player)

}

// received MainPlayer as JSON request
func ConcurrentUpdatePlayer(w http.ResponseWriter, r *http.Request) {

	addResponseHeaders(w)

	jsonBody, err := getRequestBody(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to get request body", http.StatusInternalServerError)
		return
	}

	// unmarshall Player request
	player, err := unmarshallPlayer(jsonBody)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to unmarshall player", http.StatusInternalServerError)
		return
	}

	// fetch current board
	vars := mux.Vars(r)
	gameId := vars["gameId"]
	//playerId := vars["playerId"]

	board, err := models.LoadGameBoard(gameId)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if board == nil {
		fmt.Println("Board not found:", gameId)
		http.NotFound(w, r)
		return
	}

	playerMoveRequest := new(models.PlayerMove)
	playerMoveRequest.GameId = gameId
	playerMoveRequest.PlayerToMove = *player

	playerResponseChannel := make(chan models.PlayerMoveResponse)

	playerMoveRequest.ResponseChannel = playerResponseChannel

	// send request to game channel
	var gameRequestChannel chan models.PlayerMove
	gameRequestChannel = models.GameChannels[gameId]

	if gameRequestChannel == nil {
		fmt.Println("Error no request channel found for game")
		return
	}

	// send
	gameRequestChannel <- *playerMoveRequest

	// receive response
	var playerMoveResponse models.PlayerMoveResponse

	playerMoveResponse = <-playerResponseChannel

	if playerMoveResponse.Error != nil {
		fmt.Println(playerMoveResponse.Error)
		http.Error(w, playerMoveResponse.Error.Error(), http.StatusBadRequest)
		return
	}

	/*err = playerMoveResponse.Board.SaveGameBoard()

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}*/

	returnBoardAsJson(w, &playerMoveResponse.Board)

}

// received MainPlayer as JSON request
func UpdatePlayer(w http.ResponseWriter, r *http.Request) {

	addResponseHeaders(w)

	jsonBody, err := getRequestBody(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to get request body", http.StatusInternalServerError)
		return
	}

	// unmarshall Player request
	player, err := unmarshallPlayer(jsonBody)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to unmarshall player", http.StatusInternalServerError)
		return
	}

	// fetch current board
	vars := mux.Vars(r)
	gameId := vars["gameId"]
	//playerId := vars["playerId"]

	board, err := models.LoadGameBoard(gameId)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if board == nil {
		fmt.Println("Board not found:", gameId)
		http.NotFound(w, r)
		return
	}

	err = board.MovePlayer(*player)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = board.SaveGameBoard()

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	returnBoardAsJson(w, board)

}

func unmarshallPlayer(jsonBody []byte) (*models.Player, error) {

	var player models.Player

	err := json.Unmarshal(jsonBody, &player)

	return &player, err

}

func unmarshallGameBoardSummary(jsonBody []byte) (*models.GameBoardSummary, error) {

	var board models.GameBoardSummary

	err := json.Unmarshal(jsonBody, &board)

	return &board, err

}
