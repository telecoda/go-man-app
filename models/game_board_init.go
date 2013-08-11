package models

import (
	"bufio"
	"github.com/telecoda/go-man-app/utils"
	"os"
)

/* initGameBoard loads board definition from local file */

func initGameBoard() ([][]rune, error) {

	defaultBoard := make([][]rune, BOARD_HEIGHT)

	// this path will be of the controllers package
	filePath := utils.GetAbsolutePathOfCurrentPackage("../data/maze.txt")

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)

	var r int = 0
	for {

		b, err := reader.ReadBytes('\n')
		if err == nil {
			// parse line

			b = b[:len(b)-1] // remove last new line char from bytes
			row := string(b)
			defaultBoard[r] = make([]rune, BOARD_WIDTH)
			for c, cell := range row {
				defaultBoard[r][c] = rune(cell)
				c++
			}
			r++
		} else {
			break
		}

	}

	return defaultBoard, nil

}
