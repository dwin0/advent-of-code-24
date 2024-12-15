package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/dwin0/advent-of-code-24/lib"
)

type GardenRegionName string

type GardenRegionValues struct {
	name          GardenRegionName // for debugging
	area          int
	perimeter     int
	numberOfSides int
}

type Position struct {
	rowIndex    int
	columnIndex int
}

type Garden struct {
	gardenPlotsMap        [][]string
	regionNameToPositions map[GardenRegionName][]Position
	regionValues          []GardenRegionValues
	totalPrice            int
}

func main() {
	garden := Garden{
		regionNameToPositions: make(map[GardenRegionName][]Position),
		regionValues:          []GardenRegionValues{},
	}

	err := garden.readGardenPlotsMap("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	garden.assignPlotsToGroup()
	garden.calculateRegionValues()
	garden.calculateTotalPrice()

	garden.log()
}

func (g *Garden) log() {
	// fmt.Println(g.regionNameToPositions)
	// fmt.Println(g.regionValues)
	fmt.Println(g.totalPrice)
}

const markedChar = "."

func (g *Garden) assignPlotsToGroup() {
	for rowIndex, row := range g.gardenPlotsMap {
		for columnIndex, plantChar := range row {

			// plot was already assigned to a group
			if plantChar == markedChar {
				continue
			}

			// else: traverse the newly found group

			// create new unique regionName
			regionName := (GardenRegionName)(fmt.Sprintf("%s->%d/%d", plantChar, rowIndex, columnIndex))
			position := Position{
				rowIndex:    rowIndex,
				columnIndex: columnIndex,
			}

			g.regionNameToPositions[regionName] = []Position{position}
			g.gardenPlotsMap[rowIndex][columnIndex] = markedChar

			g.traverseGardenPlotGroup(rowIndex, columnIndex, plantChar, regionName)
		}
	}
}

func (g *Garden) traverseGardenPlotGroup(rowIndex, columnIndex int, plantChar string, regionName GardenRegionName) {
	for _, rowOffset := range []int{-1, 1} {
		rowIndexToTest := rowIndex + rowOffset

		if g.isInvalidRowIndex(rowIndexToTest) {
			continue
		}

		if g.gardenPlotsMap[rowIndexToTest][columnIndex] == plantChar {
			g.regionNameToPositions[regionName] = append(g.regionNameToPositions[regionName], Position{
				rowIndex:    rowIndexToTest,
				columnIndex: columnIndex,
			})

			g.gardenPlotsMap[rowIndexToTest][columnIndex] = markedChar

			g.traverseGardenPlotGroup(rowIndexToTest, columnIndex, plantChar, regionName)
		}
	}
	for _, colOffset := range []int{-1, 1} {
		columnIndexToTest := columnIndex + colOffset

		if g.isInvalidColumnIndex(columnIndexToTest) {
			continue
		}

		if g.gardenPlotsMap[rowIndex][columnIndexToTest] == plantChar {
			g.regionNameToPositions[regionName] = append(g.regionNameToPositions[regionName], Position{
				rowIndex:    rowIndex,
				columnIndex: columnIndexToTest,
			})

			g.gardenPlotsMap[rowIndex][columnIndexToTest] = markedChar

			g.traverseGardenPlotGroup(rowIndex, columnIndexToTest, plantChar, regionName)
		}
	}
}

func (g *Garden) isInvalidRowIndex(rowIndex int) bool {
	return rowIndex < 0 || rowIndex >= len(g.gardenPlotsMap)
}

func (g *Garden) isInvalidColumnIndex(columnIndex int) bool {
	return columnIndex < 0 || columnIndex >= len(g.gardenPlotsMap[0])
}

func (g *Garden) readGardenPlotsMap(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	g.gardenPlotsMap = make([][]string, 0)

	for sc.Scan() {
		line := sc.Text()

		g.gardenPlotsMap = append(g.gardenPlotsMap, strings.Split(line, ""))

	}

	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}

func (g *Garden) calculateRegionValues() {
	for regionName, positions := range g.regionNameToPositions {
		g.regionValues = append(g.regionValues, GardenRegionValues{
			name:          regionName,
			area:          len(positions),
			perimeter:     calculatePerimeter(positions),
			numberOfSides: calculateNumberOfSides(positions),
		})
	}
}

// Part 1
func calculatePerimeter(positions []Position) int {
	perimeterTotal := 0

	for i, currentPosition := range positions {
		otherPositions := append([]Position{}, positions[:i]...)
		otherPositions = append(otherPositions, positions[i+1:]...)

		perimetersToAdd := 4
		for _, otherPosition := range otherPositions {
			if (lib.AbsInt(currentPosition.rowIndex-otherPosition.rowIndex) == 1 && lib.AbsInt(currentPosition.columnIndex-otherPosition.columnIndex) == 0) || (lib.AbsInt(currentPosition.rowIndex-otherPosition.rowIndex) == 0 && lib.AbsInt(currentPosition.columnIndex-otherPosition.columnIndex) == 1) {
				// for every plant of the same region immediately next to the current plan, 1 fence less is needed
				perimetersToAdd--
			}
		}

		perimeterTotal += perimetersToAdd
	}

	return perimeterTotal
}

// Part 2
func calculateNumberOfSides(positions []Position) int {
	/*
		Idea:
		1. Collect all sides where a fence should be placed. Collect horizontal fences (rowFences) separately from vertical fences (columnFences).
			-> For row fences, only store the columns
			-> For column fences, only store the row
		2. Sort the fences
		3. Check for continuity in the sorted fences.
			-> For row fences, check if there are non-continuous columns. Continuous fences only count as 1. Every non-continuous number increases the number of sides to place a fence.
	*/

	rowFences := make(map[string][]int)
	columnFences := make(map[string][]int)

	for i, currentPosition := range positions {
		otherPositions := append([]Position{}, positions[:i]...)
		otherPositions = append(otherPositions, positions[i+1:]...)

		// check above, below, left and right: If there's another position, no need to add a fence on this side
		for _, rowOffset := range []int{-1, 1} {
			rowIndexToTest := currentPosition.rowIndex + rowOffset

			if slices.IndexFunc(otherPositions, func(p Position) bool {
				return p.rowIndex == rowIndexToTest && p.columnIndex == currentPosition.columnIndex
			}) != -1 {
				// no fence needed on this side
			} else {
				fenceRowId := fmt.Sprintf("%d_%d", currentPosition.rowIndex, rowIndexToTest)
				existing, found := rowFences[fenceRowId]

				if found {
					rowFences[fenceRowId] = append(existing, currentPosition.columnIndex)
				} else {
					rowFences[fenceRowId] = []int{currentPosition.columnIndex}
				}
			}
		}
		for _, columnOffset := range []int{-1, 1} {
			columnIndexToTest := currentPosition.columnIndex + columnOffset

			if slices.IndexFunc(otherPositions, func(p Position) bool {
				return p.columnIndex == columnIndexToTest && p.rowIndex == currentPosition.rowIndex
			}) != -1 {
				// no fence needed on this side
			} else {
				fenceColumnId := fmt.Sprintf("%d_%d", currentPosition.columnIndex, columnIndexToTest)
				existing, found := columnFences[fenceColumnId]

				if found {
					columnFences[fenceColumnId] = append(existing, currentPosition.rowIndex)
				} else {
					columnFences[fenceColumnId] = []int{currentPosition.rowIndex}
				}
			}
		}
	}

	for rowIndex := range rowFences {
		sort.Ints(rowFences[rowIndex])
	}
	for columnIndex := range columnFences {
		sort.Ints(columnFences[columnIndex])
	}

	// Check for continuous fences. These only count as 1 side.
	continuousRowFences := 0
	for _, fenceColumns := range rowFences {
		for i, column := range fenceColumns {
			if i == 0 {
				continuousRowFences++
			} else if lib.AbsInt(column-fenceColumns[i-1]) > 1 {
				continuousRowFences++
			}
		}
	}

	continuousColumnFences := 0
	for _, fenceRows := range columnFences {
		for i, row := range fenceRows {
			if i == 0 {
				continuousColumnFences++
			} else if lib.AbsInt(row-fenceRows[i-1]) > 1 {
				continuousColumnFences++
			}
		}
	}

	return continuousRowFences + continuousColumnFences
}

func (g *Garden) calculateTotalPrice() {
	totalPrice := 0

	for _, region := range g.regionValues {
		fmt.Println(region)
		// totalPrice += region.area * region.perimeter // Part 1
		totalPrice += region.area * region.numberOfSides // Part 2
	}

	g.totalPrice = totalPrice
}
