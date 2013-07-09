package models

import (
	"fmt"
	"github.com/telecoda/go-man/utils"
	"time"
)

/* Thanks to the following source for an ASCII version of the game board
http://4coder.org/c-c-source-code/152/pacman/board.c.html
*/
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

type GameBoardSummary struct {
	Id                 string
	Name               string
	Players            []Player
	MaxGoMenAllowed    int
	MaxGoGhostsAllowed int
	State              GameState
	CreatedTime        time.Time
	LastUpdatedTime    time.Time
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

var persister = NewFilePersister()

func (board *GameBoard) SaveGameBoard() error {
	board.LastUpdatedTime = time.Now()
	return persister.Save(board)
}

func LoadGameBoard(id string) (*GameBoard, error) {
	return persister.Load(id)
}

func (board *GameBoard) convertToBoardSummary() *GameBoardSummary {

	boardSummary := new(GameBoardSummary)

	boardSummary.Id = board.Id
	boardSummary.Name = board.Name
	boardSummary.MaxGoGhostsAllowed = board.MaxGoGhostsAllowed
	boardSummary.MaxGoMenAllowed = board.MaxGoMenAllowed
	boardSummary.Players = board.Players
	boardSummary.State = board.State
	boardSummary.CreatedTime = board.CreatedTime
	boardSummary.LastUpdatedTime = board.LastUpdatedTime

	return boardSummary
}

func ReadAllGameBoards(filterByState string) (*[]GameBoardSummary, error) {

	boards, err := persister.LoadAll()

	if err != nil {
		return nil, err
	}

	fmt.Println("FilterByState:", filterByState)
	var boardSummaries []GameBoardSummary
	// convert boards to board summary
	for _, board := range boards {
		if filterByState != "" {
			if string(board.State) == filterByState {

				boardSummaries = append(boardSummaries, *board.convertToBoardSummary())
			}

		} else {
			boardSummaries = append(boardSummaries, *board.convertToBoardSummary())
		}
	}
	return &boardSummaries, nil
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

func NewGameBoard() *GameBoard {
	defaultBoard, err := initGameBoard()
	if err != nil {
		fmt.Println(err)
		return nil
	}
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

	return gameBoard
}
