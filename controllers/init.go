package controllers

import (
	"bufio"
	"fmt"
	"github.com/telecoda/go-man/models"
	"github.com/telecoda/go-man/utils"
	"os"
)

var defaultBoard [][]rune

func init() {
	//initGameBoard()
}

func initGameBoard() {

	defaultBoard = make([][]rune, models.BOARD_HEIGHT)

	// this path will be of the controllers package
	filePath := utils.GetAbsolutePathOfCurrentPackage("../data/maze.txt")

	f, err := os.Open(filePath)
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
			//fmt.Println("Processing row:", r, row)
			defaultBoard[r] = make([]rune, models.BOARD_WIDTH)
			for c, cell := range row {
				//fmt.Println("Cell:", c, cell)
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

func newGameBoard() *models.GameBoard {
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
	gameBoard.MainPlayer = *newPlayer()

	return gameBoard
}

func newPlayer() *models.Player {
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
