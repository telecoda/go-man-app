package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/telecoda/go-man/models"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateGame(t *testing.T) {

	fmt.Println("TestCreateGame started")

	ts := httptest.NewServer(http.HandlerFunc(GameCreate))
	defer ts.Close()

	res, err := http.Post(ts.URL, "application/json", nil)
	if err != nil {
		log.Fatal(err)
	}
	jsonBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("TestCreateGame ended")

	if jsonBody == nil {
		log.Fatal("No json returned")
	}

	var board models.GameBoard

	err = json.Unmarshal(jsonBody, &board)

	if err != nil {
		log.Fatal(err)
	}

	// check values of board returned

	if &board == nil {
		log.Fatal("No game board")
	}

	if len(board.Id) == 0 {
		log.Fatal("No gameboard.Id")
	}

}
