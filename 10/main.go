package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Position struct {
	rowIndex    int
	columnIndex int
}

type FloatingIsland struct {
	topographicMap [][]int
	trailheads     []Position
	sumOfAllScores int
}

func main() {
	island := FloatingIsland{
		topographicMap: [][]int{},
		trailheads:     []Position{},
	}

	err := island.readTopographicMap("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	island.countTrailheadsScore()
	island.log()
}

func (fi *FloatingIsland) log() {
	fmt.Println("SumOfAllScores: ", fi.sumOfAllScores)
}

func (fi *FloatingIsland) countTrailheadsScore() {
	for _, trailhead := range fi.trailheads {
		nineHeightPositions := fi.find9heightPositions(trailhead.rowIndex, trailhead.columnIndex, 1)
		// Part 1
		// fi.sumOfAllScores += numberOfUniquePositions(nineHeightPositions)

		// Part 2
		fi.sumOfAllScores += len(nineHeightPositions)
	}
}

func (fi *FloatingIsland) find9heightPositions(rowIndex, columnIndex, stepNumber int) []Position {
	positions := make([]Position, 0)

	for _, rowOffset := range []int{-1, 1} {
		rowIndexToTest := rowIndex + rowOffset

		if rowIndexToTest < 0 || rowIndexToTest >= len(fi.topographicMap) {
			continue
		}

		if fi.topographicMap[rowIndexToTest][columnIndex] == stepNumber {
			if stepNumber == 9 {
				positions = append(positions, Position{rowIndex: rowIndexToTest, columnIndex: columnIndex})
			} else {
				positions = append(positions, fi.find9heightPositions(rowIndexToTest, columnIndex, stepNumber+1)...)
			}
		}
	}

	for _, colOffset := range []int{-1, 1} {
		colIndexToTest := columnIndex + colOffset

		if colIndexToTest < 0 || colIndexToTest >= len(fi.topographicMap[0]) {
			continue
		}

		if fi.topographicMap[rowIndex][colIndexToTest] == stepNumber {
			if stepNumber == 9 {
				positions = append(positions, Position{rowIndex: rowIndex, columnIndex: colIndexToTest})
			} else {
				positions = append(positions, fi.find9heightPositions(rowIndex, colIndexToTest, stepNumber+1)...)
			}
		}
	}

	return positions
}

func (fi *FloatingIsland) readTopographicMap(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	for sc.Scan() {
		line := sc.Text()
		mapRow := make([]int, len(line))

		for i, heightString := range strings.Split(line, "") {
			height, conversionErr := strconv.Atoi(heightString)

			if conversionErr != nil {
				return conversionErr
			}

			if height == 0 {
				fi.trailheads = append(fi.trailheads, Position{rowIndex: len(fi.topographicMap), columnIndex: i})
			}

			mapRow[i] = height
		}

		fi.topographicMap = append(fi.topographicMap, mapRow)
	}

	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}

func numberOfUniquePositions(positions []Position) int {
	unique := make(map[string]bool, 0)

	for _, position := range positions {
		unique[fmt.Sprintf("%v_%v", position.rowIndex, position.columnIndex)] = true
	}

	return len(unique)
}
