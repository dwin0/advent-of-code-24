package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

const buttonATokens = 3
const buttonBTokens = 1

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
		minimumTokensNeeded += machine.calculateMinTokensNeeded()
	}

	fmt.Println("minimumTokensNeeded", minimumTokensNeeded)
}

func (cm *ClawMachine) calculateMinTokensNeeded() int {
	maxButtonBPresses := math.Min(math.Max(float64(cm.pricePosition.x/cm.buttonBMovement.x), float64(cm.pricePosition.y/cm.buttonBMovement.y)), 100)

	for buttonBPresses := int(maxButtonBPresses); buttonBPresses >= 0; buttonBPresses-- {
		movementXLeftForButtonA := cm.pricePosition.x - buttonBPresses*cm.buttonBMovement.x
		buttonAPressesForX := float64(movementXLeftForButtonA) / float64(cm.buttonAMovement.x)
		movementYLeftForButtonA := cm.pricePosition.y - buttonBPresses*cm.buttonBMovement.y
		buttonAPressesForY := float64(movementYLeftForButtonA) / float64(cm.buttonAMovement.y)

		if buttonAPressesForX == buttonAPressesForY && buttonAPressesForX <= 100 && float64(int(buttonAPressesForX))-buttonAPressesForX == 0 {
			return int(buttonAPressesForX)*buttonATokens + buttonBPresses*buttonBTokens
		}

	}

	return 0
}

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

		machine := ClawMachine{
			buttonAMovement: Position{x: buttonAx, y: buttonAy},
			buttonBMovement: Position{x: buttonBx, y: buttonBy},
			pricePosition:   Position{x: prizeX, y: prizeY},
		}
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
