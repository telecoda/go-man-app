package models

import (
	"fmt"
	"github.com/telecoda/go-man/utils"
	"log"
	"testing"
)

func init() {
	// delete all previous games
	utils.DeleteOldGameBoardFiles()
}

func TestCreateBoard(t *testing.T) {

	fmt.Println("TestCreateBoard started")

	board := NewGameBoard()

	if board == nil {
		log.Fatal("NewGameBoard failed to create a gameBoard")
	}

	// new board should be at state waiting for players
	if board.State != WaitingForPlayers {
		log.Fatal("A new game board should start as waiting for players")
	}

	// check player count
	if board.MaxGoMenAllowed != MAX_GOMAN_PLAYERS {
		log.Fatal("Max goman players not correct")
	}

	// check ghost count
	if board.MaxGoGhostsAllowed != MAX_GOMAN_GHOSTS {
		log.Fatal("Max goman ghosts not correct")
	}

	err := board.SaveGameBoard()

	if err != nil {
		log.Fatal("SaveGameBoard failed:", err)
	}

	fmt.Println("TestCreateBoard ended")

}

func TestAddGoManPlayer(t *testing.T) {

	fmt.Println("TestAddGoManPlayer started")

	board := NewGameBoard()

	if board == nil {
		log.Fatal("NewGameBoard failed to create a gameBoard")
	}

	// new board should be at state waiting for players
	if board.State != WaitingForPlayers {
		log.Fatal("A new game board should start as waiting for players")
	}

	// check player count
	if board.MaxGoMenAllowed != MAX_GOMAN_PLAYERS {
		log.Fatal("Max goman players not correct")
	}

	// check ghost count
	if board.MaxGoGhostsAllowed != MAX_GOMAN_GHOSTS {
		log.Fatal("Max goman ghosts not correct")
	}

	err := board.SaveGameBoard()

	if err != nil {
		log.Fatal("SaveGameBoard failed:", err)
	}

	err = board.DestroyGameBoard()
	if err != nil {
		log.Fatal("DestroyGameBoard failed:", err)
	}

	fmt.Println("TestAddGoManPlayer ended")

}

func TestIsMoveValidWorksWithValidXMove(t *testing.T) {

	fmt.Println("TestIsMoveValidWorksWithValidXMove started")

	existingLocation := &Point{10, 10}
	newLocation := &Point{11, 10}

	if !isMoveValid(existingLocation, newLocation) {
		log.Fatal("isMoveValid should allow this move")
	}

	fmt.Println("TestIsMoveValidWorksWithValidXMove ended")

}

func TestIsMoveValidFailsWithInvalidXMove(t *testing.T) {

	fmt.Println("TestIsMoveValidFailsWithInvalidXMove started")

	existingLocation := &Point{10, 10}
	newLocation := &Point{13, 10}

	if isMoveValid(existingLocation, newLocation) {
		log.Fatal("isMoveValid should NOT allow this move")
	}

	fmt.Println("TestIsMoveValidFailsWithInvalidXMove ended")

}

func TestIsMoveValidWorksWithValidYMove(t *testing.T) {

	fmt.Println("TestIsMoveValidWorksWithValidXMove started")

	existingLocation := &Point{10, 10}
	newLocation := &Point{10, 11}

	if !isMoveValid(existingLocation, newLocation) {
		log.Fatal("isMoveValid should allow this move")
	}

	fmt.Println("TestIsMoveValidWorksWithValidYMove ended")

}

func TestIsMoveValidFailsWithInvalidYMove(t *testing.T) {

	fmt.Println("TestIsMoveValidFailsWithInvalidYMove started")

	existingLocation := &Point{10, 10}
	newLocation := &Point{10, 7}

	if isMoveValid(existingLocation, newLocation) {
		log.Fatal("isMoveValid should NOT allow this move")
	}

	fmt.Println("TestIsMoveValidFailsWithInvalidYMove ended")

}

func TestIsMoveValidFailsWithInvalidXYMove(t *testing.T) {

	fmt.Println("TestIsMoveValidFailsWithInvalidXYMove started")

	existingLocation := &Point{10, 10}
	newLocation := &Point{11, 11}

	if isMoveValid(existingLocation, newLocation) {
		log.Fatal("isMoveValid should NOT allow this move")
	}

	fmt.Println("TestIsMoveValidFailsWithInvalidXYMove ended")

}
