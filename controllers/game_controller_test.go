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

	fmt.Println("TestCreateGame ended")

}

func TestIsMoveValidWorksWithValidXMove(t *testing.T) {

	fmt.Println("TestIsMoveValidWorksWithValidXMove started")

	existingLocation := &models.Point{10, 10}
	newLocation := &models.Point{11, 10}

	if !isMoveValid(existingLocation, newLocation) {
		log.Fatal("isMoveValid should allow this move")
	}

	fmt.Println("TestIsMoveValidWorksWithValidXMove ended")

}

func TestIsMoveValidFailsWithInvalidXMove(t *testing.T) {

	fmt.Println("TestIsMoveValidFailsWithInvalidXMove started")

	existingLocation := &models.Point{10, 10}
	newLocation := &models.Point{13, 10}

	if isMoveValid(existingLocation, newLocation) {
		log.Fatal("isMoveValid should NOT allow this move")
	}

	fmt.Println("TestIsMoveValidFailsWithInvalidXMove ended")

}

func TestIsMoveValidWorksWithValidYMove(t *testing.T) {

	fmt.Println("TestIsMoveValidWorksWithValidXMove started")

	existingLocation := &models.Point{10, 10}
	newLocation := &models.Point{10, 11}

	if !isMoveValid(existingLocation, newLocation) {
		log.Fatal("isMoveValid should allow this move")
	}

	fmt.Println("TestIsMoveValidWorksWithValidYMove ended")

}

func TestIsMoveValidFailsWithInvalidYMove(t *testing.T) {

	fmt.Println("TestIsMoveValidFailsWithInvalidYMove started")

	existingLocation := &models.Point{10, 10}
	newLocation := &models.Point{10, 7}

	if isMoveValid(existingLocation, newLocation) {
		log.Fatal("isMoveValid should NOT allow this move")
	}

	fmt.Println("TestIsMoveValidFailsWithInvalidYMove ended")

}

func TestIsMoveValidFailsWithInvalidXYMove(t *testing.T) {

	fmt.Println("TestIsMoveValidFailsWithInvalidXYMove started")

	existingLocation := &models.Point{10, 10}
	newLocation := &models.Point{11, 11}

	if isMoveValid(existingLocation, newLocation) {
		log.Fatal("isMoveValid should NOT allow this move")
	}

	fmt.Println("TestIsMoveValidFailsWithInvalidXYMove ended")

}
