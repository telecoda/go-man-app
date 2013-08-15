package models

import (
	"fmt"
	"testing"
)

var playingBoardId string
var playingBoardPlayerId string

func TestCreateBoard(t *testing.T) {

	fmt.Println("TestCreateBoard started")

	board, err := NewGameBoard("test game", 1, 4, 1)

	if err != nil {
		t.Errorf("NewGameBoard failed:", err)
		return
	}

	if board == nil {
		t.Errorf("NewGameBoard failed to create a gameBoard")
		return
	}

	// new board should be at state NewGame
	if board.State != NewGame {
		t.Errorf("A new game board should start at state NewGame")
		return
	}

	// check player count
	if board.MaxGoMenAllowed != 1 {
		t.Errorf("Max goman players not correct")
		return
	}

	// check ghost count
	if board.MaxGoGhostsAllowed != 4 {
		t.Errorf("Max goman ghosts not correct")
		return
	}

	err = board.SaveGameBoard()

	if err != nil {
		t.Errorf("SaveGameBoard failed:", err)
		return
	}

	fmt.Println("TestCreateBoard ended")

}

func TestAddGoManPlayerWorksWithValidBoard(t *testing.T) {

	fmt.Println("TestAddGoManPlayerWorksWithValidBoard started")

	board, err := NewGameBoard("test game", 1, 4, 1)

	if err != nil {
		t.Errorf("NewGameBoard failed:", err)
		return
	}

	if board == nil {
		t.Errorf("NewGameBoard failed to create a gameBoard")
		return
	}

	newPlayer := new(Player)
	newPlayer.Name = "Rob"
	newPlayer.Type = GoMan

	addedPlayer, err := board.AddPlayer(newPlayer)

	if err != nil {
		t.Errorf("Error adding player to board:", err.Error)
		return
	}

	if addedPlayer == nil {
		t.Errorf("Failed to add player to game")
		return
	}

	if addedPlayer.Id == "" {
		t.Errorf("Added player does not have id")
		return
	}

	if addedPlayer.Name != "Rob" {
		t.Errorf("Player has wrong name")
		return
	}

	if addedPlayer.Type != GoMan {
		t.Errorf("Player has wrong type")
		return
	}

	if len(board.Players) != 1 {
		t.Errorf("Board should have 1 player")
		return
	}

	fmt.Println("TestAddGoManPlayerWorksWithValidBoard ended")

}

func TestAddGoManPlayerFailsIfTooManyGoMen(t *testing.T) {

	fmt.Println("TestAddGoManPlayerFailsIfTooManyGoMen - started")

	board, err := NewGameBoard("test game", 1, 4, 1)

	if err != nil {
		t.Errorf("NewGameBoard failed:", err)
		return
	}

	if board == nil {
		t.Errorf("NewGameBoard failed to create a gameBoard")
		return
	}

	newPlayer1 := new(Player)
	newPlayer1.Name = "Rob"
	newPlayer1.Type = GoMan

	addedPlayer1, err := board.AddPlayer(newPlayer1)

	if err != nil {
		t.Errorf("Error adding player to board:", err.Error)
		return
	}

	if addedPlayer1 == nil {
		t.Errorf("Failed to add player to game")
		return
	}

	newPlayer2 := new(Player)
	newPlayer2.Name = "Bob"
	newPlayer2.Type = GoMan

	addedPlayer2, err := board.AddPlayer(newPlayer2)

	fmt.Println("Error expected, here it is:", err)
	if err == nil {
		t.Errorf("Adding a second GoMan player SHOULD cause an error")
		return
	}

	if addedPlayer2 != nil {
		t.Errorf("Second GoMan player should NOT be added")
		return
	}

	fmt.Println("TestAddGoManPlayerFailsIfTooManyGoMen - ended")

}

func TestAddGoManPlayerFailsIfTooManyGoGhosts(t *testing.T) {

	fmt.Println("TestAddGoGhostFailsIfTooManyGoGhosts - started")

	board, err := NewGameBoard("test game", 1, 4, 1)

	if err != nil {
		t.Errorf("NewGameBoard failed:", err)
		return
	}

	if board == nil {
		t.Errorf("NewGameBoard failed to create a gameBoard")
		return
	}

	newGhost1 := new(Player)
	newGhost1.Name = "Blinky"
	newGhost1.Type = GoGhost

	addedGhost1, err := board.AddPlayer(newGhost1)

	if err != nil {
		t.Errorf("Error adding player to board:", err.Error)
		return
	}

	if addedGhost1 == nil {
		t.Errorf("Failed to add ghost to game")
		return
	}

	newGhost2 := new(Player)
	newGhost2.Name = "Pinky"
	newGhost2.Type = GoGhost

	addedGhost2, err := board.AddPlayer(newGhost2)

	if err != nil {
		t.Errorf("Error adding player to board:", err.Error)
		return
	}

	if addedGhost2 == nil {
		t.Errorf("Failed to add ghost to game")
		return
	}

	newGhost3 := new(Player)
	newGhost3.Name = "Inky"
	newGhost3.Type = GoGhost

	addedGhost3, err := board.AddPlayer(newGhost3)

	if err != nil {
		t.Errorf("Error adding player to board:", err.Error)
		return
	}

	if addedGhost3 == nil {
		t.Errorf("Failed to add ghost to game")
		return
	}

	newGhost4 := new(Player)
	newGhost4.Name = "Clyde"
	newGhost4.Type = GoGhost

	addedGhost4, err := board.AddPlayer(newGhost4)

	if err != nil {
		t.Errorf("Error adding player to board:", err.Error)
		return
	}

	if addedGhost4 == nil {
		t.Errorf("Failed to add ghost to game")
		return
	}

	newGhost5 := new(Player)
	newGhost5.Name = "Earl"
	newGhost5.Type = GoGhost

	addedGhost5, err := board.AddPlayer(newGhost5)

	fmt.Println("Error expected, here it is:", err)
	if err == nil {
		t.Errorf("Adding a fifth GoGhost player SHOULD cause an error")
		return
	}

	if addedGhost5 != nil {
		t.Errorf("Fifth GoGhost player should NOT be added")
		return
	}

	fmt.Println("TestAddGoGhostFailsIfTooManyGoGhosts - ended")

}

func TestAddPlayerFailsWithInvalidType(t *testing.T) {

	fmt.Println("TestAddPlayerFailsWithInvalidType - started")

	board, err := NewGameBoard("test game", 1, 4, 1)

	if err != nil {
		t.Errorf("NewGameBoard failed:", err)
		return
	}

	if board == nil {
		t.Errorf("NewGameBoard failed to create a gameBoard")
		return
	}

	newPlayer := new(Player)
	newPlayer.Name = "Joe"
	newPlayer.Type = "invalid" // use a non valid constant

	addedPlayer, err := board.AddPlayer(newPlayer)

	if err == nil {
		t.Errorf("Adding a player with an unknown type SHOULD return an error")
		return
	}

	if addedPlayer != nil {
		t.Errorf("Player should NOT have been added")
		return
	}

	fmt.Println("TestAddPlayerFailsWithInvalidType - ended")

}

func TestIsMoveValidWorksWithValidXMove(t *testing.T) {

	fmt.Println("TestIsMoveValidWorksWithValidXMove - started")

	existingLocation := Point{10, 10}
	newLocation := Point{11, 10}

	if !isMoveValid(existingLocation, newLocation) {
		t.Errorf("isMoveValid should allow this move")
		return
	}

	fmt.Println("TestIsMoveValidWorksWithValidXMove - ended")

}

func TestIsMoveValidFailsWithInvalidXMove(t *testing.T) {

	fmt.Println("TestIsMoveValidFailsWithInvalidXMove - started")

	existingLocation := Point{10, 10}
	newLocation := Point{13, 10}

	if isMoveValid(existingLocation, newLocation) {
		t.Errorf("isMoveValid should NOT allow this move")
		return
	}

	fmt.Println("TestIsMoveValidFailsWithInvalidXMove - ended")

}

func TestIsMoveValidWorksWithValidYMove(t *testing.T) {

	fmt.Println("TestIsMoveValidWorksWithValidXMove - started")

	existingLocation := Point{10, 10}
	newLocation := Point{10, 11}

	if !isMoveValid(existingLocation, newLocation) {
		t.Errorf("isMoveValid should allow this move")
		return
	}

	fmt.Println("TestIsMoveValidWorksWithValidYMove - ended")

}

func TestIsMoveValidFailsWithInvalidYMove(t *testing.T) {

	fmt.Println("TestIsMoveValidFailsWithInvalidYMove - started")

	existingLocation := Point{10, 10}
	newLocation := Point{10, 7}

	if isMoveValid(existingLocation, newLocation) {
		t.Errorf("isMoveValid should NOT allow this move")
		return
	}

	fmt.Println("TestIsMoveValidFailsWithInvalidYMove - ended")

}

func TestIsMoveValidFailsWithInvalidXYMove(t *testing.T) {

	fmt.Println("TestIsMoveValidFailsWithInvalidXYMove - started")

	existingLocation := Point{10, 10}
	newLocation := Point{11, 11}

	if isMoveValid(existingLocation, newLocation) {
		t.Errorf("isMoveValid should NOT allow this move")
		return
	}

	fmt.Println("TestIsMoveValidFailsWithInvalidXYMove - ended")

}

func TestMovePlayerWithValidMoveWorks(t *testing.T) {

	setup()
	defer tearDown()

	fmt.Println("TestMovePlayerWithValidMoveWorks started")

	gameId := playingBoardId
	playerId := playingBoardPlayerId

	// fetch board
	board, err := LoadGameBoard(gameId)

	if err != nil {
		t.Errorf("Error fetching game:", err.Error)
		return
	}

	if board == nil {
		t.Errorf("Error: GameBoard not returned")
		return
	}

	player := board.getPlayer(playerId)

	if player == nil {
		t.Errorf("Player not found in game")
		return
	}

	// move player right

	player.Location.X++

	board.ConcurrentMovePlayer(*player)
	// move player
	if err != nil {
		t.Errorf("Error moving player on board:", err.Error)
		return
	}

	// fetch moved player from board
	movedPlayer := board.getPlayer(playerId)

	if movedPlayer == nil {
		t.Errorf("Moved Player not found in game")
		return
	}

	// check player has actually moved
	if (movedPlayer.Location.X != player.Location.X) || (movedPlayer.Location.Y != player.Location.Y) {
		t.Errorf("Player has not moved")
		return
	}

	fmt.Println("TestMovePlayerWithValidMoveWorks ended")

}

/* helper functions */

func setup() {
	fmt.Println("Test setup")
	GamePersister.DeleteAll()

	addTestGames()
}

func tearDown() {
	fmt.Println("Test teardown")
	deleteAllGames()
}

func deleteAllGames() {
	// delete all the games in the games Persister
	GamePersister.DeleteAll()
}

func addTestGames() {
	// create board at playing state
	var playingBoard *GameBoard
	playingBoard, err := NewGameBoard("test game", 1, 4, 1)

	playingBoardId = playingBoard.Id
	playingBoard.State = PlayingGame
	// add player
	newPlayer := &Player{Name: "Player", Type: GoMan}
	addedPlayer, err := playingBoard.AddPlayer(newPlayer)

	if err != nil {
		panic(err)
		return
	}

	playingBoardPlayerId = addedPlayer.Id

	// create board at waiting state
	GamePersister.Create(playingBoard)
	var board2 *GameBoard
	board2, err = NewGameBoard("test game", 1, 4, 1)

	GamePersister.Create(board2)

	var board3 *GameBoard
	board3, err = NewGameBoard("test game", 1, 4, 1)

	GamePersister.Create(board3)

}
