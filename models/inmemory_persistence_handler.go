package models

import (
	"errors"
	"fmt"
)

var games map[string]*GameBoard

func initGames() {

	games = make(map[string]*GameBoard)

}

func InMemoryPersister() *Persister {
	initGames()
	return &Persister{Name: "In Memory Persister"}
}

func (p *Persister) Create(board *GameBoard) (err error) {

	// create an in memory instance of the game
	games[board.Id] = board
	return p.Update(board)
}

func (p *Persister) Update(board *GameBoard) (err error) {

	// create an in memory instance of the game
	games[board.Id] = board

	return nil

}

func (p *Persister) Read(id string) (*GameBoard, error) {

	board := games[id]

	if board == nil {
		fmt.Println("Board not found")
		fmt.Println("Total boards:", len(games))
		return nil, errors.New("Board not found")
	}

	return board, nil
}

func (p *Persister) ReadAll() ([]GameBoard, error) {

	fmt.Println("Loading gameboards:")

	var boards []GameBoard

	for id, _ := range games {

		board := games[id]
		boards = append(boards, *board)

	}

	return boards, nil

}

func (p *Persister) Delete(id string) error {

	fmt.Println("Destroying gameboard:", id)

	delete(games, id)

	return nil
}

func (p *Persister) DeleteAll() error {

	initGames()

	return nil
}
