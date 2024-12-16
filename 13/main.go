package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

const buttonATokens = 3
const buttonBTokens = 1
const deltaPrize = 10000000000000 // Part 2

type Position struct {
	x int
	y int
}

type ClawMachine struct {
	buttonAMovement Position
	buttonBMovement Position
	pricePosition   Position
}

type Lobby struct {
	clawMachines []ClawMachine
}

// https://en.wikipedia.org/wiki/System_of_linear_equations
// https://www.reddit.com/r/golang/comments/zs0hqy/how_to_solve_a_set_of_simultaneous_equations_in_go/

func main() {
	lobby := Lobby{
		clawMachines: []ClawMachine{},
	}

	err := lobby.readClawMachines("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	minimumTokensNeeded := 0
	for _, machine := range lobby.clawMachines {
		tokenNeeded, calculationErr := machine.calculateMinTokensNeeded()
		if calculationErr != nil {
			log.Fatal(calculationErr)
		}

		minimumTokensNeeded += tokenNeeded
	}

	fmt.Println("minimumTokensNeeded", minimumTokensNeeded)
}

func (cm *ClawMachine) calculateMinTokensNeeded() (int, error) {
	amountOfEquations := 2
	amountOfVariables := 2

	A := mat.NewDense(amountOfEquations, amountOfVariables, []float64{
		float64(cm.buttonAMovement.x),
		float64(cm.buttonBMovement.x),
		float64(cm.buttonAMovement.y),
		float64(cm.buttonBMovement.y),
	})
	b := mat.NewVecDense(amountOfEquations, []float64{
		float64(cm.pricePosition.x),
		float64(cm.pricePosition.y)},
	)

	var x mat.VecDense
	if err := x.SolveVec(A, b); err != nil {
		return 0, err
	}

	solutionAB := x.RawVector().Data
	solutionA := int(math.Round(solutionAB[0]))
	solutionB := int(math.Round(solutionAB[1]))

	// You estimate that each button would need to be pressed no more than 100 times to win a prize.
	// if solutionA < 0 || solutionA > 100 || solutionB < 0 || solutionB > 100 { Part 1
	if solutionA < 0 || solutionB < 0 { // Part 2
		// ignore result
		return 0, nil
	}

	if solutionA*cm.buttonAMovement.x+solutionB*cm.buttonBMovement.x != cm.pricePosition.x || solutionA*cm.buttonAMovement.y+solutionB*cm.buttonBMovement.y-cm.pricePosition.y != 0 {
		// ignore result
		return 0, nil
	}

	return solutionA*buttonATokens + solutionB*buttonBTokens, nil
}

// Part 1
// func (cm *ClawMachine) calculateMinTokensNeeded() int {
// 	maxButtonBPresses := math.Min(math.Max(float64(cm.pricePosition.x/cm.buttonBMovement.x), float64(cm.pricePosition.y/cm.buttonBMovement.y)), 100)

// 	for buttonBPresses := int(maxButtonBPresses); buttonBPresses >= 0; buttonBPresses-- {
// 		movementXLeftForButtonA := cm.pricePosition.x - buttonBPresses*cm.buttonBMovement.x
// 		buttonAPressesForX := float64(movementXLeftForButtonA) / float64(cm.buttonAMovement.x)
// 		movementYLeftForButtonA := cm.pricePosition.y - buttonBPresses*cm.buttonBMovement.y
// 		buttonAPressesForY := float64(movementYLeftForButtonA) / float64(cm.buttonAMovement.y)

// 		if buttonAPressesForX == buttonAPressesForY && buttonAPressesForX <= 100 && float64(int(buttonAPressesForX))-buttonAPressesForX == 0 {
// 			return int(buttonAPressesForX)*buttonATokens + buttonBPresses*buttonBTokens
// 		}

// 	}

// 	return 0
// }

func (l *Lobby) readClawMachines(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	lines := []string{}

	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		lines = append(lines, line)

		// Wait until all instructions for a machine have been read
		if len(lines) != 3 {
			continue
		}

		buttonAx, buttonAy, buttonAerr := parseInputLine(lines[0])
		buttonBx, buttonBy, buttonBerr := parseInputLine(lines[1])
		prizeX, prizeY, prizeErr := parseInputLine(lines[2])

		if buttonAerr != nil {
			return buttonAerr
		}
		if buttonBerr != nil {
			return buttonBerr
		}
		if prizeErr != nil {
			return prizeErr
		}

		// Part 2
		machine := ClawMachine{
			buttonAMovement: Position{x: buttonAx, y: buttonAy},
			buttonBMovement: Position{x: buttonBx, y: buttonBy},
			pricePosition:   Position{x: prizeX + deltaPrize, y: prizeY + deltaPrize},
		}

		// Part 1
		// machine := ClawMachine{
		// 	buttonAMovement: Position{x: buttonAx, y: buttonAy},
		// 	buttonBMovement: Position{x: buttonBx, y: buttonBy},
		// 	pricePosition:   Position{x: prizeX, y: prizeY},
		// }
		l.clawMachines = append(l.clawMachines, machine)
		lines = []string{}
	}

	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}

func parseInputLine(line string) (x, y int, err error) {
	regex, regexErr := regexp.Compile(`(?:Button [A-B]|Prize): X[+,=](\d*), Y[+,=](\d*)`)
	if regexErr != nil {
		return 0, 0, regexErr
	}

	stringValues := regex.FindStringSubmatch(line)

	// index 0 is the full string match (to ignore)
	x, errX := strconv.Atoi(stringValues[1])
	if errX != nil {
		return 0, 0, errX
	}

	y, errY := strconv.Atoi(stringValues[2])
	if errY != nil {
		return 0, 0, errY
	}

	return x, y, nil
}
