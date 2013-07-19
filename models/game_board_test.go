package models

import (
	"fmt"
	"github.com/telecoda/go-man/utils"
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
		t.Errorf("NewGameBoard failed to create a gameBoard")
	}

	// new board should be at state NewGame
	if board.State != NewGame {
		t.Errorf("A new game board should start at state NewGame")
	}

	// check player count
	if board.MaxGoMenAllowed != MAX_GOMAN_PLAYERS {
		t.Errorf("Max goman players not correct")
	}

	// check ghost count
	if board.MaxGoGhostsAllowed != MAX_GOMAN_GHOSTS {
		t.Errorf("Max goman ghosts not correct")
	}

	err := board.SaveGameBoard()

	if err != nil {
		t.Errorf("SaveGameBoard failed:", err)
	}

	fmt.Println("TestCreateBoard ended")

}

func TestAddGoManPlayerWorksWithValidBoard(t *testing.T) {

	fmt.Println("TestAddGoManPlayerWorksWithValidBoard started")

	board := NewGameBoard()

	if board == nil {
		t.Errorf("NewGameBoard failed to create a gameBoard")
	}

	newPlayer := new(Player)
	newPlayer.Name = "Rob"
	newPlayer.Type = GoMan

	addedPlayer, err := board.AddPlayer(newPlayer)

	if err != nil {
		t.Errorf("Error adding player to board:", err.Error)
	}

	if addedPlayer == nil {
		t.Errorf("Failed to add player to game")
	}

	if addedPlayer.Id == "" {
		t.Errorf("Added player does not have id")
	}

	if addedPlayer.Name != "Rob" {
		t.Errorf("Player has wrong name")
	}

	if addedPlayer.Type != GoMan {
		t.Errorf("Player has wrong type")
	}

	if len(board.Players) != 1 {
		t.Errorf("Board should have 1 player")
	}

	fmt.Println("TestAddGoManPlayerWorksWithValidBoard ended")

}

func TestAddGoManPlayerFailsIfTooManyGoMen(t *testing.T) {

	fmt.Println("TestAddGoManPlayerFailsIfTooManyGoMen - started")

	board := NewGameBoard()

	if board == nil {
		t.Errorf("NewGameBoard failed to create a gameBoard")
	}

	newPlayer1 := new(Player)
	newPlayer1.Name = "Rob"
	newPlayer1.Type = GoMan

	addedPlayer1, err := board.AddPlayer(newPlayer1)

	if err != nil {
		t.Errorf("Error adding player to board:", err.Error)
	}

	if addedPlayer1 == nil {
		t.Errorf("Failed to add player to game")
	}

	newPlayer2 := new(Player)
	newPlayer2.Name = "Bob"
	newPlayer2.Type = GoMan

	addedPlayer2, err := board.AddPlayer(newPlayer2)

	fmt.Println("Error expected, here it is:", err)
	if err == nil {
		t.Errorf("Adding a second GoMan player SHOULD cause an error")
	}

	if addedPlayer2 != nil {
		t.Errorf("Second GoMan player should NOT be added")
	}

	fmt.Println("TestAddGoManPlayerFailsIfTooManyGoMen - ended")

}

func TestAddGoManPlayerFailsIfTooManyGoGhosts(t *testing.T) {

	fmt.Println("TestAddGoGhostFailsIfTooManyGoGhosts - started")

	board := NewGameBoard()

	if board == nil {
		t.Errorf("NewGameBoard failed to create a gameBoard")
	}

	newGhost1 := new(Player)
	newGhost1.Name = "Blinky"
	newGhost1.Type = GoGhost

	addedGhost1, err := board.AddPlayer(newGhost1)

	if err != nil {
		t.Errorf("Error adding player to board:", err.Error)
	}

	if addedGhost1 == nil {
		t.Errorf("Failed to add ghost to game")
	}

	newGhost2 := new(Player)
	newGhost2.Name = "Pinky"
	newGhost2.Type = GoGhost

	addedGhost2, err := board.AddPlayer(newGhost2)

	if err != nil {
		t.Errorf("Error adding player to board:", err.Error)
	}

	if addedGhost2 == nil {
		t.Errorf("Failed to add ghost to game")
	}

	newGhost3 := new(Player)
	newGhost3.Name = "Inky"
	newGhost3.Type = GoGhost

	addedGhost3, err := board.AddPlayer(newGhost3)

	if err != nil {
		t.Errorf("Error adding player to board:", err.Error)
	}

	if addedGhost3 == nil {
		t.Errorf("Failed to add ghost to game")
	}

	newGhost4 := new(Player)
	newGhost4.Name = "Clyde"
	newGhost4.Type = GoGhost

	addedGhost4, err := board.AddPlayer(newGhost4)

	if err != nil {
		t.Errorf("Error adding player to board:", err.Error)
	}

	if addedGhost4 == nil {
		t.Errorf("Failed to add ghost to game")
	}

	newGhost5 := new(Player)
	newGhost5.Name = "Earl"
	newGhost5.Type = GoGhost

	addedGhost5, err := board.AddPlayer(newGhost5)

	fmt.Println("Error expected, here it is:", err)
	if err == nil {
		t.Errorf("Adding a fifth GoGhost player SHOULD cause an error")
	}

	if addedGhost5 != nil {
		t.Errorf("Fifth GoGhost player should NOT be added")
	}

	fmt.Println("TestAddGoGhostFailsIfTooManyGoGhosts - ended")

}

func TestAddPlayerFailsWithInvalidType(t *testing.T) {

	fmt.Println("TestAddPlayerFailsWithInvalidType - started")

	board := NewGameBoard()

	if board == nil {
		t.Errorf("NewGameBoard failed to create a gameBoard")
	}

	newPlayer := new(Player)
	newPlayer.Name = "Joe"
	newPlayer.Type = "invalid" // use a non valid constant

	addedPlayer, err := board.AddPlayer(newPlayer)

	if err == nil {
		t.Errorf("Adding a player with an unknown type SHOULD return an error")
	}

	if addedPlayer != nil {
		t.Errorf("Player should NOT have been added")
	}

	fmt.Println("TestAddPlayerFailsWithInvalidType - ended")

}

func TestIsMoveValidWorksWithValidXMove(t *testing.T) {

	fmt.Println("TestIsMoveValidWorksWithValidXMove - started")

	existingLocation := Point{10, 10}
	newLocation := Point{11, 10}

	if !isMoveValid(existingLocation, newLocation) {
		t.Errorf("isMoveValid should allow this move")
	}

	fmt.Println("TestIsMoveValidWorksWithValidXMove - ended")

}

func TestIsMoveValidFailsWithInvalidXMove(t *testing.T) {

	fmt.Println("TestIsMoveValidFailsWithInvalidXMove - started")

	existingLocation := Point{10, 10}
	newLocation := Point{13, 10}

	if isMoveValid(existingLocation, newLocation) {
		t.Errorf("isMoveValid should NOT allow this move")
	}

	fmt.Println("TestIsMoveValidFailsWithInvalidXMove - ended")

}

func TestIsMoveValidWorksWithValidYMove(t *testing.T) {

	fmt.Println("TestIsMoveValidWorksWithValidXMove - started")

	existingLocation := Point{10, 10}
	newLocation := Point{10, 11}

	if !isMoveValid(existingLocation, newLocation) {
		t.Errorf("isMoveValid should allow this move")
	}

	fmt.Println("TestIsMoveValidWorksWithValidYMove - ended")

}

func TestIsMoveValidFailsWithInvalidYMove(t *testing.T) {

	fmt.Println("TestIsMoveValidFailsWithInvalidYMove - started")

	existingLocation := Point{10, 10}
	newLocation := Point{10, 7}

	if isMoveValid(existingLocation, newLocation) {
		t.Errorf("isMoveValid should NOT allow this move")
	}

	fmt.Println("TestIsMoveValidFailsWithInvalidYMove - ended")

}

func TestIsMoveValidFailsWithInvalidXYMove(t *testing.T) {

	fmt.Println("TestIsMoveValidFailsWithInvalidXYMove - started")

	existingLocation := Point{10, 10}
	newLocation := Point{11, 11}

	if isMoveValid(existingLocation, newLocation) {
		t.Errorf("isMoveValid should NOT allow this move")
	}

	fmt.Println("TestIsMoveValidFailsWithInvalidXYMove - ended")

}

func TestMovePlayerWithValidMoveWorks(t *testing.T) {

	setup()
	defer tearDown()

	fmt.Println("TestMovePlayerWithValidMoveWorks started")

	gameId := "game_at_playing_state"
	playerId := "goman_player"

	// fetch board
	board, err := LoadGameBoard(gameId)

	if err != nil {
		t.Errorf("Error fetching game:", err.Error)
	}

	if board == nil {
		t.Errorf("Error: GameBoard not returned")
	}

	player := board.getPlayerFromServer(playerId)

	if player == nil {
		t.Errorf("Player not found in game")
	}

	// move player right

	player.Location.X++

	board.MovePlayer(*player)
	// move player
	if err != nil {
		t.Errorf("Error moving player on board:", err.Error)
	}

	// fetch moved player from board
	movedPlayer := board.getPlayerFromServer(playerId)

	if movedPlayer == nil {
		t.Errorf("Moved Player not found in game")
	}

	// check player has actually moved
	if (movedPlayer.Location.X != player.Location.X) || (movedPlayer.Location.Y != player.Location.Y) {
		t.Errorf("Player has not moved")
	}

	fmt.Println("TestMovePlayerWithValidMoveWorks ended")

}

/* helper functions */

func setup() {
	fmt.Println("Test setup")
	deleteAllGames()
	copyGameFixtures()
}

func tearDown() {
	fmt.Println("Test teardown")
	deleteAllGames()
}

func deleteAllGames() {
	utils.DeleteOldGameBoardFiles()
}

func copyGameFixtures() {
	utils.CopyGameDataFixtures("gamedata_fixtures", "gamedata")
}
