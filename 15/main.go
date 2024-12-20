package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const robotSymbol = "@"
const boxSmallSymbol = "O"
const boxLargeSymbol1 = "["
const boxLargeSymbol2 = "]"
const wallSymbol = "#"
const emptySpaceSymbol = "."

type Position struct {
	rowIndex int
	colIndex int
}

type Warehouse struct {
	boxesAndWallsMap [][]string
	robotPosition    Position
	robotMovements   []string
}

func main() {
	warehouse := Warehouse{}

	err := warehouse.readMap("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	warehouse.transformMap() // Part 2
	warehouse.findRobot()

	for _, movement := range warehouse.robotMovements {
		// fmt.Scanln() // cause the program to wait for user input (manual debugging)
		warehouse.move(movement)
		// warehouse.log()
	}

	warehouse.log()
	// fmt.Println(warehouse.calculateGpsSum(boxSmallSymbol)) // Part 1
	fmt.Println(warehouse.calculateGpsSum(boxLargeSymbol1)) // Part 2
}

func (w *Warehouse) log() {
	mapWithRobot := copyMap(w.boxesAndWallsMap)
	mapWithRobot[w.robotPosition.rowIndex][w.robotPosition.colIndex] = robotSymbol

	fmt.Println(mapWithRobot)
}

func (w *Warehouse) calculateGpsSum(boxSymbolToCheck string) int {
	sum := 0

	for rowIndex, row := range w.boxesAndWallsMap {
		for colIndex, cell := range row {
			if cell == boxSymbolToCheck {
				sum += 100*rowIndex + colIndex
			}
		}
	}

	return sum
}

func (w *Warehouse) move(step string) {
	switch step {
	case "^":
		{
			w.moveRobot(-1, 0)
		}
	case "v":
		{
			w.moveRobot(1, 0)
		}
	case "<":
		{
			w.moveRobot(0, -1)
		}
	case ">":
		{
			w.moveRobot(0, 1)
		}
	}
}

// Part 2
func (w *Warehouse) moveRobot(rowOffset, colOffset int) {
	newRow := w.robotPosition.rowIndex + rowOffset
	newCol := w.robotPosition.colIndex + colOffset

	if w.boxesAndWallsMap[newRow][newCol] == wallSymbol {
		return
	}

	if w.boxesAndWallsMap[newRow][newCol] == boxLargeSymbol1 || w.boxesAndWallsMap[newRow][newCol] == boxLargeSymbol2 {
		if colOffset == 0 {
			/** First check if all boxes CAN be moved, only then move them
			[
			[# # # # # # # # # # # # # #]
			[# # . . . [ ] . # # . . # #]
			[# # . . . . . [ ] . . . # #]
			[# # . . . . [ ] . . . . # #]
			[# # . . . . . @ . . . . # #]
			[# # . . . . . . . . . . # #]
			[# # # # # # # # # # # # # #]
			]
			*/

			allBoxesCanBeMoved := w.moveBoxVertical(w.boxesAndWallsMap[newRow][newCol], newRow, newCol, rowOffset == -1, false)
			if allBoxesCanBeMoved {
				w.moveBoxVertical(w.boxesAndWallsMap[newRow][newCol], newRow, newCol, rowOffset == -1, true)
			} else {
				return
			}
		} else {
			if !w.moveBoxHorizontal(newRow, newCol, colOffset == -1) {
				return
			}
		}
	}

	w.robotPosition = Position{
		rowIndex: newRow,
		colIndex: newCol,
	}
}

// Part 2
func (w *Warehouse) moveBoxVertical(boxSymbol string, rowIndex, colIndex int, up, performMove bool) bool {
	boxPosition1 := Position{}
	boxPosition2 := Position{}
	position1ToCheck := Position{}
	position2ToCheck := Position{}

	moveTheBox := func() {
		if performMove {
			w.moveBox(boxPosition1, boxPosition2, position1ToCheck, position2ToCheck)
		}
	}

	if boxSymbol == boxLargeSymbol1 {
		boxPosition1 = Position{rowIndex: rowIndex, colIndex: colIndex}
		boxPosition2 = Position{rowIndex: rowIndex, colIndex: colIndex + 1}
	} else {
		boxPosition1 = Position{rowIndex: rowIndex, colIndex: colIndex - 1}
		boxPosition2 = Position{rowIndex: rowIndex, colIndex: colIndex}
	}

	if up {
		position1ToCheck = Position{rowIndex: boxPosition1.rowIndex - 1, colIndex: boxPosition1.colIndex}
		position2ToCheck = Position{rowIndex: boxPosition2.rowIndex - 1, colIndex: boxPosition2.colIndex}
	} else {
		position1ToCheck = Position{rowIndex: boxPosition1.rowIndex + 1, colIndex: boxPosition1.colIndex}
		position2ToCheck = Position{rowIndex: boxPosition2.rowIndex + 1, colIndex: boxPosition2.colIndex}
	}

	if w.boxesAndWallsMap[position1ToCheck.rowIndex][position1ToCheck.colIndex] == emptySpaceSymbol && w.boxesAndWallsMap[position2ToCheck.rowIndex][position2ToCheck.colIndex] == emptySpaceSymbol {
		/**
		..
		[] <- current box
		*/

		moveTheBox()
		return true
	} else if w.boxesAndWallsMap[position1ToCheck.rowIndex][position1ToCheck.colIndex] == wallSymbol || w.boxesAndWallsMap[position2ToCheck.rowIndex][position2ToCheck.colIndex] == wallSymbol {
		/**
		.# <- wall
		[] <- current box
		*/

		return false
	} else if w.boxesAndWallsMap[position1ToCheck.rowIndex][position1ToCheck.colIndex] == boxLargeSymbol1 {
		// box directly in front of current box
		/**
		??
		[]
		[] <- current box
		*/

		if w.moveBoxVertical(boxLargeSymbol1, position1ToCheck.rowIndex, position1ToCheck.colIndex, up, performMove) {
			moveTheBox()
			return true
		} else {
			return false
		}
	} else if w.boxesAndWallsMap[position1ToCheck.rowIndex][position1ToCheck.colIndex] == boxLargeSymbol2 && w.boxesAndWallsMap[position2ToCheck.rowIndex][position2ToCheck.colIndex] == boxLargeSymbol1 {
		// 2 boxes diagonal in front of current box
		/**
		 ??
		[][]
		 [] <- current box
		*/

		canMoveLeftBoxInFront := w.moveBoxVertical(boxLargeSymbol2, position1ToCheck.rowIndex, position1ToCheck.colIndex, up, performMove)
		canMoveRightBoxInFront := w.moveBoxVertical(boxLargeSymbol1, position2ToCheck.rowIndex, position2ToCheck.colIndex, up, performMove)

		if canMoveLeftBoxInFront && canMoveRightBoxInFront {
			moveTheBox()
			return true
		} else {
			return false
		}
	} else if w.boxesAndWallsMap[position1ToCheck.rowIndex][position1ToCheck.colIndex] == boxLargeSymbol2 {
		// box diagonal in front of current box
		/**
		 ??
		[].
		 [] <- current box
		*/

		canMoveLeftBoxInFront := w.moveBoxVertical(boxLargeSymbol2, position1ToCheck.rowIndex, position1ToCheck.colIndex, up, performMove)

		if canMoveLeftBoxInFront {
			moveTheBox()
			return true
		} else {
			return false
		}
	} else if w.boxesAndWallsMap[position2ToCheck.rowIndex][position2ToCheck.colIndex] == boxLargeSymbol1 {
		// box diagonal in front of current box
		/**
		??
		.[]
		[] <- current box
		*/

		canMoveRightBoxInFront := w.moveBoxVertical(boxLargeSymbol1, position2ToCheck.rowIndex, position2ToCheck.colIndex, up, performMove)

		if canMoveRightBoxInFront {
			moveTheBox()
			return true
		} else {
			return false
		}
	} else {
		log.Fatal("Unknown case", w.boxesAndWallsMap[position1ToCheck.rowIndex][position1ToCheck.colIndex], w.boxesAndWallsMap[position2ToCheck.rowIndex][position2ToCheck.colIndex])
		return false
	}
}
func (w *Warehouse) moveBoxHorizontal(rowIndex, colIndex int, left bool) bool {
	boxPosition1 := Position{}
	boxPosition2 := Position{}
	positionToCheck := Position{}
	safePosition := Position{} // no need to check second position, as this is where the box current is located
	newBoxPosition1 := Position{}
	newBoxPosition2 := Position{}

	if left {
		boxPosition1 = Position{rowIndex: rowIndex, colIndex: colIndex - 1}
		boxPosition2 = Position{rowIndex: rowIndex, colIndex: colIndex}

		positionToCheck = Position{rowIndex: boxPosition1.rowIndex, colIndex: boxPosition1.colIndex - 1}
		safePosition = Position{rowIndex: boxPosition1.rowIndex, colIndex: boxPosition1.colIndex}

		newBoxPosition1 = positionToCheck
		newBoxPosition2 = safePosition
	} else {
		boxPosition1 = Position{rowIndex: rowIndex, colIndex: colIndex}
		boxPosition2 = Position{rowIndex: rowIndex, colIndex: colIndex + 1}

		positionToCheck = Position{rowIndex: boxPosition2.rowIndex, colIndex: boxPosition2.colIndex + 1}
		safePosition = Position{rowIndex: boxPosition2.rowIndex, colIndex: boxPosition2.colIndex}

		newBoxPosition1 = safePosition
		newBoxPosition2 = positionToCheck
	}

	if w.boxesAndWallsMap[positionToCheck.rowIndex][positionToCheck.colIndex] == emptySpaceSymbol {
		//		.[] <- current box			[]. <- current box

		w.moveBox(boxPosition1, boxPosition2, newBoxPosition1, newBoxPosition2)
		return true
	} else if w.boxesAndWallsMap[positionToCheck.rowIndex][positionToCheck.colIndex] == wallSymbol {
		//		#[]							[]#

		return false
	} else if left && w.boxesAndWallsMap[positionToCheck.rowIndex][positionToCheck.colIndex] == boxLargeSymbol2 {
		//		[][] <- current box

		if w.moveBoxHorizontal(positionToCheck.rowIndex, positionToCheck.colIndex, left) {
			w.moveBox(boxPosition1, boxPosition2, newBoxPosition1, newBoxPosition2)
			return true
		} else {
			return false
		}
	} else if !left && w.boxesAndWallsMap[positionToCheck.rowIndex][positionToCheck.colIndex] == boxLargeSymbol1 {
		// 		current box -> [][]

		if w.moveBoxHorizontal(positionToCheck.rowIndex, positionToCheck.colIndex, left) {
			w.moveBox(boxPosition1, boxPosition2, newBoxPosition1, newBoxPosition2)
			return true
		} else {
			return false
		}
	} else {
		log.Fatal("Unknown case", w.boxesAndWallsMap[positionToCheck.rowIndex][positionToCheck.colIndex], left)
		return false
	}
}

// Part 2
func (w *Warehouse) moveBox(from1, from2, to1, to2 Position) {
	w.boxesAndWallsMap[from1.rowIndex][from1.colIndex] = emptySpaceSymbol
	w.boxesAndWallsMap[from2.rowIndex][from2.colIndex] = emptySpaceSymbol
	w.boxesAndWallsMap[to1.rowIndex][to1.colIndex] = boxLargeSymbol1
	w.boxesAndWallsMap[to2.rowIndex][to2.colIndex] = boxLargeSymbol2
}

// Part 1
// func (w *Warehouse) moveRobot(rowOffset, colOffset int) {
// 	newRow := w.robotPosition.rowIndex + rowOffset
// 	newCol := w.robotPosition.colIndex + colOffset

// 	if w.boxesAndWallsMap[newRow][newCol] == wallSymbol {
// 		return
// 	}

// 	if w.boxesAndWallsMap[newRow][newCol] == boxSmallSymbol {
// 		if rowOffset == -1 {
// 			for i := newRow - 1; i >= 0; i-- {
// 				if w.boxesAndWallsMap[i][newCol] == wallSymbol {
// 					return
// 				} else if w.boxesAndWallsMap[i][newCol] == emptySpaceSymbol {
// 					w.boxesAndWallsMap[newRow][newCol] = emptySpaceSymbol
// 					w.boxesAndWallsMap[i][newCol] = boxSmallSymbol
// 					break
// 				}
// 			}
// 		} else if rowOffset == 1 {
// 			for i := newRow + 1; i < len(w.boxesAndWallsMap); i++ {
// 				if w.boxesAndWallsMap[i][newCol] == wallSymbol {
// 					return
// 				} else if w.boxesAndWallsMap[i][newCol] == emptySpaceSymbol {
// 					w.boxesAndWallsMap[newRow][newCol] = emptySpaceSymbol
// 					w.boxesAndWallsMap[i][newCol] = boxSmallSymbol
// 					break
// 				}
// 			}
// 		} else if colOffset == -1 {
// 			for i := newCol - 1; i >= 0; i-- {
// 				if w.boxesAndWallsMap[newRow][i] == wallSymbol {
// 					return
// 				} else if w.boxesAndWallsMap[newRow][i] == emptySpaceSymbol {
// 					w.boxesAndWallsMap[newRow][newCol] = emptySpaceSymbol
// 					w.boxesAndWallsMap[newRow][i] = boxSmallSymbol
// 					break
// 				}
// 			}
// 		} else {
// 			for i := newCol + 1; i < len(w.boxesAndWallsMap[0]); i++ {
// 				if w.boxesAndWallsMap[newRow][i] == wallSymbol {
// 					return
// 				} else if w.boxesAndWallsMap[newRow][i] == emptySpaceSymbol {
// 					w.boxesAndWallsMap[newRow][newCol] = emptySpaceSymbol
// 					w.boxesAndWallsMap[newRow][i] = boxSmallSymbol
// 					break
// 				}
// 			}
// 		}
// 	}

// 	w.robotPosition = Position{
// 		rowIndex: newRow,
// 		colIndex: newCol,
// 	}
// }

func (w *Warehouse) readMap(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	w.boxesAndWallsMap = make([][]string, 0)
	w.robotMovements = make([]string, 0)
	isReadingMap := true

	for sc.Scan() {
		line := sc.Text()

		if line == "" {
			isReadingMap = false
			continue
		}

		if isReadingMap {
			w.boxesAndWallsMap = append(w.boxesAndWallsMap, strings.Split(line, ""))
		} else {
			w.robotMovements = append(w.robotMovements, strings.Split(line, "")...)
		}
	}

	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}

func (w *Warehouse) findRobot() {
	for rowIndex, row := range w.boxesAndWallsMap {
		for colIndex, cell := range row {
			if cell == robotSymbol {
				w.robotPosition = Position{
					rowIndex: rowIndex,
					colIndex: colIndex,
				}
				w.boxesAndWallsMap[rowIndex][colIndex] = emptySpaceSymbol
			}
		}
	}
}

func (w *Warehouse) transformMap() {
	transformedMap := make([][]string, len(w.boxesAndWallsMap))

	for rowIndex, row := range w.boxesAndWallsMap {
		transformedMap[rowIndex] = make([]string, 2*len(row))

		for colIndex := 0; colIndex < 2*len(row); colIndex += 2 {
			cell := row[colIndex/2]

			if cell == wallSymbol {
				transformedMap[rowIndex][colIndex] = wallSymbol
				transformedMap[rowIndex][colIndex+1] = wallSymbol
			} else if cell == boxSmallSymbol {
				transformedMap[rowIndex][colIndex] = boxLargeSymbol1
				transformedMap[rowIndex][colIndex+1] = boxLargeSymbol2
			} else if cell == emptySpaceSymbol {
				transformedMap[rowIndex][colIndex] = emptySpaceSymbol
				transformedMap[rowIndex][colIndex+1] = emptySpaceSymbol
			} else if cell == robotSymbol {
				transformedMap[rowIndex][colIndex] = robotSymbol
				transformedMap[rowIndex][colIndex+1] = emptySpaceSymbol
			}
		}
	}

	w.boxesAndWallsMap = transformedMap
}

func copyMap(original [][]string) [][]string {
	copiedMap := make([][]string, len(original))

	for i, row := range original {
		copiedMap[i] = make([]string, len(row))
		copy(copiedMap[i], row)
	}

	return copiedMap
}
