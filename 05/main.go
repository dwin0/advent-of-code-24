package main

import (
	"bufio"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"os"
)

func main() {
	rules := readRules("./rules.txt")
	pageBeforeLookupMap, pageAfterLookupMap := createRuleLookupMap(rules)

	updatePages := readUpdatePages("./update-page-numbers.txt")

	sumOfMiddleNumbers := 0

	for _, singleUpdate := range updatePages {
		if isUpdateValid(singleUpdate, &pageBeforeLookupMap, &pageAfterLookupMap) {
			// --- Part 1 ---
			// sumOfMiddleNumbers += findMiddleNumber(singleUpdate)
		} else {
			// --- Part 2 ---
			validUpdate := sortUpdate(&singleUpdate, &pageBeforeLookupMap, &pageAfterLookupMap)
			sumOfMiddleNumbers += findMiddleNumber(*validUpdate)
		}
	}

	fmt.Println("Sum is", sumOfMiddleNumbers)
}

func sortUpdate(singleUpdate *[]int, beforeRules *map[int][]int, afterRules *map[int][]int) *[]int {
	for i, currentPage := range *singleUpdate {
		pagesBefore := (*singleUpdate)[:i]
		pagesAfter := (*singleUpdate)[i+1:]

		currentPageMustBeBeforeThesePages, foundBeforeRules := (*beforeRules)[currentPage]
		currentPageMustBeAfterThesePages, foundAfterRules := (*afterRules)[currentPage]

		if foundBeforeRules {
			for _, page := range pagesBefore {
				if slices.Contains(currentPageMustBeBeforeThesePages, page) {
					indexOfPage := slices.Index(*singleUpdate, page)

					(*singleUpdate)[i] = page
					(*singleUpdate)[indexOfPage] = currentPage

					return sortUpdate(singleUpdate, beforeRules, afterRules)
				}
			}
		}

		if foundAfterRules {
			for _, page := range pagesAfter {
				if slices.Contains(currentPageMustBeAfterThesePages, page) {
					indexOfPage := slices.Index(*singleUpdate, page)

					(*singleUpdate)[i] = page
					(*singleUpdate)[indexOfPage] = currentPage

					return sortUpdate(singleUpdate, beforeRules, afterRules)
				}
			}
		}
	}

	return singleUpdate
}

func isUpdateValid(singleUpdate []int, beforeRules *map[int][]int, afterRules *map[int][]int) bool {
	for i, currentPage := range singleUpdate {
		pagesBefore := singleUpdate[:i]
		pagesAfter := singleUpdate[i+1:]

		currentPageMustBeBeforeThesePages, foundBeforeRules := (*beforeRules)[currentPage]
		currentPageMustBeAfterThesePages, foundAfterRules := (*afterRules)[currentPage]

		if foundBeforeRules {
			for _, page := range pagesBefore {
				if slices.Contains(currentPageMustBeBeforeThesePages, page) {
					return false
				}
			}
		}

		if foundAfterRules {
			for _, page := range pagesAfter {
				if slices.Contains(currentPageMustBeAfterThesePages, page) {
					return false
				}
			}
		}
	}

	return true
}

func findMiddleNumber(singleUpdate []int) int {
	return singleUpdate[len(singleUpdate)/2]
}

func createRuleLookupMap(rules [][]int) (beforeRules map[int][]int, afterRules map[int][]int) {
	beforeRules = make(map[int][]int)
	afterRules = make(map[int][]int)

	for _, pageRules := range rules {
		pageBefore := pageRules[0]
		pageAfter := pageRules[1]

		beforePages, foundBefore := beforeRules[pageBefore]
		afterPages, foundAfter := afterRules[pageAfter]

		if foundBefore {
			beforeRules[pageBefore] = append(beforePages, pageAfter)
		} else {
			beforeRules[pageBefore] = []int{pageAfter}
		}

		if foundAfter {
			afterRules[pageAfter] = append(afterPages, pageBefore)
		} else {
			afterRules[pageAfter] = []int{pageBefore}
		}
	}

	return
}

func readRules(filePath string) [][]int {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	rules := make([][]int, 0)
	lineNumber := 0

	for sc.Scan() {
		line := sc.Text()
		rules = append(rules, make([]int, 2))
		page1, page2, ok := strings.Cut(line, "|")

		if !ok {
			log.Fatal("could not parse rule on line", line)
		}

		page1Num, errPage1 := strconv.Atoi(page1)
		page2Num, errPage2 := strconv.Atoi(page2)

		if errPage1 != nil || errPage2 != nil {
			log.Fatal("error parsing number on line", line)
		}

		rules[lineNumber] = []int{page1Num, page2Num}

		lineNumber++
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	return rules
}

func readUpdatePages(filePath string) [][]int {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	updatePages := make([][]int, 0)

	for sc.Scan() {
		line := sc.Text()
		singleUpdate := []int{}

		for _, num := range strings.Split(line, ",") {
			pageNum, err := strconv.Atoi(num)
			if err != nil {
				log.Fatal(err)
			}
			singleUpdate = append(singleUpdate, pageNum)
		}

		updatePages = append(updatePages, singleUpdate)
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	return updatePages
}
