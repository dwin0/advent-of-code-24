package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"os"
)

var startingDirection = "up"

// Part 1
// func main() {
// 	guardMap := readMap("./map-small.txt")
// 	rowIndex, colIndex := findGuard(&guardMap)
// 	markSignInMap(&guardMap, rowIndex, colIndex, "X")
// 	direction := startingDirection

// 	for {
// 		updatedRowIndex, updatedColIndex, updatedDirection, hasLeftArea := guardMove(&guardMap, rowIndex, colIndex, direction)
// 		rowIndex = updatedRowIndex
// 		colIndex = updatedColIndex
// 		direction = updatedDirection

// 		if hasLeftArea {
// 			break
// 		} else {
// 			markSignInMap(&guardMap, rowIndex, colIndex, "X")
// 			fmt.Println(guardMap)
// 		}
// 	}

// 	fmt.Println("X in map", countXInMap(&guardMap))
// }

var signForDirection = map[string]string{
	"up":    "|",
	"down":  "|",
	"left":  "-",
	"right": "-",
}

// Part 2
func main() {
	guardMap := readMap("./map-small.txt")
	rowIndex, colIndex := findGuard(&guardMap)
	direction := startingDirection
	markSignInMap(&guardMap, rowIndex, colIndex, signForDirection[direction])

	for {
		updatedRowIndex, updatedColIndex, updatedDirection, hasLeftArea := guardMove(&guardMap, rowIndex, colIndex, direction)

		if hasLeftArea {
			break
		} else {
			if direction != updatedDirection {
				markSignInMap(&guardMap, rowIndex, colIndex, "+")
			}

			markSignInMap(&guardMap, updatedRowIndex, updatedColIndex, signForDirection[updatedDirection])
		}

		rowIndex = updatedRowIndex
		colIndex = updatedColIndex
		direction = updatedDirection
		fmt.Println(guardMap)
	}

}

func guardMove(guardMap *[][]string, rowIndex int, colIndex int, direction string) (newRowIndex int, newColIndex int, newDirection string, leftArea bool) {
	switch direction {
	case "up":
		{
			if rowIndex-1 < 0 {
				return rowIndex, colIndex, direction, true
			}

			if (*guardMap)[rowIndex-1][colIndex] == "#" {
				fmt.Println("found #")
				return rowIndex, colIndex + 1, "right", false
			}
			return rowIndex - 1, colIndex, direction, false
		}
	case "down":
		{
			if rowIndex+1 == len(*guardMap) {
				return rowIndex, colIndex, direction, true
			}

			if (*guardMap)[rowIndex+1][colIndex] == "#" {
				fmt.Println("found #")
				return rowIndex, colIndex - 1, "left", false
			}
			return rowIndex + 1, colIndex, direction, false
		}
	case "left":
		{
			if colIndex-1 < 0 {
				return rowIndex, colIndex, direction, true
			}

			if (*guardMap)[rowIndex][colIndex-1] == "#" {
				fmt.Println("found #")
				return rowIndex - 1, colIndex, "up", false
			}
			return rowIndex, colIndex - 1, direction, false
		}
	case "right":
		{
			if colIndex+1 == len((*guardMap)[0]) {
				return rowIndex, colIndex, direction, true
			}

			if (*guardMap)[rowIndex][colIndex+1] == "#" {
				fmt.Println("found #")
				return rowIndex + 1, colIndex, "down", false
			}
			return rowIndex, colIndex + 1, direction, false
		}
	default:
		{
			log.Fatal("Unknown direction", direction)
			return
		}
	}
}

func countXInMap(guardMap *[][]string) int {
	totalX := 0

	for _, row := range *guardMap {
		for _, char := range row {
			if char == "X" {
				totalX++
			}
		}
	}

	return totalX
}

func markSignInMap(guardMap *[][]string, rowIndex int, colIndex int, sign string) {
	// Part 1
	// (*guardMap)[rowIndex][colIndex] = sign

	existingSign := (*guardMap)[rowIndex][colIndex]
	if existingSign != "." && existingSign != sign {
		// TODO: we have been here before, check if there are # around
		// additionally check if we are close to a wall and handle this case

		(*guardMap)[rowIndex][colIndex] = "+"
	} else {
		(*guardMap)[rowIndex][colIndex] = sign
	}
}

func findGuard(guardMap *[][]string) (rowIndex int, colIndex int) {
	for iRow, row := range *guardMap {
		for iCol, char := range row {
			if char == "^" {
				return iRow, iCol
			}
		}
	}

	log.Fatal("No guard found")
	return
}

func readMap(filePath string) [][]string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	guardMap := make([][]string, 0)

	for sc.Scan() {
		line := sc.Text()
		symbols := strings.Split(line, "")
		guardMap = append(guardMap, symbols)

	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	return guardMap
}
