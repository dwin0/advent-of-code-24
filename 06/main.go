package main

import (
	"bufio"
	"fmt"
	"log"
	"slices"
	"strconv"
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

// Part 2 (not working)

var visitedPlaces = make(map[string][]string)

func main() {
	guardMap := readMap("./map.txt")
	guardRowIndex, guardColIndex := findGuard(&guardMap)

	totalPossibleLoops := 0

	for iRow, row := range guardMap {
		for iCol, val := range row {
			if val == "#" || val == "^" {
				continue
			}

			if loopByPlacingObstruction(copyMap(guardMap), iRow, iCol, guardRowIndex, guardColIndex) {
				totalPossibleLoops++
				fmt.Println(iRow, iCol)
			}
		}
	}

	fmt.Println("totalPossibleLoops", totalPossibleLoops)
}

func copyMap(guardMap [][]string) [][]string {
	duplicate := make([][]string, len(guardMap))
	for i := range guardMap {
		duplicate[i] = make([]string, len(guardMap[i]))
		copy(duplicate[i], guardMap[i])
	}
	return duplicate
}

func loopByPlacingObstruction(guardMap [][]string, rowIndex int, colIndex int, guardRowPosition int, guardColPosition int) bool {
	guardMap[rowIndex][colIndex] = "#"
	direction := startingDirection
	addToVisitedPlaces(guardRowPosition, guardColPosition, direction)
	defer clearVisitedPlaces()

	for {
		updatedRowIndex, updatedColIndex, updatedDirection, hasLeftArea := guardMove(&guardMap, guardRowPosition, guardColPosition, direction)

		if hasLeftArea {
			return false
		} else if isLoop(updatedRowIndex, updatedColIndex, updatedDirection) {
			// fmt.Println(updatedRowIndex, updatedColIndex, updatedDirection)
			return true
		} else {
			// fmt.Println(updatedRowIndex, updatedColIndex, updatedDirection)
			addToVisitedPlaces(updatedRowIndex, updatedColIndex, updatedDirection)
		}

		guardRowPosition = updatedRowIndex
		guardColPosition = updatedColIndex
		direction = updatedDirection
	}
}

func addToVisitedPlaces(rowIndex int, colIndex int, direction string) {
	rIndex := strconv.Itoa(rowIndex)
	cIndex := strconv.Itoa(colIndex)

	existingDirections, found := visitedPlaces[rIndex+"_"+cIndex]
	if found {
		visitedPlaces[rIndex+"_"+cIndex] = append(existingDirections, direction)
	} else {
		visitedPlaces[rIndex+"_"+cIndex] = []string{direction}
	}
}

func clearVisitedPlaces() {
	visitedPlaces = make(map[string][]string)
}

func isLoop(rowIndex int, colIndex int, direction string) bool {
	rIndex := strconv.Itoa(rowIndex)
	cIndex := strconv.Itoa(colIndex)

	existingDirections, found := visitedPlaces[rIndex+"_"+cIndex]
	if found {
		return slices.Contains(existingDirections, direction)
	}

	return false
}

func guardMove(guardMap *[][]string, rowIndex int, colIndex int, direction string) (newRowIndex int, newColIndex int, newDirection string, leftArea bool) {
	switch direction {
	case "up":
		{
			if rowIndex-1 < 0 {
				return rowIndex, colIndex, direction, true
			}

			if (*guardMap)[rowIndex-1][colIndex] == "#" {
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
