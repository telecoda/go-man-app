package models

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/telecoda/go-man/utils"
	"math"
	"os"
)

type Point struct {
	X, Y int
}

type GameState int

const (
	WaitingForPlayers GameState = iota
	PlayingGame
	MainPlayerDied
	BoardClear
	GameOver
)

type PlayerType int

const (
	GoMan PlayerType = iota
	GoGhost
)

type PlayerState int

const (
	Alive PlayerState = iota
	Dead
	Spawing
)

type Player struct {
	Location      Point
	Id            string
	Type          PlayerType
	State         PlayerState
	Name          string
	cpuControlled bool
}

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
	BoardCells         [][]rune
}

// dimensions
const BOARD_WIDTH int = 28
const BOARD_HEIGHT int = 24

const PLAYER_START_X = 13
const PLAYER_START_Y = 14

const MAX_GOMAN_PLAYERS int = 1
const MAX_GOMAN_GHOSTS int = 4

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
	return persister.Save(board)
}

func LoadGameBoard(id string) (*GameBoard, error) {
	return persister.Load(id)
}

func (board *GameBoard) DestroyGameBoard() error {
	fmt.Println("Destroying gameBoard:", board.Id)
	return persister.Destroy(board.Id)
}

func (board *GameBoard) MovePlayer(player *Player) error {

	// only allow moves when game playing
	if board.State != PlayingGame {
		return errors.New("Not ready, please wait")
	}

	// check if player belongs to this game
	playerServerState := board.getPlayerFromServer(player.Id)
	if playerServerState == nil {
		return errors.New("You are not a player in this game.")
	}

	// check move is valid
	if !isMoveValid(&playerServerState.Location, &player.Location) {
		return errors.New("Cheat, invalid move")
	}

	cell := board.GetCellAtLocation(&player.Location)

	switch cell {
	case WALL:
		return errors.New("Invalid move, you can't walk through walls")
	case PILL:
		board.eatPillAtLocation(&player.Location)
		break
	}

	// update board with player
	playerServerState.Location = player.Location

	return nil
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

func (board *GameBoard) getPlayerFromServer(playerId string) *Player {

	for _, player := range board.Players {
		if player.Id == playerId {
			return &player
		}
	}

	return nil
}

func isMoveValid(existingLocation *Point, newLocation *Point) bool {

	// player can only move in one direction at a time
	// player can only move one cell at a time

	distX := math.Abs(float64(existingLocation.X - newLocation.X))
	distY := math.Abs(float64(existingLocation.Y - newLocation.Y))

	// moved more than one cell
	if distX > 1 || distY > 1 {
		return false
	}

	// moved more than one direction
	if distX > 0 && distY > 0 {
		return false
	}

	// valid move
	return true
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
	//gameBoard.Players = make([]Player, MAX_GOMAN_GHOSTS+MAX_GOMAN_PLAYERS)

	gameBoard.UpdatePillsRemaining()

	// init players
	//gameBoard.MainPlayer = *newPlayer()

	return gameBoard
}

func (board *GameBoard) AddPlayer(newPlayer *Player) (*Player, error) {

	if newPlayer.Type != GoMan && newPlayer.Type != GoGhost {
		return nil, errors.New("Invalid player type")
	}

	ghostCount := board.countGhosts()
	goMenCount := board.countGoMen()

	if newPlayer.Type == GoGhost && ghostCount >= MAX_GOMAN_GHOSTS {
		return nil, errors.New("Cannot add anymore ghosts to game")
	}

	if newPlayer.Type == GoMan && goMenCount >= MAX_GOMAN_PLAYERS {
		return nil, errors.New("Cannot add anymore go-men to game")
	}

	newPlayer.Location = Point{PLAYER_START_X, PLAYER_START_Y}
	newPlayer.Id, _ = utils.GenUUID()
	newPlayer.State = Alive
	board.Players = append(board.Players, *newPlayer)

	return newPlayer, nil
}

func (board *GameBoard) countGhosts() int {
	totalGhosts := 0
	for _, player := range board.Players {
		if player.Type == GoGhost {
			totalGhosts++
		}
	}

	return totalGhosts
}

func (board *GameBoard) countGoMen() int {
	totalGoMen := 0
	for _, player := range board.Players {
		if player.Type == GoMan {
			totalGoMen++
		}
	}

	return totalGoMen
}
