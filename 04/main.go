package main

import (
	"bufio"
	"fmt"
	"log"

	"os"
)

// --- Part 1 ---
// func main() {
// 	wordSearch := readWordSearch("./input.txt")
// 	totalWords := 0

// 	for rowIndex, line := range wordSearch {
// 		for columnIndex, char := range line {
// 			if char == "X" {
// 				fmt.Println("Found", "X", rowIndex, columnIndex)
// 				totalWords += countXmasAtPosition(&wordSearch, rowIndex, columnIndex)
// 			}
// 		}
// 	}

// 	fmt.Println("totalWords", totalWords)
// }

// func countXmasAtPosition(wordSearch *[][]string, xRowIndex int, xColumnIndex int) int {
// 	rangeToSearch := []int{-1, 0, 1}
// 	found := 0

// 	for _, offsetRow := range rangeToSearch {
// 		for _, offsetCol := range rangeToSearch {
// 			if offsetRow == 0 && offsetCol == 0 {
// 				continue
// 			}

// 			rowToCheck := xRowIndex + offsetRow
// 			columnToCheck := xColumnIndex + offsetCol

// 			if !isValidRange(wordSearch, rowToCheck, columnToCheck) {
// 				continue
// 			}

// 			if (*wordSearch)[rowToCheck][columnToCheck] == "M" {
// 				fmt.Println("Found", "M", rowToCheck, columnToCheck)

// 				// search "AS" in same direction
// 				if isValidRange(wordSearch, rowToCheck+offsetRow, columnToCheck+offsetCol) && (*wordSearch)[rowToCheck+offsetRow][columnToCheck+offsetCol] == "A" {
// 					fmt.Println("Found", "A", rowToCheck+offsetRow, columnToCheck+offsetCol)

// 					if isValidRange(wordSearch, rowToCheck+offsetRow+offsetRow, columnToCheck+offsetCol+offsetCol) && (*wordSearch)[rowToCheck+offsetRow+offsetRow][columnToCheck+offsetCol+offsetCol] == "S" {
// 						fmt.Println("Found", "S", rowToCheck+offsetRow+offsetRow, columnToCheck+offsetCol+offsetCol)
// 						found++
// 					}
// 				}
// 			}
// 		}
// 	}

//		return found
//	}
//
// --- Part 2 ---
func main() {
	wordSearch := readWordSearch("./input.txt")
	totalWords := 0

	for rowIndex, line := range wordSearch {
		for columnIndex, char := range line {
			if char == "A" {
				fmt.Println("Found", "A", rowIndex, columnIndex)
				if hasXmasCrossAtPosition(&wordSearch, rowIndex, columnIndex) {
					totalWords++
				}
			}
		}
	}

	fmt.Println("totalWords", totalWords)
}

func hasXmasCrossAtPosition(wordSearch *[][]string, xRowIndex int, xColumnIndex int) bool {
	rangeToSearch := []int{-1, 1}

	for _, offsetRow := range rangeToSearch {
		for _, offsetCol := range rangeToSearch {

			if !isValidRange(wordSearch, xRowIndex+offsetRow, xColumnIndex+offsetCol) || !isValidRange(wordSearch, xRowIndex-offsetRow, xColumnIndex-offsetCol) {
				continue
			}

			if (*wordSearch)[xRowIndex+offsetRow][xColumnIndex+offsetCol] == "M" {
				if (*wordSearch)[xRowIndex-offsetRow][xColumnIndex-offsetCol] == "S" {
					// check opposite
					char := (*wordSearch)[xRowIndex-offsetRow][xColumnIndex+offsetCol]

					if char == "M" {
						if (*wordSearch)[xRowIndex+offsetRow][xColumnIndex-offsetCol] == "S" {
							return true
						}
					} else if char == "S" {
						if (*wordSearch)[xRowIndex+offsetRow][xColumnIndex-offsetCol] == "M" {
							return true
						}
					}

				}
			}
		}
	}

	return false
}

func isValidRange(wordSearch *[][]string, rowToCheck int, columnToCheck int) bool {
	return !(rowToCheck < 0 || columnToCheck < 0 || rowToCheck > len(*wordSearch)-1 || columnToCheck > len((*wordSearch)[0])-1)
}

func readWordSearch(filePath string) [][]string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	wordSearch := make([][]string, 0)
	lineNumber := 0

	for sc.Scan() {
		line := sc.Text()
		wordSearch = append(wordSearch, make([]string, len(line)))

		for i, char := range line {
			wordSearch[lineNumber][i] = string(char)
		}

		lineNumber++
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	return wordSearch
}
