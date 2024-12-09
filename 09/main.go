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

type Disk struct {
	diskMap       []int
	fileBlocks    []string
	freePositions []int
	checksum      int
}

func main() {
	disk := Disk{
		diskMap:       make([]int, 0),
		fileBlocks:    make([]string, 0),
		freePositions: make([]int, 0),
	}
	err := disk.readDiskMap("./input-small.txt")

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
	// fmt.Println(m.freePositions)
	fmt.Println("Checksum: ", m.checksum)
}

func (m *Disk) moveFileBlocks() {
	freePositionsIndex := 0
	numberOfFreePositions := len(m.freePositions)

	for i := len(m.fileBlocks) - 1; i >= 0; i-- {
		if m.fileBlocks[i] == "." {
			continue
		}

		if freePositionsIndex+1 == numberOfFreePositions {
			return
		}

		nextFreePosition := m.freePositions[freePositionsIndex]
		if nextFreePosition > i {
			return
		}

		val := m.fileBlocks[i]
		m.fileBlocks[nextFreePosition] = val
		m.fileBlocks[i] = "."

		freePositionsIndex++
	}
}

func (m *Disk) calculateChecksum() error {
	for position, fileIdString := range m.fileBlocks {
		if fileIdString == "." {
			return nil
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

			for range val {
				m.fileBlocks = append(m.fileBlocks, strconv.Itoa(fileId))
			}
		} else {
			freeRange := makeRange(len(m.fileBlocks), len(m.fileBlocks)+val-1)
			m.freePositions = append(m.freePositions, freeRange...)

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
