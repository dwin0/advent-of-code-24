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

var cache = make(map[string]int)

type Pluto struct {
	stones []int
}

func main() {
	pluto := Pluto{}

	err := pluto.readStonesInput("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Part 1
	// for range 25 {
	// 	err = pluto.blink()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	totalStones := 0

	for _, stone := range pluto.stones {
		totalStones += findNumberOfStones(stone, 75)
	}

	fmt.Println(totalStones)
}

func findNumberOfStones(stone, blink int) int {
	size, found := cache[fmt.Sprintf("%v_%v", stone, blink)]

	if found {
		return size
	}

	if blink == 0 {
		return 1
	}

	if stone == 0 {
		result := findNumberOfStones(1, blink-1)
		cache[fmt.Sprintf("%v_%v", stone, blink)] = result
		return result
	}

	if number := strconv.Itoa(stone); len(number)%2 == 0 {
		stone1 := number[:len(number)/2]
		stone2 := number[len(number)/2:]

		stone1Number, conversion1Err := strconv.Atoi(stone1)
		stone2Number, conversion2Err := strconv.Atoi(stone2)

		if conversion1Err != nil {
			log.Fatal(conversion1Err)
		}
		if conversion2Err != nil {
			log.Fatal(conversion2Err)
		}

		res1 := findNumberOfStones(stone1Number, blink-1)
		res2 := findNumberOfStones(stone2Number, blink-1)
		cache[fmt.Sprintf("%v_%v", stone, blink)] = res1 + res2
		return res1 + res2
	}

	res := findNumberOfStones(stone*2024, blink-1)
	cache[fmt.Sprintf("%v_%v", stone, blink)] = res
	return res
}

// Part 1
// func (p *Pluto) blink() error {
// 	updatedStones := []int{}

// 	for _, stoneNumber := range p.stones {
// 		if stoneNumber == 0 {
// 			updatedStones = append(updatedStones, 1)
// 			continue
// 		}

// 		if number := strconv.Itoa(stoneNumber); len(number)%2 == 0 {
// 			stone1 := number[:len(number)/2]
// 			stone2 := number[len(number)/2:]

// 			stone1Number, conversion1Err := strconv.Atoi(stone1)
// 			stone2Number, conversion2Err := strconv.Atoi(stone2)

// 			if conversion1Err != nil {
// 				return conversion1Err
// 			}
// 			if conversion2Err != nil {
// 				return conversion2Err
// 			}

// 			updatedStones = append(updatedStones, stone1Number, stone2Number)
// 			continue
// 		}

// 		updatedStones = append(updatedStones, stoneNumber*2024)
// 	}

// 	p.stones = updatedStones

// 	return nil
// }

func (p *Pluto) readStonesInput(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	amountOfLines := 0

	for sc.Scan() {
		line := sc.Text()
		p.stones = make([]int, 0)

		for _, stoneNumberAsString := range strings.Split(line, " ") {
			stoneNumber, conversionErr := strconv.Atoi(stoneNumberAsString)

			if conversionErr != nil {
				return conversionErr
			}

			p.stones = append(p.stones, stoneNumber)
		}

		amountOfLines++
	}

	if err := sc.Err(); err != nil {
		return err
	}

	if amountOfLines > 1 {
		return errors.New("input file has too many lines")
	}

	return nil
}
