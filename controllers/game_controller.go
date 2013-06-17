package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/telecoda/go-man/models"
	"github.com/telecoda/go-man/utils"
	"net/http"
	"os"
)

/* Thanks to the following source for an ASCII version of the game board
http://4coder.org/c-c-source-code/152/pacman/board.c.html

*/

var defaultBoard [][]byte

func GameList(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, Response{"success": true, "message": "Here are the current games", "method": r.Method})
}

func GameCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var board = NewGameBoard()

	board.SaveGameBoard()

	bJson, err := json.Marshal(board)

	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Fprint(w, string(bJson))
	}
}

func init() {
	initGameBoard()
}

func initGameBoard() {

	defaultBoard = make([][]byte, models.BOARD_HEIGHT)

	// read data from maze.dat
	f, err := os.Open("data/maze.txt")
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}
	r := bufio.NewReader(f)

	var i int = 0
	for {

		b, err := r.ReadBytes('\n')
		if err == nil {
			// parse line

			b = b[:len(b)-1] // remove last new line char from bytes
			defaultBoard[i] = b
			fmt.Println(string(defaultBoard[i]))
			i++
		} else {
			break
		}

	}

}

func NewGameBoard() *models.GameBoard {
	initGameBoard()
	gameBoard := new(models.GameBoard)

	id, err := utils.GenUUID()
	if err != nil {
		fmt.Println("Error generating guid")
		return nil
	}
	gameBoard.Id = id
	gameBoard.Name = "Init name"
	gameBoard.BoardCells = defaultBoard

	// init players
	gameBoard.MainPlayer = *NewPlayer()

	return gameBoard
}

func NewPlayer() *models.Player {
	//return &models.Player{Location: {0, 0}, Id: 1, Type: models.PlayerType.GoMan}
	//player := models.Player{
	//	Location: models.Point{models.PLAYER_START_X, models.PLAYER_START_Y},
	//	Id:       1,
	//}

	player := new(models.Player)
	player.Location = models.Point{models.PLAYER_START_X, models.PLAYER_START_Y}

	player.Id, _ = utils.GenUUID()

	return player
}

func GameById(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	gameId := vars["gameId"]

	fmt.Println("Getting game board", gameId)
	var board, err = models.LoadGameBoard(gameId)

	fmt.Println("Loaded board", board)

	bJson, err := json.Marshal(board)

	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Fprint(w, string(bJson))
	}
}

func UpdatePlayer(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	// fetch latest board
	vars := mux.Vars(r)
	gameId := vars["gameId"]
	// get player from body of PUT request
	player := new(models.Player)

	fmt.Println("Getting game board", gameId)
	var board, err = models.LoadGameBoard(gameId)

	if board == nil || err != nil {
		http.NotFound(w, r)
	}

	// verify play belongs to this gameboard
	if playerInGame(board, player) {

	}

}

func playerInGame(board *models.GameBoard, player *models.Player) bool {
	return true
}
