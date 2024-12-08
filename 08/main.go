package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"os"
)

type Coordinate struct {
	row    int
	column int
}

type Map struct {
	antennaCoordinates map[string][]Coordinate
	antinodesMap       [][]string
	totalAntinodes     int
}

func main() {
	m := Map{}
	err := m.readMapInput("./input.txt")

	if err != nil {
		log.Fatal(err)
	}

	m.calculatePositionsOfAntinode()
	m.countAntinodes()
	m.log()
}

func (m *Map) log() {
	fmt.Println(m.antennaCoordinates)
	fmt.Println(m.antinodesMap)
	fmt.Println("Total antinodes: ", m.totalAntinodes)
}

func (m *Map) countAntinodes() {
	for _, row := range m.antinodesMap {
		for _, char := range row {
			if char == "#" {
				m.totalAntinodes++
			}
		}
	}
}

func (m *Map) calculatePositionsOfAntinode() {
	for _, coordinates := range m.antennaCoordinates {

		// skip single nodes
		if len(coordinates) == 1 {
			continue
		}

		for i := 0; i < len(coordinates); i++ {
			c1 := coordinates[i]
			// Part 2
			m.antinodesMap[c1.row][c1.column] = "#"

			for j := i + 1; j < len(coordinates); j++ {
				c2 := coordinates[j]

				rowDistance := c2.row - c1.row
				colDistance := c2.column - c1.column
				antinode1 := Coordinate{row: c1.row - rowDistance, column: c1.column - colDistance}
				antinode2 := Coordinate{row: c2.row + rowDistance, column: c2.column + colDistance}

				// Part 2
				for m.isInMapBoundary(antinode1) {
					m.antinodesMap[antinode1.row][antinode1.column] = "#"
					antinode1 = Coordinate{row: antinode1.row - rowDistance, column: antinode1.column - colDistance}
				}
				for m.isInMapBoundary(antinode2) {
					m.antinodesMap[antinode2.row][antinode2.column] = "#"
					antinode2 = Coordinate{row: antinode2.row + rowDistance, column: antinode2.column + colDistance}
				}

				// Part 1
				// if m.isInMapBoundary(antinode1) {
				// 	m.antinodesMap[antinode1.row][antinode1.column] = "#"
				// }
				// if m.isInMapBoundary(antinode2) {
				// 	m.antinodesMap[antinode2.row][antinode2.column] = "#"
				// }
			}
		}
	}
}

func (m *Map) isInMapBoundary(antinode Coordinate) bool {
	return antinode.row >= 0 && antinode.row < len(m.antinodesMap) && antinode.column >= 0 && antinode.column < len(m.antinodesMap[0])
}

func (m *Map) readMapInput(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	rowIndex := 0

	m.antennaCoordinates = make(map[string][]Coordinate, 0)
	m.antinodesMap = make([][]string, 0)

	for sc.Scan() {
		line := sc.Text()
		coordinates := strings.Split(line, "")

		for colIndex, char := range coordinates {
			if char != "." {
				nextAntenna := Coordinate{row: rowIndex, column: colIndex}
				existingAntennas, ok := m.antennaCoordinates[char]

				if ok {
					m.antennaCoordinates[char] = append(existingAntennas, nextAntenna)
				} else {
					m.antennaCoordinates[char] = []Coordinate{nextAntenna}
				}
			}
		}

		m.antinodesMap = append(m.antinodesMap, strings.Split(strings.Repeat(".", len(coordinates)), ""))
		rowIndex++
	}

	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}
