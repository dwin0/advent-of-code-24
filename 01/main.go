package main

import (
	"bufio"
	"fmt"
	"log"

	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/dwin0/advent-of-code-24/lib"
)

func main() {
	left, right := readTextLeftAndRightNumbers("./input.txt")

	if len(left) != len(right) {
		log.Fatal("Arrays have a different length")
	}

	// --- Part 1 ---
	sort.Ints(left)
	sort.Ints(right)

	totalDistance := 0

	for i, l := range left {
		r := right[i]

		distance := lib.AbsInt(l - r)
		totalDistance += distance
	}

	fmt.Printf("diff is %d\n", totalDistance)

	// --- Part 2 ---
	similarityScore := 0

	numberOfItemsInRightListLookup := lib.NumberOfItemsLookupMap(right)

	for _, valueLeft := range left {
		numberOfItemsInRightList := numberOfItemsInRightListLookup[valueLeft]
		similarityScore += valueLeft * numberOfItemsInRightList
	}

	fmt.Printf("similarityScore is %d\n", similarityScore)
}

func readTextLeftAndRightNumbers(filePath string) (left []int, right []int) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	left = make([]int, 0)
	right = make([]int, 0)

	for sc.Scan() {
		line := sc.Text()
		lineContent := strings.Fields(line)

		l, errL := strconv.Atoi(lineContent[0])
		r, errR := strconv.Atoi(lineContent[1])

		if errL != nil {
			log.Fatal(errL)
		}
		if errR != nil {
			log.Fatal(errR)
		}

		left = append(left, l)
		right = append(right, r)
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	return
}
