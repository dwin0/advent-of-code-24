package main

import (
	"bufio"
	"fmt"
	"log"
	"math"

	"os"
	"strconv"
	"strings"

	"github.com/dwin0/advent-of-code-24/lib"
)

// --- Part 1 ---
// func main() {
// 	reports := readLevels("./input.txt")

// 	totalSafeReports := 0

// 	for _, report := range reports {
// 		isSafeReport := true
// 		signOfReport := 0

// 		for indexLevel := range report {
// 			if len(report)-1 == indexLevel {
// 				break
// 			}

// 			levelsToCompare := report[indexLevel : indexLevel+2]
// 			currentLevel := levelsToCompare[0]
// 			nextLevel := levelsToCompare[1]

// 			diff := int(nextLevel) - int(currentLevel)
// 			diffAbs := lib.AbsInt(diff)

// 			if diffAbs < 1 || diffAbs > 3 {
// 				isSafeReport = false
// 				break
// 			}

// 			isDecreasing := math.Signbit(float64(diff))

// 			if signOfReport != 0 {
// 				if signOfReport == 1 && isDecreasing {
// 					isSafeReport = false
// 					break
// 				}
// 				if signOfReport == -1 && !isDecreasing {
// 					isSafeReport = false
// 					break
// 				}
// 			}

// 			if isDecreasing {
// 				signOfReport = -1
// 			} else {
// 				signOfReport = 1
// 			}
// 		}

// 		if isSafeReport {
// 			totalSafeReports++
// 		}
// 	}

// 	fmt.Printf("totalSafeReports is %d\n", totalSafeReports)
// }

// --- Part 2 ---
func main() {
	reports := readLevels("./input.txt")

	totalSafeReports := 0

	for _, report := range reports {
		if isSafeReport(report, true) {
			totalSafeReports++
		}
	}

	fmt.Printf("totalSafeReports is %d\n", totalSafeReports)
}

func isSafeReport(report []uint8, isFullReport bool) bool {
	isSafe := true
	signOfReport := 0
	indexLevel := 0

	for indexLevel = range report {
		if len(report)-1 == indexLevel {
			break
		}

		currentLevel := report[indexLevel]
		nextLevel := report[indexLevel+1]

		diff := int(nextLevel) - int(currentLevel)
		diffAbs := lib.AbsInt(diff)

		if diffAbs < 1 || diffAbs > 3 {
			isSafe = false
			break
		}

		isDecreasing := math.Signbit(float64(diff))

		if signOfReport != 0 {
			if signOfReport == 1 && isDecreasing {
				isSafe = false
				break
			}
			if signOfReport == -1 && !isDecreasing {
				isSafe = false
				break
			}
		}

		if isDecreasing {
			signOfReport = -1
		} else {
			signOfReport = 1
		}
	}

	if isSafe {
		return true
	}

	if !isFullReport {
		return false
	}

	if indexLevel > 0 {
		return isSafeReport(lib.RemoveIndex(report, indexLevel-1), false) || isSafeReport(lib.RemoveIndex(report, indexLevel), false) || isSafeReport(lib.RemoveIndex(report, indexLevel+1), false)
	} else {
		return isSafeReport(lib.RemoveIndex(report, indexLevel), false) || isSafeReport(lib.RemoveIndex(report, indexLevel+1), false)
	}

}

func readLevels(filePath string) [][]uint8 {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	reports := make([][]uint8, 0)
	lineNumber := 0

	for sc.Scan() {
		reportLine := sc.Text()
		levels := strings.Fields(reportLine)
		reports = append(reports, make([]uint8, len(levels)))

		for i, stringValue := range levels {
			intValue, err := strconv.Atoi(stringValue)
			if err != nil {
				log.Fatal(err)
			}

			reports[lineNumber][i] = uint8(intValue)
		}

		lineNumber++
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	return reports
}
