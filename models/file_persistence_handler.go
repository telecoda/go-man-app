package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type GameBoardPersistence interface {
	Save(model *GameBoard) (err error)
	Load(id string) (model *GameBoard, err error)
}

type Persister struct {
	Name string
}

func NewFilePersister() *Persister {
	return &Persister{Name: "File Persister"}
}

// Match matches registered routes against the request.
func (p *Persister) Save(board *GameBoard) (err error) {

	fmt.Println("Saving gameboard as JSON:", board.Id)

	filename := "gamedata/" + board.Id + ".json"

	// convert to JSON for saving to file (binary would be quicker...)
	bJson, err := json.Marshal(board)

	err = ioutil.WriteFile(filename, bJson, 0600)

	if err != nil {
		fmt.Println("Error saving file", err)
		return err
	}

	fmt.Println("Saved gameboard")
	return nil

}

func (p *Persister) Load(id string) (*GameBoard, error) {

	fmt.Println("Loading gameboard:", id)

	filename := "gamedata/" + id + ".json"

	bJson, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error saving file", err)
		return nil, err
	}

	// convert from JSON to object to probably conver back to JSON...

	var board GameBoard

	err = json.Unmarshal(bJson, &board)

	fmt.Println("Loaded gameboard")
	return &board, nil
}
