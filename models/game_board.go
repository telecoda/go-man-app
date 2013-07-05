package models

import (
	"bufio"
	"fmt"
	"github.com/telecoda/go-man/utils"
	"os"
	"time"
)

type Point struct {
	X, Y int
}

type GameState string

const (
	WaitingForPlayers GameState = "waiting"
	PlayingGame                 = "playing"
	BoardClear                  = "clear"
	GameOver                    = "over"
)

type GameBoard struct {
	Id                 string
	Name               string
	PillsRemaining     int
	Score              int
	Lives              int
	Players            []Player
	MaxGoMenAllowed    int
	MaxGoGhostsAllowed int
	State              GameState
	CreatedTime        time.Time
	LastUpdatedTime    time.Time
	BoardCells         [][]rune
}

// dimensions
const BOARD_WIDTH int = 28
const BOARD_HEIGHT int = 24

// cell types
const WALL = '#'
const PILL = '.'
const POWER_PILL = 'P'
const BONUS = '$'

// points
const PILL_POINTS = 10

var defaultBoard [][]rune

var persister = NewFilePersister()

func (board *GameBoard) SaveGameBoard() error {
	board.LastUpdatedTime = time.Now()
	return persister.Save(board)
}

func LoadGameBoard(id string) (*GameBoard, error) {
	return persister.Load(id)
}

func (board *GameBoard) DestroyGameBoard() error {
	fmt.Println("Destroying gameBoard:", board.Id)
	return persister.Destroy(board.Id)
}

func (board *GameBoard) eatPillAtLocation(location *Point) {
	board.Score += PILL_POINTS
	board.PillsRemaining--
	board.ClearCellAtLocation(location)
}

func (board *GameBoard) GetCellAtLocation(checkLocation *Point) rune {

	return board.BoardCells[checkLocation.Y][checkLocation.X]

}

func (board *GameBoard) ClearCellAtLocation(checkLocation *Point) {

	board.BoardCells[checkLocation.Y][checkLocation.X] = ' '

}

func (board *GameBoard) UpdatePillsRemaining() {
	count := 0
	for _, row := range board.BoardCells {
		for _, cell := range row {
			if cell == PILL || cell == POWER_PILL {
				count++
			}
		}
	}

	board.PillsRemaining = count

}

func initGameBoard() {

	defaultBoard = make([][]rune, BOARD_HEIGHT)

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
			defaultBoard[r] = make([]rune, BOARD_WIDTH)
			for c, cell := range row {
				defaultBoard[r][c] = rune(cell)
				c++
			}
			r++
		} else {
			break
		}

	}

}

func NewGameBoard() *GameBoard {
	initGameBoard()
	gameBoard := new(GameBoard)

	id, err := utils.GenUUID()
	if err != nil {
		fmt.Println("Error generating guid")
		return nil
	}
	gameBoard.Id = id
	gameBoard.Name = "Init name"
	gameBoard.Score = 0
	gameBoard.Lives = 3
	gameBoard.BoardCells = defaultBoard
	gameBoard.State = WaitingForPlayers
	gameBoard.MaxGoGhostsAllowed = MAX_GOMAN_GHOSTS
	gameBoard.MaxGoMenAllowed = MAX_GOMAN_PLAYERS
	gameBoard.CreatedTime = time.Now()
	gameBoard.UpdatePillsRemaining()

	// init players
	//gameBoard.MainPlayer = *newPlayer()

	return gameBoard
}
