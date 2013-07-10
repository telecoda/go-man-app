package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var filePath = GetAbsolutePathOfCurrentPackage(".." + string(os.PathSeparator) + "gamedata" + string(os.PathSeparator))

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
			err := os.Remove(filePath + string(os.PathSeparator) + dir.Name())
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func CopyGameDataFixtures(fromDir string, toDir string) {
	fromPath := GetAbsolutePathOfCurrentPackage(".." + string(os.PathSeparator) + fromDir)
	toPath := GetAbsolutePathOfCurrentPackage(".." + string(os.PathSeparator) + toDir)

	dirs, err := ioutil.ReadDir(fromPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, dir := range dirs {
		name := dir.Name()
		if !dir.IsDir() && strings.HasSuffix(name, ".json") {
			fmt.Println("Copying file:", name)
			oldFile, err := os.Open(fromPath + string(os.PathSeparator) + dir.Name())
			if err != nil {
				fmt.Println(err)
				break
			}
			newFile, err := os.Create(toPath + string(os.PathSeparator) + dir.Name())
			if err != nil {
				fmt.Println(err)
				break
			}

			readBuffer := make([]byte, 5000)

			bytesRead, err := oldFile.Read(readBuffer)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("Bytes read", bytesRead)

			// truncate buffer to actual size read
			readBuffer = readBuffer[:bytesRead]

			bytesWritten, err := newFile.Write(readBuffer)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("Bytes written", bytesWritten)
		}
	}

}
