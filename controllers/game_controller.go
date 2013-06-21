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

var defaultBoard [][]rune

func AddResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	// allow cross origin requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func GameList(w http.ResponseWriter, r *http.Request) {
	AddResponseHeaders(w)
	fmt.Fprint(w, Response{"success": true, "message": "Here are the current games", "method": r.Method})
}

func GameCreate(w http.ResponseWriter, r *http.Request) {

	AddResponseHeaders(w)
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

	defaultBoard = make([][]rune, models.BOARD_HEIGHT)

	// read data from maze.dat
	f, err := os.Open("data/maze.txt")
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}
	reader := bufio.NewReader(f)

	var r int = 0
	for {

		b, err := reader.ReadBytes('\n')
		if err == nil {
			// parse line

			b = b[:len(b)-1] // remove last new line char from bytes
			row := string(b)
			fmt.Println("Processing row:", r, row)
			defaultBoard[r] = make([]rune, models.BOARD_WIDTH)
			for c, cell := range row {
				fmt.Println("Cell:", c, cell)
				defaultBoard[r][c] = rune(cell)
				c++
				//fmt.Println(defaultBoard[r])

			}
			r++
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

	AddResponseHeaders(w)

	vars := mux.Vars(r)
	gameId := vars["gameId"]

	fmt.Println("Getting game board", gameId)
	var board, err = models.LoadGameBoard(gameId)

	//fmt.Println("Loaded board", board)

	bJson, err := json.Marshal(board)

	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Fprint(w, string(bJson))
	}
}

func UpdatePlayer(w http.ResponseWriter, r *http.Request) {

	AddResponseHeaders(w)

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
