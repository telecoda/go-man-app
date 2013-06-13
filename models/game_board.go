package models

import (
	"bufio"
	"fmt"
	"github.com/telecoda/go-man/utils"
	"os"
)

type GameBoard struct {
	Id         string
	Name       string
	BoardCells [][]byte
}

const BOARD_WIDTH int = 28
const BOARD_HEIGHT int = 24

/* Thanks to the following source for an ASCII version of the game board
http://4coder.org/c-c-source-code/152/pacman/board.c.html

*/

var defaultBoard [][]byte

var persister = NewFilePersister()

func init() {
	initGameBoard()
}

func initGameBoard() {

	defaultBoard = make([][]byte, BOARD_HEIGHT)

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

func NewGameBoard() *GameBoard {
	gameBoard := new(GameBoard)

	id, err := utils.GenUUID()
	if err != nil {
		fmt.Println("Error generating guid")
		return nil
	}
	gameBoard.Id = id
	gameBoard.Name = "Init name"
	gameBoard.BoardCells = defaultBoard

	return gameBoard
}

func (model *GameBoard) SaveGameBoard() {
	persister.Save(model)
}

func LoadGameBoard(id string) (*GameBoard, error) {
	return persister.Load(id)
}
