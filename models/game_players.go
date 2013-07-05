package models

import (
	"errors"
	"github.com/telecoda/go-man/utils"
	"math"
)

/* this file contains the player specific function */

type PlayerType string

const (
	GoMan   PlayerType = "goman"
	GoGhost            = "goghost"
)

type PlayerState string

const (
	Alive    PlayerState = "alive"
	Dead                 = "dead"
	Spawning             = "spawning"
)

type Player struct {
	Location      Point
	Id            string
	Type          PlayerType
	State         PlayerState
	Name          string
	cpuControlled bool
}

const PLAYER_START_X = 13
const PLAYER_START_Y = 14

const MAX_GOMAN_PLAYERS int = 1
const MAX_GOMAN_GHOSTS int = 4

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
