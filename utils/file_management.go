package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var filePath = GetAbsolutePathOfCurrentPackage("../gamedata/")

func DeleteOldGameBoardFiles() {
	fmt.Println("Deleting old game boards")
	dirs, err := ioutil.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, dir := range dirs {
		name := dir.Name()
		if !dir.IsDir() && strings.HasSuffix(name, ".json") {
			fmt.Println("Deleting file:", name)
			err := os.Remove(filePath + "/" + dir.Name())
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
