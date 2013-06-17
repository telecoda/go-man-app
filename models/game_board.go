package models

type Point struct {
	X, Y int
}

type PlayerType int

const (
	GoMan = iota
	GoGhost
)

type PlayerState int

const (
	Alive = iota
	Dead
	Spawing
)

type Player struct {
	Location Point
	Id       int
	Type     PlayerType
}

type GameBoard struct {
	Id         string
	Name       string
	MainPlayer Player
	BoardCells [][]byte
}

const BOARD_WIDTH int = 28
const BOARD_HEIGHT int = 24

const PLAYER_START_X = 13
const PLAYER_START_Y = 14

var persister = NewFilePersister()

func (model *GameBoard) SaveGameBoard() {
	persister.Save(model)
}

func LoadGameBoard(id string) (*GameBoard, error) {
	return persister.Load(id)
}
