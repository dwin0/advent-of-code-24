package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var errLoop = errors.New("guard stuck in a loop")

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

var directionName = map[Direction]string{
	Up:    "up",
	Right: "right",
	Down:  "down",
	Left:  "left",
}

func (d Direction) String() string {
	return directionName[d]
}

type Position struct {
	rowIndex int
	colIndex int
}

type Guard struct {
	position      Position
	direction     Direction
	visitedPlaces map[string]bool
}

func (g *Guard) update(rowIndex, colIndex int, direction Direction) error {
	g.position = Position{rowIndex: rowIndex, colIndex: colIndex}
	g.direction = direction

	placeId := fmt.Sprintf("%d_%d_%s", rowIndex, colIndex, direction)

	visitedBefore := g.visitedPlaces[placeId]
	if visitedBefore {
		return errLoop
	}

	g.visitedPlaces[placeId] = true
	return nil
}

type ManufacturingLab struct {
	labMap        [][]string
	guard         Guard
	possibleLoops int
}

func main() {
	guard := Guard{position: Position{colIndex: -1, rowIndex: -1}, visitedPlaces: make(map[string]bool, 0)}
	manufacturingLab := ManufacturingLab{
		labMap: [][]string{},
		guard:  guard,
	}

	err := manufacturingLab.readLabMap("./map-small.txt")
	if err != nil {
		log.Fatal(err)
	}

	err = manufacturingLab.findGuard()
	if err != nil {
		log.Fatal(err)
	}

	// Part 1
	// for {
	// 	done, err := manufacturingLab.guardMove()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	if done {
	// 		break
	// 	}
	// }

	// Part 2
	err = manufacturingLab.countPossibleLoops()
	if err != nil {
		log.Fatal(err)
	}
	manufacturingLab.log()
}

func (ml *ManufacturingLab) log() {
	// fmt.Println(ml.labMap)
	// fmt.Println("distinct positions", ml.countXInMap())
	fmt.Println("possibleLoops", ml.possibleLoops)
	fmt.Println(ml.guard.visitedPlaces)
}

func (ml *ManufacturingLab) countPossibleLoops() error {
	possibleLoops := 0

	for rowIndex, row := range ml.labMap {
		for columnIndex, char := range row {
			if char == "." {
				hasLoop, err := ml.hasLoopWithObstacleAt(rowIndex, columnIndex)
				fmt.Println(hasLoop, err)

				if err != nil {
					return err
				}

				if hasLoop {
					possibleLoops++
				}

			}
		}
	}

	ml.possibleLoops = possibleLoops
	return nil
}

func (ml *ManufacturingLab) hasLoopWithObstacleAt(rowIndex, colIndex int) (bool, error) {
	originalMap := ml.labMap
	mapCopy := copyLabMap(ml.labMap)
	mapCopy[rowIndex][colIndex] = "#"
	ml.labMap = mapCopy

	originalPosition := Position{rowIndex: ml.guard.position.rowIndex, colIndex: ml.guard.position.colIndex}
	originalDirection := ml.guard.direction

	defer func() {
		ml.guard.visitedPlaces = make(map[string]bool, 0)
		ml.labMap = originalMap
		ml.guard.update(originalPosition.rowIndex, originalPosition.colIndex, originalDirection)
	}()

	for {
		done, err := ml.guardMove()

		if err == errLoop {
			return true, nil
		}

		if err != nil {
			return false, err
		}

		if done {
			return false, nil
		}
	}
}

func copyLabMap(originalMap [][]string) (copy [][]string) {
	copy = make([][]string, len(originalMap))

	for rowIndex, row := range originalMap {
		copy[rowIndex] = make([]string, len(row))

		for columnIndex, char := range row {
			copy[rowIndex][columnIndex] = char
		}
	}

	return copy
}

func (ml *ManufacturingLab) guardMove() (done bool, err error) {
	position := ml.guard.position
	direction := ml.guard.direction

	switch direction {
	case Up:
		{
			if position.rowIndex-1 < 0 {
				return true, nil
			}

			if ml.labMap[position.rowIndex-1][position.colIndex] == "#" {
				err := ml.guard.update(position.rowIndex, position.colIndex, Right)
				return false, err
			}

			ml.markXInMap(position.rowIndex-1, position.colIndex)
			err := ml.guard.update(position.rowIndex-1, position.colIndex, Up)
			return false, err
		}
	case Down:
		{
			if position.rowIndex+1 == len(ml.labMap) {
				return true, nil
			}

			if ml.labMap[position.rowIndex+1][position.colIndex] == "#" {
				err := ml.guard.update(position.rowIndex, position.colIndex, Left)
				return false, err
			}

			ml.markXInMap(position.rowIndex+1, position.colIndex)
			err := ml.guard.update(position.rowIndex+1, position.colIndex, Down)
			return false, err
		}
	case Left:
		{
			if position.colIndex-1 < 0 {
				return true, nil
			}

			if ml.labMap[position.rowIndex][position.colIndex-1] == "#" {
				err := ml.guard.update(position.rowIndex, position.colIndex, Up)
				return false, err
			}

			ml.markXInMap(position.rowIndex, position.colIndex-1)
			err := ml.guard.update(position.rowIndex, position.colIndex-1, Left)
			return false, err
		}
	case Right:
		{
			if position.colIndex+1 == len(ml.labMap[0]) {
				return true, nil
			}

			if ml.labMap[position.rowIndex][position.colIndex+1] == "#" {
				err := ml.guard.update(position.rowIndex, position.colIndex, Down)
				return false, err
			}

			ml.markXInMap(position.rowIndex, position.colIndex+1)
			err := ml.guard.update(position.rowIndex, position.colIndex+1, Right)
			return false, err
		}
	default:
		{
			return false, errors.New("invalid direction")
		}
	}
}

func (ml *ManufacturingLab) countXInMap() int {
	totalX := 0

	for _, row := range ml.labMap {
		for _, char := range row {
			if char == "X" {
				totalX++
			}
		}
	}

	return totalX
}

func (ml *ManufacturingLab) markXInMap(rowIndex int, colIndex int) {
	ml.labMap[rowIndex][colIndex] = "X"
}

func (ml *ManufacturingLab) findGuard() error {
	for iRow, row := range ml.labMap {
		for iCol, char := range row {
			if char == "^" {
				ml.guard.update(iRow, iCol, Up)
				ml.markXInMap(iRow, iCol)
				return nil
			}
		}
	}

	return errors.New("no guard found")
}

func (ml *ManufacturingLab) readLabMap(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	for sc.Scan() {
		line := sc.Text()
		ml.labMap = append(ml.labMap, strings.Split(line, ""))
	}

	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}
