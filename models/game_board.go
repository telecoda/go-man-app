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
	Id       string
	Type     PlayerType
	State    PlayerState
}

type GameBoard struct {
	Id             string
	Name           string
	PillsRemaining int
	MainPlayer     Player
	BoardCells     [][]rune
}

const BOARD_WIDTH int = 28
const BOARD_HEIGHT int = 24

const PLAYER_START_X = 13
const PLAYER_START_Y = 14

// cell types
const WALL = '#'
const PILL = '.'
const POWER_PILL = 'P'
const BONUS = '$'

var persister = NewFilePersister()

func (model *GameBoard) SaveGameBoard() error {
	return persister.Save(model)
}

func LoadGameBoard(id string) (*GameBoard, error) {
	return persister.Load(id)
}

func (model *GameBoard) UpdatePillsRemaining() {
	count := 0
	for _, row := range model.BoardCells {
		for _, cell := range row {
			if cell == PILL || cell == POWER_PILL {
				count++
			}
		}
	}

	model.PillsRemaining = count

}
