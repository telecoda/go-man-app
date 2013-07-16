package models

import (
	"encoding/json"
	"fmt"
	"github.com/telecoda/go-man/utils"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

	filePath := utils.GetAbsolutePathOfCurrentPackage("../gamedata/" + board.Id + ".json")

	// convert to JSON for saving to file (binary would be quicker...)
	bJson, err := json.Marshal(board)

	err = ioutil.WriteFile(filePath, bJson, 0600)

	if err != nil {
		fmt.Println("Error saving file", err)
		return err
	}

	fmt.Println("Saved gameboard")
	return nil

}

func (p *Persister) Load(id string) (*GameBoard, error) {

	filePath := utils.GetAbsolutePathOfCurrentPackage("../gamedata/" + id + ".json")

	bJson, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error loading file", err)
		return nil, err
	}

	// convert from JSON to object to probably conver back to JSON...

	var board GameBoard

	err = json.Unmarshal(bJson, &board)

	return &board, err
}

func (p *Persister) LoadAll() ([]GameBoard, error) {

	fmt.Println("Loading gameboards:")

	filePath := utils.GetAbsolutePathOfCurrentPackage("../gamedata/")

	dirs, err := ioutil.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var boards []GameBoard

	for _, dir := range dirs {
		name := dir.Name()
		if !dir.IsDir() && strings.HasSuffix(name, ".json") {
			id := name[:len(name)-5] // remove .json from filename
			board, err := LoadGameBoard(id)
			if err != nil {
				return nil, err
			}
			boards = append(boards, *board)
		}
	}

	fmt.Println("Loaded gameboards")
	return boards, err
}

func (p *Persister) Destroy(id string) error {

	fmt.Println("Destroying gameboard:", id)

	filePath := utils.GetAbsolutePathOfCurrentPackage("../gamedata/" + id + ".json")

	return os.Remove(filePath)
}
