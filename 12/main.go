package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dwin0/advent-of-code-24/lib"
)

type GardenRegionName string

type GardenRegionValues struct {
	name      GardenRegionName // for debugging
	area      int
	perimeter int
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
			name:      regionName,
			area:      len(positions),
			perimeter: calculatePerimeter(positions),
		})
	}
}

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

func (g *Garden) calculateTotalPrice() {
	totalPrice := 0

	for _, region := range g.regionValues {
		fmt.Println(region)
		totalPrice += region.area * region.perimeter
	}

	g.totalPrice = totalPrice
}
