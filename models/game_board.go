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
	NewGame           GameState = "new"
	WaitingForPlayers           = "waiting"
	PlayingGame                 = "playing"
	GameWon                     = "won"
	GameOver                    = "over"
)

type GameBoard struct {
	Id                 string
	Name               string
	PillsRemaining     int
	Players            map[string]*Player
	MaxGoMenAllowed    int
	MaxGoGhostsAllowed int
	WaitForPlayersSeconds int
	State              GameState
	PowerPillsActive   int
	CreatedTime        time.Time
	LastUpdatedTime    time.Time
	GameStartTime      time.Time
	BoardCells         [][]rune
}

type GameBoardSummary struct {
	Id                 string
	Name               string
	Players            map[string]*Player
	MaxGoMenAllowed    int
	MaxGoGhostsAllowed int
	WaitForPlayersSeconds int
	State              GameState
	CreatedTime        time.Time
	LastUpdatedTime    time.Time
	GameStartTime      time.Time
}

// dimensions
const BOARD_WIDTH int = 28
const BOARD_HEIGHT int = 24

const INITIAL_GAMES_HOSTED = 10

// cell types
const WALL = '#'
const PILL = '.'
const POWER_PILL = 'P'
const BONUS = '$'

// points
const PILL_POINTS = 10
const POWER_PILL_POINTS = 50
const POWER_PILL_WAIT_SECONDS = 5

var GamePersister = InMemoryPersister()

func (board *GameBoard) CreateGameBoard() error {
	return GamePersister.Create(board)
}

func (board *GameBoard) SaveGameBoard() error {
	board.LastUpdatedTime = time.Now()
	return GamePersister.Update(board)
}

func LoadGameBoard(id string) (*GameBoard, error) {
	return GamePersister.Read(id)
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
	boardSummary.GameStartTime = board.GameStartTime

	return boardSummary
}

func ReadAllGameBoards(filterByState string) (*[]GameBoardSummary, error) {

	boards, err := GamePersister.ReadAll()

	if err != nil {
		return nil, err
	}

	fmt.Println("FilterByState:", filterByState)

	// at least return an empty list
	var boardSummaries = make([]GameBoardSummary, 0)

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
	return GamePersister.Delete(board.Id)
}

func (board *GameBoard) eatPillAtLocation(location Point) {
	board.PillsRemaining--
	board.ClearCellAtLocation(location)
}

func (board *GameBoard) eatPowerPillAtLocation(location Point) {
	go board.activatePowerPill()
	// start power pill timer...
	board.PillsRemaining--
	board.ClearCellAtLocation(location)
}

func (board *GameBoard) activatePowerPill() {
	// increment active powerpill count
	// this is to gracefully cope with multiple activations
	board.PowerPillsActive++

	// wait till end of pill
	time.Sleep(time.Duration(POWER_PILL_WAIT_SECONDS) * time.Second)

	board.PowerPillsActive--

}

func (board *GameBoard) GetCellAtLocation(checkLocation Point) rune {

	return board.BoardCells[checkLocation.Y][checkLocation.X]

}

func (board *GameBoard) ClearCellAtLocation(checkLocation Point) {

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

func NewGameBoard(gameName string, maxGoMen int, maxGoGhosts int, waitForPlayersSeconds int) *GameBoard {
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
	gameBoard.Name = gameName
	gameBoard.BoardCells = defaultBoard
	gameBoard.State = NewGame
	gameBoard.WaitForPlayersSeconds = waitForPlayersSeconds
	gameBoard.MaxGoMenAllowed = maxGoMen
	gameBoard.MaxGoGhostsAllowed = maxGoGhosts
	gameBoard.Players = make(map[string]*Player)
	gameBoard.CreatedTime = time.Now()
	gameBoard.GameStartTime = gameBoard.CreatedTime.Add(time.Duration(waitForPlayersSeconds) * time.Second)
	gameBoard.UpdatePillsRemaining()

	return gameBoard
}
