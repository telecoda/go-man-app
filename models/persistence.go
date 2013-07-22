package models

type GameBoardPersistence interface {
	Create(model *GameBoard) (err error)
	Read(id string) (model *GameBoard, err error)
	ReadAll() ([]GameBoard, error) 
	Update(model *GameBoard) (err error)
	Delete(id string) (model *GameBoard, err error)
}

type Persister struct {
	Name string
}