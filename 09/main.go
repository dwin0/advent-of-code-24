package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type FreePosition struct {
	index   int
	lengths int
}

type FilePosition struct {
	index   int
	value   int
	lengths int
}

func (fp FreePosition) String() string {
	return fmt.Sprintf("index %d, lengths %d;", fp.index, fp.lengths)
}

func (fp FilePosition) String() string {
	return fmt.Sprintf("index %d, value %d, lengths %d;", fp.index, fp.value, fp.lengths)
}

type Disk struct {
	diskMap             []int
	fileBlocks          []string       // Part 1
	filePositions       []FilePosition // Part 2
	freePositionIndices []int          // Part 1
	freePositions       []FreePosition // Part 2
	checksum            int
}

func main() {
	disk := Disk{
		diskMap:             make([]int, 0),
		fileBlocks:          make([]string, 0),
		freePositionIndices: make([]int, 0),
		freePositions:       make([]FreePosition, 0),
		filePositions:       make([]FilePosition, 0),
	}

	err := disk.readDiskMap("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	disk.unpackDiskMap()
	disk.moveFileBlocks()

	err = disk.calculateChecksum()
	if err != nil {
		log.Fatal(err)
	}

	disk.log()
}

func (m *Disk) log() {
	// fmt.Println(m.diskMap)
	// fmt.Println(m.fileBlocks)
	// fmt.Println(m.freePositionIndices)
	// fmt.Println(m.freePositions)
	// fmt.Println(m.filePositions)
	fmt.Println("Checksum: ", m.checksum)
}

// Part 1
// func (m *Disk) moveFileBlocks() {
// 	freePositionsIndex := 0
// 	numberOfFreePositions := len(m.freePositionIndices)

// 	for i := len(m.fileBlocks) - 1; i >= 0; i-- {
// 		if m.fileBlocks[i] == "." {
// 			continue
// 		}

// 		if freePositionsIndex+1 == numberOfFreePositions {
// 			return
// 		}

// 		nextFreePosition := m.freePositionIndices[freePositionsIndex]
// 		if nextFreePosition > i {
// 			return
// 		}

// 		val := m.fileBlocks[i]
// 		m.fileBlocks[nextFreePosition] = val
// 		m.fileBlocks[i] = "."

// 		freePositionsIndex++
// 	}
// }

// Part 2
func (m *Disk) moveFileBlocks() {
	for i := len(m.filePositions) - 1; i >= 0; i-- {
		fileToMove := m.filePositions[i]

		for fpIndex, freePosition := range m.freePositions {
			if fileToMove.lengths <= freePosition.lengths && fileToMove.index > freePosition.index {

				for l := range fileToMove.lengths {
					m.fileBlocks[freePosition.index+l] = strconv.Itoa(fileToMove.value)
					m.fileBlocks[fileToMove.index+l] = "."
				}
				m.freePositions[fpIndex].lengths -= fileToMove.lengths
				m.freePositions[fpIndex].index += fileToMove.lengths
				// no need to update m.filePositions, as there's only 1 try to move the files

				break
			}
		}
	}
}

func (m *Disk) calculateChecksum() error {
	for position, fileIdString := range m.fileBlocks {
		// Part 1: Early return
		// if fileIdString == "." {
		// 	return nil
		// }

		if fileIdString == "." {
			continue
		}

		fileId, conversionError := strconv.Atoi(fileIdString)

		if conversionError != nil {
			return conversionError
		}

		m.checksum += position * fileId
	}

	return nil
}

func (m *Disk) unpackDiskMap() {
	for i, val := range m.diskMap {
		if i%2 == 0 {
			// can be larger than 10 (meaning 2 digits)
			fileId := i / 2

			m.filePositions = append(m.filePositions, FilePosition{index: len(m.fileBlocks), value: fileId, lengths: val})

			for range val {
				m.fileBlocks = append(m.fileBlocks, strconv.Itoa(fileId))
			}
		} else if val == 0 {
			// empty space with lengths 0
			continue
		} else {
			freeRange := makeRange(len(m.fileBlocks), len(m.fileBlocks)+val-1)
			m.freePositionIndices = append(m.freePositionIndices, freeRange...)
			m.freePositions = append(m.freePositions, FreePosition{index: len(m.fileBlocks), lengths: val})

			emptySpace := strings.Split(strings.Repeat(".", val), "")
			m.fileBlocks = append(m.fileBlocks, emptySpace...)
		}
	}
}

// https://stackoverflow.com/questions/39868029/how-to-generate-a-sequence-of-numbers
func makeRange(min int, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func (m *Disk) readDiskMap(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	lineNumber := 1

	for sc.Scan() {
		if lineNumber > 1 {
			return errors.New("too many input lines, should only be 1")
		}

		line := sc.Text()

		for _, number := range strings.Split(line, "") {
			num, conversionErr := strconv.Atoi(number)

			if conversionErr != nil {
				return conversionErr
			}

			m.diskMap = append(m.diskMap, num)
		}

		lineNumber++
	}

	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}
