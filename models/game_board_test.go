package models

import (
	"fmt"
	"log"
	"testing"
)

func TestIsMoveValidWorksWithValidXMove(t *testing.T) {

	fmt.Println("TestIsMoveValidWorksWithValidXMove started")

	existingLocation := &Point{10, 10}
	newLocation := &Point{11, 10}

	if !IsMoveValid(existingLocation, newLocation) {
		log.Fatal("isMoveValid should allow this move")
	}

	fmt.Println("TestIsMoveValidWorksWithValidXMove ended")

}

func TestIsMoveValidFailsWithInvalidXMove(t *testing.T) {

	fmt.Println("TestIsMoveValidFailsWithInvalidXMove started")

	existingLocation := &Point{10, 10}
	newLocation := &Point{13, 10}

	if IsMoveValid(existingLocation, newLocation) {
		log.Fatal("isMoveValid should NOT allow this move")
	}

	fmt.Println("TestIsMoveValidFailsWithInvalidXMove ended")

}

func TestIsMoveValidWorksWithValidYMove(t *testing.T) {

	fmt.Println("TestIsMoveValidWorksWithValidXMove started")

	existingLocation := &Point{10, 10}
	newLocation := &Point{10, 11}

	if !IsMoveValid(existingLocation, newLocation) {
		log.Fatal("isMoveValid should allow this move")
	}

	fmt.Println("TestIsMoveValidWorksWithValidYMove ended")

}

func TestIsMoveValidFailsWithInvalidYMove(t *testing.T) {

	fmt.Println("TestIsMoveValidFailsWithInvalidYMove started")

	existingLocation := &Point{10, 10}
	newLocation := &Point{10, 7}

	if IsMoveValid(existingLocation, newLocation) {
		log.Fatal("isMoveValid should NOT allow this move")
	}

	fmt.Println("TestIsMoveValidFailsWithInvalidYMove ended")

}

func TestIsMoveValidFailsWithInvalidXYMove(t *testing.T) {

	fmt.Println("TestIsMoveValidFailsWithInvalidXYMove started")

	existingLocation := &Point{10, 10}
	newLocation := &Point{11, 11}

	if IsMoveValid(existingLocation, newLocation) {
		log.Fatal("isMoveValid should NOT allow this move")
	}

	fmt.Println("TestIsMoveValidFailsWithInvalidXYMove ended")

}
