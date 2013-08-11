package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/telecoda/go-man-app/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var playingBoardId string
var playingBoardPlayerId string

var waitingBoardId string
var waitingBoardPlayerId string

func TestCreateGame(t *testing.T) {

	setup()
	defer tearDown()

	fmt.Println("TestCreateGame started")

	ts := httptest.NewServer(http.HandlerFunc(GameCreate))
	defer ts.Close()

	res, err := http.Post(ts.URL, "application/json", nil)
	if err != nil {
		t.Errorf("Error POSTing request to API:", err)
	}
	jsonBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Errorf("Failed to read JSON response", err.Error)
	}

	if jsonBody == nil {
		t.Errorf("No json returned")
	}

	var board models.GameBoard

	err = json.Unmarshal(jsonBody, &board)

	if err != nil {
		t.Errorf("Failed to unmarshal JSON response", err.Error)
	}

	// check values of board returned

	if &board == nil {
		t.Errorf("No game board")
	}

	if len(board.Id) == 0 {
		t.Errorf("No gameboard.Id")
	}

	err = board.DestroyGameBoard()
	if err != nil {
		t.Errorf("DestroyGameBoard failed:", err)
	}

	fmt.Println("TestCreateGame ended")

}

func TestGameList(t *testing.T) {

	setup()
	defer tearDown()

	fmt.Println("TestGameList started")

	ts := httptest.NewServer(http.HandlerFunc(GameList))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Errorf("Error GETting request from API:", err)
	}
	jsonBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Errorf("Failed to read JSON response", err.Error)
	}

	if jsonBody == nil {
		t.Errorf("No json returned")
	}

	var boards []models.GameBoardSummary

	fmt.Println(boards)

	err = json.Unmarshal(jsonBody, &boards)

	if err != nil {
		t.Errorf("Failed to unmarshal JSON response", err.Error)
	}

	// check values of boards returned

	if len(boards) != 2 {
		t.Errorf("There should be 2 boards but we received %d", len(boards))
	}

	if &boards == nil {
		t.Errorf("No game boards")
	}

	fmt.Println("TestGameList ended")

}

func TestGameListFilteredByState(t *testing.T) {

	setup()
	defer tearDown()

	fmt.Println("TestGameListFilteredByState started")

	ts := httptest.NewServer(http.HandlerFunc(GameList))
	defer ts.Close()

	requestURL := ts.URL + "?state=waiting"
	res, err := http.Get(requestURL)
	if err != nil {
		t.Errorf("Error GETting request from API:", err)
	}
	jsonBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Errorf("Failed to read JSON response", err.Error)
	}

	if jsonBody == nil {
		t.Errorf("No json returned")
	}

	var boards []models.GameBoardSummary

	fmt.Println(boards)

	err = json.Unmarshal(jsonBody, &boards)

	if err != nil {
		t.Errorf("Failed to unmarshal JSON response", err.Error)
	}

	// check values of boards returned

	if len(boards) != 1 {
		t.Errorf("There should be 1 boards but we received %d", len(boards))
	}

	if &boards == nil {
		t.Errorf("No game boards")
	}

	fmt.Println("TestGameListFilteredByState ended")

}

func setup() {
	fmt.Println("Test setup")
	models.GamePersister.DeleteAll()

	addTestGames()
}

func tearDown() {
	fmt.Println("Test teardown")
	deleteAllGames()
}

func deleteAllGames() {
	// delete all the games in the games persister
	models.GamePersister.DeleteAll()
}

func addTestGames() {
	// create board at playing state
	var playingBoard = models.NewGameBoard()
	playingBoardId = playingBoard.Id
	playingBoard.State = models.PlayingGame
	// add player
	newPlayer := &models.Player{Name: "Player", Type: models.GoMan}
	addedPlayer, err := playingBoard.AddPlayer(newPlayer)

	if err != nil {
		panic(err)
		return
	}

	playingBoardPlayerId = addedPlayer.Id

	models.GamePersister.Create(playingBoard)

	// create board at waiting state
	var waitingBoard = models.NewGameBoard()
	waitingBoardId = waitingBoard.Id
	waitingBoard.State = models.WaitingForPlayers
	models.GamePersister.Create(waitingBoard)

}
