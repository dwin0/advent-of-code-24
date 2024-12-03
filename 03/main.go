package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

// --- Part 1 ---
// func main() {
// 	regexMul, _ := regexp.Compile(`mul\(\d{1,3}\,\d{1,3}\)`)
// 	regexNumbers, _ := regexp.Compile(`mul\((\d{1,3})\,(\d{1,3})\)`)
// 	multiplicationsFound := make([]string, 0)
// 	sum := 0

// 	for _, instructionLine := range readMemory("./input.txt") {
// 		multiplicationsFound = append(multiplicationsFound, regexMul.FindAllString(instructionLine, -1)...)
// 	}

// 	for _, multiplication := range multiplicationsFound {
// 		// example: [mul(2,4) 2 4]
// 		submatch := regexNumbers.FindStringSubmatch(multiplication)
// 		if len(submatch) != 3 {
// 			log.Fatal("submatch has not correct shape", submatch)
// 		}

// 		multiplicand, err1 := strconv.Atoi(submatch[1])
// 		multiplier, err2 := strconv.Atoi(submatch[2])
// 		if err1 != nil || err2 != nil {
// 			log.Fatalf("numbers are not in expected format.\nmultiplicand: %v\n%v", multiplicand, multiplier)
// 		}

// 		sum += multiplicand * multiplier
// 	}

// 	fmt.Printf("Sum is: %d\n", sum)
// }

// --- Part 2 ---
func main() {
	regexMul, regexMulErr := regexp.Compile(`(do\(\))|(don\'t\(\))|(mul\(\d{1,3}\,\d{1,3}\))`)
	regexNumbers, regexNumbersErr := regexp.Compile(`mul\((\d{1,3})\,(\d{1,3})\)`)

	if regexMulErr != nil || regexNumbersErr != nil {
		log.Fatalf("Could not compile regex.\nregexMulErr: %v\nregexNumbersErr: %v", regexMulErr, regexNumbersErr)
	}

	instructionsFound := make([]string, 0)
	sum := 0
	isNextInstructionActive := true

	for _, instructionLine := range readMemory("./input.txt") {
		instructionsFound = append(instructionsFound, regexMul.FindAllString(instructionLine, -1)...)
	}

	for _, instruction := range instructionsFound {
		if instruction == "do()" {
			isNextInstructionActive = true
			continue
		}

		if instruction == "don't()" {
			isNextInstructionActive = false
			continue
		}

		if !isNextInstructionActive {
			continue
		}

		// example: [mul(2,4) 2 4]
		submatch := regexNumbers.FindStringSubmatch(instruction)
		if len(submatch) != 3 {
			log.Fatal("submatch has not correct shape", submatch)
		}

		multiplicand, err1 := strconv.Atoi(submatch[1])
		multiplier, err2 := strconv.Atoi(submatch[2])
		if err1 != nil || err2 != nil {
			log.Fatalf("numbers are not in expected format.\nmultiplicand: %v\n%v", multiplicand, multiplier)
		}

		sum += multiplicand * multiplier
	}

	fmt.Printf("Sum is: %d\n", sum)
}

func readMemory(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	memoryInstructionLines := make([]string, 0)

	for sc.Scan() {
		line := sc.Text()
		memoryInstructionLines = append(memoryInstructionLines, line)
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	return memoryInstructionLines
}
