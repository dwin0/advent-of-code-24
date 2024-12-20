package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type Position struct {
	rowIndex int
	colIndex int
}

type Direction int

// clockwise direction
const (
	Up    Direction = 0
	Right Direction = 90
	Down  Direction = 180
	Left  Direction = 270
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

const startDirection Direction = Right
const startSymbol = "S"
const endSymbol = "E"
const wallSymbol = "#"
const moveCosts = 1
const rotationCosts = 1000

type ReindeerOlympics struct {
	maze          [][]string
	startPosition Position
	endPosition   Position
}

// https://www.redblobgames.com/pathfinding/a-star/introduction.html
func main() {
	reindeerOlympics := ReindeerOlympics{}

	err := reindeerOlympics.readMaze("./input-small.txt")
	if err != nil {
		log.Fatal(err)
	}

	err = reindeerOlympics.findStartAndEnd()
	if err != nil {
		log.Fatal(err)
	}

	costsPerTile := reindeerOlympics.findLowestScore()
	fmt.Println("lowestScore", costsPerTile[reindeerOlympics.endPosition])
	// fmt.Println("numberOfTilesInBestPathThroughMaze", reindeerOlympics.findNumberOfTilesInBestPathThroughMaze(&costsPerTile))
}

// Part 2: not working
// func (ro *ReindeerOlympics) findNumberOfTilesInBestPathThroughMaze(costsPerTile *map[Position]int) int {
// 	tilesInBestPathThroughMaze := make(map[Position]bool, 0)
// 	tilesInBestPathThroughMaze[ro.endPosition] = true
// 	tilesInBestPathThroughMaze[ro.startPosition] = true
// 	nextTilesToCheck := []Position{ro.endPosition}

// 	for len(nextTilesToCheck) > 0 {
// 		tileToCheck := nextTilesToCheck[0]
// 		otherTiles := nextTilesToCheck[1:]

// 		if tileToCheck == ro.startPosition {
// 			break
// 		}

// 		neightbours := slices.DeleteFunc(ro.getNeighbours(tileToCheck), func(p Position) bool {
// 			_, found := tilesInBestPathThroughMaze[p]
// 			isWall := (*costsPerTile)[p] == 0
// 			return found || isWall
// 		})

// 		fmt.Println(neightbours)
// 		fmt.Scanln()

// 		neightboursWithLowestScore := []Position{}
// 		for _, neighbour := range neightbours {
// 			if len(neightboursWithLowestScore) == 0 {
// 				neightboursWithLowestScore = append(neightboursWithLowestScore, neighbour)
// 				continue
// 			}
// 			if (*costsPerTile)[neightboursWithLowestScore[0]] < (*costsPerTile)[neighbour] {
// 				continue
// 			}
// 			if (*costsPerTile)[neightboursWithLowestScore[0]] == (*costsPerTile)[neighbour] {
// 				neightboursWithLowestScore = append(neightboursWithLowestScore, neighbour)
// 				continue
// 			}
// 			if (*costsPerTile)[neightboursWithLowestScore[0]] > (*costsPerTile)[neighbour] {
// 				neightboursWithLowestScore = []Position{neighbour}
// 			}
// 		}

// 		for _, neighbour := range neightboursWithLowestScore {
// 			tilesInBestPathThroughMaze[neighbour] = true
// 		}

// 		nextTilesToCheck = append(otherTiles, neightboursWithLowestScore...)
// 	}

// 	return len(tilesInBestPathThroughMaze)
// }

func (ro *ReindeerOlympics) findLowestScore() map[Position]int {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{
		position:  ro.startPosition,
		direction: startDirection,
		priority:  0,
	})
	costsSoFar := make(map[Position]int)
	costsSoFar[ro.startPosition] = 0

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Item)

		if current.position == ro.endPosition {
			break
		}

		for _, next := range ro.getNeighbours(current.position) {
			newCost := costsSoFar[current.position] + getCost(current.position, next, current.direction)

			costsSoFarForNext, found := costsSoFar[next]
			if !found || newCost < costsSoFarForNext {
				costsSoFar[next] = newCost
				heap.Push(&pq, &Item{
					position:  next,
					priority:  newCost,
					direction: getDirection(current.position, next),
				})
			}
		}
	}

	return costsSoFar
}

func getDirection(currentPosition, nextPosition Position) Direction {
	rowOffset := currentPosition.rowIndex - nextPosition.rowIndex
	colOffset := currentPosition.colIndex - nextPosition.colIndex

	if rowOffset == 1 {
		return Up
	}
	if rowOffset == -1 {
		return Down
	}
	if colOffset == 1 {
		return Left
	}
	return Right
}

func getCost(currentPosition, nextPosition Position, currentDirection Direction) int {
	rowOffset := currentPosition.rowIndex - nextPosition.rowIndex
	colOffset := currentPosition.colIndex - nextPosition.colIndex

	if rowOffset != 0 && (currentDirection == Up || currentDirection == Down) {
		return moveCosts
	} else if colOffset != 0 && (currentDirection == Left || currentDirection == Right) {
		return moveCosts
	} else {
		return rotationCosts + moveCosts
	}
}

func (ro *ReindeerOlympics) getNeighbours(position Position) []Position {
	possibleNeighbours := []Position{
		{rowIndex: position.rowIndex + 1, colIndex: position.colIndex},
		{rowIndex: position.rowIndex - 1, colIndex: position.colIndex},
		{rowIndex: position.rowIndex, colIndex: position.colIndex + 1},
		{rowIndex: position.rowIndex, colIndex: position.colIndex - 1},
	}

	return slices.DeleteFunc(possibleNeighbours, func(p Position) bool {
		return ro.maze[p.rowIndex][p.colIndex] == wallSymbol
	})
}

func (ro *ReindeerOlympics) readMaze(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	ro.maze = make([][]string, 0)

	for sc.Scan() {
		line := sc.Text()
		ro.maze = append(ro.maze, strings.Split(line, ""))
	}

	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}

func (ro *ReindeerOlympics) findStartAndEnd() error {
	for rowIndex, row := range ro.maze {
		for colIndex, cell := range row {
			if cell == startSymbol {
				ro.startPosition = Position{
					rowIndex: rowIndex,
					colIndex: colIndex,
				}
			} else if cell == endSymbol {
				ro.endPosition = Position{
					rowIndex: rowIndex,
					colIndex: colIndex,
				}
			}

			if ro.startPosition != (Position{}) && ro.endPosition != (Position{}) {
				return nil
			}
		}
	}

	return errors.New("not found im maze")
}
