package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"os"
)

type Equations struct {
	equationLines         [][]int
	sumOfCorrectEquations int
}

func main() {
	equations := Equations{
		sumOfCorrectEquations: 0,
	}
	err := equations.readEquations("./input.txt")

	if err != nil {
		log.Fatal(err)
	}

	err = equations.tryOperators()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("sumOfCorrectEquations", equations.sumOfCorrectEquations)
}

func (equations *Equations) tryOperators() error {
	for _, eq := range equations.equationLines {
		res := eq[0]
		numbers := eq[1:]

		correct, err := tryEquationLine(res, numbers)

		if err != nil {
			return err
		}

		if correct {
			equations.sumOfCorrectEquations += res
		}
	}

	return nil
}

func (equations *Equations) readEquations(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	equations.equationLines = make([][]int, 0)

	for sc.Scan() {
		line := sc.Text()

		lineWithoutColon := strings.Replace(line, ":", "", 1)
		inputNumbers := strings.Split(lineWithoutColon, " ")

		equationLine := make([]int, len(inputNumbers))

		for i, numberAsString := range inputNumbers {
			number, err := strconv.Atoi(numberAsString)

			if err != nil {
				return err
			}

			equationLine[i] = number
		}

		equations.equationLines = append(equations.equationLines, equationLine)
	}

	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}

func tryEquationLine(result int, numbers []int) (bool, error) {
	// Part 1 (Part 2 has more operators)
	/* 1 2 3 4
	[
		[+ + +],
		[+ * +],
		[* + *],
		[+ + *],
		[+ * *],
		[* * +],
		[* * *],
	]
	*/

	// Part 1
	// operators := []string{"+", "*"}

	// Part 2
	operators := []string{"+", "*", "||"}
	permutations := getPermutations(operators, len(numbers)-1)

	for _, permutation := range permutations {
		intermediateResult := numbers[0]

		for i, operator := range permutation {
			nextNumber := numbers[i+1]

			if operator == "+" {
				intermediateResult += nextNumber
			} else if operator == "*" {
				intermediateResult *= nextNumber
			} else {
				res, err := strconv.Atoi(strconv.Itoa(intermediateResult) + strconv.Itoa(nextNumber))

				if err != nil {
					return false, err
				}

				intermediateResult = res
			}
		}

		if result == intermediateResult {
			return true, nil
		}
	}

	return false, nil
}

// https://www.geeksforgeeks.org/print-all-the-permutation-of-length-l-using-the-elements-of-an-array-iterative/
func getPermutations(arr []string, L int) [][]string {
	lenArr := len(arr)
	permutations := make([][]string, 0)

	// There can be (len)^L permutations
	numPermutations := int(math.Pow(float64(lenArr), float64(L)))
	for i := 0; i < numPermutations; i++ {
		// Convert i to len-th base
		permutations = append(permutations, convertToLenThBase(i, arr, lenArr, L))
	}

	return permutations
}

// Convert the number to Lth base and print the sequence
func convertToLenThBase(n int, arr []string, lenArr, L int) []string {
	permutation := make([]string, L)

	// Sequence is of length L
	for i := 0; i < L; i++ {
		// Print the ith element of the sequence
		permutation[i] += arr[n%lenArr] // Print current element in base-L
		n = n / lenArr                  // Update n for the next iteration
	}

	return permutation
}
